package genai

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/common/cplogs"
	"github.com/botlabs-gg/yagpdb/v2/common/featureflags"
	"github.com/botlabs-gg/yagpdb/v2/genai/models"
	"github.com/botlabs-gg/yagpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"github.com/botlabs-gg/yagpdb/v2/web"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"goji.io"
	"goji.io/pat"
)

//go:embed assets/genai.html
var PageHTML string

type ConextKey int

type FormData struct {
	Enabled        bool   `json:"enabled" schema:"enabled"`
	Provider       int    `json:"provider" schema:"provider"`
	Model          string `json:"model" schema:"model"`
	Key            string `json:"key" schema:"key"`
	BaseCmdEnabled bool   `json:"base_cmd_enabled" schema:"base_cmd_enabled"`
	ResetToken     bool   `json:"reset_token" schema:"reset_token"`
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

	// Alll handlers here require guild channels present
	genaiMux.Use(web.RequireServerAdminMiddleware)
	genaiMux.Use(baseData)

	// Get just renders the template, so let the renderhandler do all the work
	genaiMux.Handle(pat.Get(""), web.RenderHandler(nil, "cp_genai"))
	genaiMux.Handle(pat.Get("/"), web.RenderHandler(nil, "cp_genai"))

	genaiMux.Handle(pat.Post(""), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), FormData{}))
	genaiMux.Handle(pat.Post("/"), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), FormData{}))
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
	newConfFakeKey := *newConf

	tmpl["GenAIConfig"] = &newConfFakeKey

	if !ok {
		return tmpl
	}

	var saveNewKey bool
	provider := GenAIProviderFromID(newConf.Provider)
	conf, err := GetConfig(guild.ID)
	if err != nil {
		conf = &models.GenaiConfig{}
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
			if web.CheckErr(tmpl, provider.ValidateAPIToken(partialGuildState, string(newConf.Key)), "Your API token is invalid.", web.CtxLogger(ctx).Error) {
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
