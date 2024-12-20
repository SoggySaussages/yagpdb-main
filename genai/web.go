package genai

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/botlabs-gg/sgpdb/v2/common"
	"github.com/botlabs-gg/sgpdb/v2/common/cplogs"
	"github.com/botlabs-gg/sgpdb/v2/common/featureflags"
	"github.com/botlabs-gg/sgpdb/v2/genai/models"
	"github.com/botlabs-gg/sgpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/sgpdb/v2/lib/dstate"
	"github.com/botlabs-gg/sgpdb/v2/web"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"goji.io"
	"goji.io/pat"
)

//go:embed assets/genai.html
var PageHTML string

//go:embed assets/genai_commands.html
var CommandsPageHTML string

type ConextKey int

type FormData struct {
	Enabled        bool
	Provider       int
	Model          string
	Key            string
	BaseCmdEnabled bool
	ResetToken     bool
}

type CommandFormData struct {
	ID                      int64
	GuildID                 int64
	Enabled                 bool
	Name                    string
	Aliases                 string
	Prompt                  string
	AllowInput              bool
	MaxTokens               int `valid:"1,512"`
	AutodeleteResponse      bool
	AutodeleteTrigger       bool
	AutodeleteResponseDelay int     `valid:"0,2678400"`
	AutodeleteTriggerDelay  int     `valid:"0,2678400"`
	Channels                []int64 `valid:"channel,true"`
	Categories              []int64 `valid:"channel,true"`
	ChannelsWhitelistMode   bool
	Roles                   []int64 `valid:"role,true"`
	RolesWhitelistMode      bool
}

const (
	ConextKeyConfig ConextKey = iota
)

var panelLogKey = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "genai_settings_updated", FormatString: "Updated genai settings"})

func (p *Plugin) InitWeb() {
	web.AddHTMLTemplate("genai/assets/genai.html", PageHTML)
	web.AddSidebarItem(web.SidebarCategoryGenAI, &web.SidebarItem{
		Name: "General",
		URL:  "genai",
		Icon: "fas fa-cog",
	})

	genaiMux := goji.SubMux()
	web.CPMux.Handle(pat.New("/genai/*"), genaiMux)
	web.CPMux.Handle(pat.New("/genai"), genaiMux)
	web.CPMux.Handle(pat.New("/genai/commands"), genaiMux)

	genaiMux.Use(web.RequireServerAdminMiddleware)
	genaiMux.Use(baseData)

	// Get just renders the template, so let the renderhandler do all the work
	genaiMux.Handle(pat.Get(""), web.RenderHandler(nil, "cp_genai"))
	genaiMux.Handle(pat.Get("/"), web.RenderHandler(nil, "cp_genai"))

	genaiMux.Handle(pat.Post(""), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), FormData{}))
	genaiMux.Handle(pat.Post("/"), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), FormData{}))

	web.AddHTMLTemplate("genai/assets/genai_commands.html", CommandsPageHTML)
	web.AddSidebarItem(web.SidebarCategoryGenAI, &web.SidebarItem{
		Name: "Commands",
		URL:  "genai/commands",
		Icon: "fas fa-terminal",
	})
	getHandler := web.ControllerHandler(HandleCommands, "cp_genai_commands")

	genaiMux.Handle(pat.Get("/commands"), getHandler)
	genaiMux.Handle(pat.Get("/commands/"), getHandler)

	genaiMux.Handle(pat.Post("/commands/new"),
		web.ControllerPostHandler(HandleCreateCommand, getHandler, CommandFormData{}))

	genaiMux.Handle(pat.Post("/commands/:commandID/update"),
		web.ControllerPostHandler(CommandMiddleware(HandleUpdateCommand), getHandler, CommandFormData{}))

	genaiMux.Handle(pat.Post("/commands/:commandID/delete"),
		web.ControllerPostHandler(CommandMiddleware(HandleDeleteCommand), getHandler, nil))
}

// Adds the current config to the context
func baseData(inner http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		guild, tmpl := web.GetBaseCPContextData(r.Context())
		config, err := GetConfig(guild.ID)
		if web.CheckErr(tmpl, err, "Failed retrieving genai config :'(", web.CtxLogger(r.Context()).Error) {
			web.LogIgnoreErr(web.Templates.ExecuteTemplate(w, "cp_genai", tmpl))
			return
		}
		if len(config.Key) > 0 {
			tmpl["KeySet"] = true
		}
		tmpl["GenAIConfig"] = config
		tmpl["GenAIProviders"] = GenAIProviders
		inner.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ConextKeyConfig, config)))
	}

	return http.HandlerFunc(mw)
}

func HandlePostGenAI(w http.ResponseWriter, r *http.Request) interface{} {
	ctx := r.Context()
	guild, tmpl := web.GetBaseCPContextData(ctx)
	tmpl["VisibleURL"] = "/manage/" + discordgo.StrID(guild.ID) + "/genai/"

	ok := ctx.Value(common.ContextKeyFormOk).(bool)
	formData := ctx.Value(common.ContextKeyParsedForm).(*FormData)
	newConf := &models.GenaiConfig{
		GuildID:        guild.ID,
		Enabled:        formData.Enabled,
		Provider:       formData.Provider,
		Model:          formData.Model,
		BaseCMDEnabled: formData.BaseCmdEnabled,
	}
	if !formData.Enabled {
		formData.Key = ""
		formData.ResetToken = true
	}

	if !ok {
		return tmpl
	}

	var saveNewKey bool
	provider := GenAIProviderFromID(newConf.Provider)

	// see if selected model is a part of selected provider's map, revert to a
	// default if not
	var found bool
	for _, model := range *provider.ModelMap() {
		if model == newConf.Model {
			found = true
			break
		}
	}
	if !found {
		newConf.Model = provider.DefaultModel()
	}

	newConfFakeKey := *newConf
	tmpl["GenAIConfig"] = &newConfFakeKey

	conf, err := GetConfig(guild.ID)
	if err != nil {
		conf = &models.GenaiConfig{}
	}

	if conf.Provider != formData.Provider {
		// provider changed, reset token
		formData.Key = ""
		formData.ResetToken = true
	}

	partialGuildState := &dstate.GuildState{ID: guild.ID, OwnerID: guild.OwnerID}
	if formData.Key != "" {
		newConf.Key, err = encryptAPIToken(partialGuildState, formData.Key)
		if web.CheckErr(tmpl, err, "Failed encrypting your API token to save", web.CtxLogger(ctx).Error) {
			return tmpl
		}
	}

	blacklistColumns := []string{}
	keyChanged := !bytes.Equal(newConf.Key, conf.Key) && len(newConf.Key) > 0
	saveNewKey = provider.KeyRequired() && (formData.ResetToken || keyChanged)
	if saveNewKey {
		if len(newConf.Key) > 0 {
			if web.CheckErr(tmpl, provider.ValidateAPIToken(partialGuildState, string(formData.Key)), "Your API token is invalid.", web.CtxLogger(ctx).Error) {
				return tmpl
			}
		}
	} else {
		newConf.Key = conf.Key
		blacklistColumns = append(blacklistColumns, "key")
	}

	if conf.GuildID == 0 { // config has never been saved
		err = newConf.InsertG(r.Context(), boil.Infer())
	} else {
		_, err = newConf.UpdateG(ctx, boil.Blacklist(blacklistColumns...))
	}
	if web.CheckErr(tmpl, err, "Failed saving config :'(", web.CtxLogger(ctx).Error) {
		return tmpl
	}

	err = featureflags.UpdatePluginFeatureFlags(guild.ID, &Plugin{})
	if err != nil {
		web.CtxLogger(ctx).WithError(err).Error("failed updating feature flags")
	}

	// Ensure the "reset token" button appears on the webpage
	tmpl["KeySet"] = len(newConf.Key) > 0

	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKey))

	return tmpl.AddAlerts(web.SucessAlert("Saved settings"))
}

// Servers the command page with current config
func HandleCommands(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ctx := r.Context()
	activeGuild, templateData := web.GetBaseCPContextData(ctx)

	commands, err := models.GenaiCommands(models.GenaiCommandWhere.GuildID.EQ(activeGuild.ID)).AllG(ctx)
	if err != nil {
		return templateData, err
	}

	templateData["Commands"] = commands

	templateData["VisibleURL"] = "/manage/" + discordgo.StrID(activeGuild.ID) + "/genai/commands"

	return templateData, nil
}

// Command handlers
func CommandMiddleware(inner func(w http.ResponseWriter, r *http.Request, override *models.GenaiCommand) (web.TemplateData, error)) web.ControllerHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
		activeGuild := r.Context().Value(common.ContextKeyCurrentGuild).(*dstate.GuildSet)

		var command *models.GenaiCommand
		var err error

		id := pat.Param(r, "commandID")
		idParsed, _ := strconv.ParseInt(id, 10, 64)
		command, err = models.GenaiCommands(qm.Where("guild_id = ? AND id = ?", activeGuild.ID, idParsed)).OneG(r.Context())

		if err != nil {
			return nil, web.NewPublicError("Command not found, someone else deleted it in the meantime perhaps? Check control panel logs")
		}

		tmpl, err := inner(w, r, command)
		featureflags.MarkGuildDirty(activeGuild.ID)
		return tmpl, err
	}
}

func HandleCreateCommand(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	activeGuild, templateData := web.GetBaseCPContextData(r.Context())
	formData := r.Context().Value(common.ContextKeyParsedForm).(*CommandFormData)

	count, err := models.GenaiCommands(models.GenaiCommandWhere.GuildID.EQ(activeGuild.ID)).CountG(r.Context())
	if err != nil {
		return templateData, errors.WithMessage(err, "count")
	}

	if count > 4 {
		return templateData.AddAlerts(web.ErrorAlert("Max 5 GenAI commands supported.")), nil
	}

	localID, err := common.GenLocalIncrID(activeGuild.ID, "genai_command")
	if err != nil {
		return templateData, errors.WrapIf(err, "error generating local id")
	}

	triggers := []string{formData.Name}
	if formData.Aliases != "" {
		triggers = append(triggers, strings.Split(formData.Aliases, ",")...)
	}

	model := &models.GenaiCommand{
		ID:                      localID,
		GuildID:                 activeGuild.ID,
		Enabled:                 true,
		Triggers:                triggers,
		Prompt:                  formData.Prompt,
		AllowInput:              formData.AllowInput,
		MaxTokens:               formData.MaxTokens,
		AutodeleteResponse:      formData.AutodeleteResponse,
		AutodeleteTrigger:       formData.AutodeleteTrigger,
		AutodeleteResponseDelay: formData.AutodeleteResponseDelay,
		AutodeleteTriggerDelay:  formData.AutodeleteTriggerDelay,
		Channels:                append(formData.Channels, formData.Categories...),
		ChannelsWhitelistMode:   formData.ChannelsWhitelistMode,
		Roles:                   formData.Roles,
		RolesWhitelistMode:      formData.RolesWhitelistMode,
	}

	err = model.InsertG(r.Context(), boil.Infer())
	if err == nil {
		err = featureflags.UpdatePluginFeatureFlags(activeGuild.ID, &Plugin{})
		if err != nil {
			web.CtxLogger(r.Context()).WithError(err).Error("failed updating feature flags")
		}
		//	go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyNewChannelOverride))
	}
	return templateData, errors.WithMessage(err, "InsertG")
}

func HandleUpdateCommand(w http.ResponseWriter, r *http.Request, currentCommand *models.GenaiCommand) (web.TemplateData, error) {
	_, templateData := web.GetBaseCPContextData(r.Context())

	formData := r.Context().Value(common.ContextKeyParsedForm).(*CommandFormData)

	triggers := []string{formData.Name}
	if formData.Aliases != "" {
		triggers = append(triggers, strings.Split(formData.Aliases, ",")...)
	}

	currentCommand.Enabled = formData.Enabled
	currentCommand.Triggers = triggers
	currentCommand.Prompt = formData.Prompt
	currentCommand.AllowInput = formData.AllowInput
	currentCommand.MaxTokens = formData.MaxTokens
	currentCommand.AutodeleteResponse = formData.AutodeleteResponse
	currentCommand.AutodeleteTrigger = formData.AutodeleteTrigger
	currentCommand.AutodeleteResponseDelay = formData.AutodeleteResponseDelay
	currentCommand.AutodeleteTriggerDelay = formData.AutodeleteTriggerDelay
	currentCommand.Channels = append(formData.Channels, formData.Categories...)
	currentCommand.ChannelsWhitelistMode = formData.ChannelsWhitelistMode
	currentCommand.Roles = formData.Roles
	currentCommand.RolesWhitelistMode = formData.RolesWhitelistMode

	_, err := currentCommand.UpdateG(r.Context(), boil.Infer())
	//	if err == nil {
	//		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyUpdatedChannelOverride))
	//	}
	return templateData, errors.WithMessage(err, "UpdateG")
}

func HandleDeleteCommand(w http.ResponseWriter, r *http.Request, currentOverride *models.GenaiCommand) (web.TemplateData, error) {
	activeGuild, templateData := web.GetBaseCPContextData(r.Context())

	_, err := currentOverride.DeleteG(r.Context())
	//	if rows > 0 {
	//		go cplogs.RetryAddEntry(web.NewLogEntryFromContext(r.Context(), panelLogKeyRemovedChannelOverride))
	//	}

	ffErr := featureflags.UpdatePluginFeatureFlags(activeGuild.ID, &Plugin{})
	if ffErr != nil {
		web.CtxLogger(r.Context()).WithError(ffErr).Error("failed updating feature flags")
	}

	return templateData, errors.WithMessage(err, "DeleteG")
}

var _ web.PluginWithServerHomeWidget = (*Plugin)(nil)

func (p *Plugin) LoadServerHomeWidget(w http.ResponseWriter, r *http.Request) (web.TemplateData, error) {
	ag, templateData := web.GetBaseCPContextData(r.Context())

	templateData["WidgetTitle"] = "Generative AI"
	templateData["SettingsPath"] = "/genai"

	config, err := GetConfig(ag.ID)
	if err != nil {
		return templateData, err
	}

	format := `<ul>
	<li>GenAI status: %s</li>
	<li>GenAI provider: <code>%s</code></li>
	<li>GenAI model: <code>#%s</code></li>
	<li>GenAI base command: <code>#%s</code>%s</li>
</ul>`

	status := web.EnabledDisabledSpanStatus(config.Enabled)

	if config.Enabled {
		templateData["WidgetEnabled"] = true
	} else {
		templateData["WidgetDisabled"] = true
	}

	provider := GenAIProviderFromID(config.Provider)
	baseCmdStr := "disabled"
	if config.BaseCMDEnabled {
		baseCmdStr = "enabled"
	}
	templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, status, provider.String(), config.Model, baseCmdStr, web.Indicator(config.BaseCMDEnabled)))

	return templateData, nil
}
