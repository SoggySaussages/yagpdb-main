package genai

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/common/cplogs"
	"github.com/botlabs-gg/yagpdb/v2/common/featureflags"
	"github.com/botlabs-gg/yagpdb/v2/common/pubsub"
	"github.com/botlabs-gg/yagpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"github.com/botlabs-gg/yagpdb/v2/web"
	"goji.io"
	"goji.io/pat"
)

//go:embed assets/genai.html
var PageHTML string

type ConextKey int

type FormData struct {
	Enabled        bool   `json:"enabled" schema:"enabled"`
	Provider       uint   `json:"provider" schema:"provider"`
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
		if config.Key != "" {
			config.Key = "key-hidden-for-security"
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
	newConf := &Config{
		Enabled:        formData.Enabled,
		Provider:       GenAIProviderID(formData.Provider),
		Model:          formData.Model,
		BaseCmdEnabled: formData.BaseCmdEnabled,
	}
	newConfFakeKey := *newConf
	if formData.Key != "" {
		newConfFakeKey.Key = "key-hidden-for-security"
	}

	tmpl["GenAIConfig"] = &newConfFakeKey

	if !ok {
		return tmpl
	}

	var saveNewKey bool
	provider := GenAIProviderFromID(newConf.Provider)
	conf, err := GetConfig(guild.ID)
	if err != nil {
		conf = &Config{}
	}

	if formData.Key != "" {
		newConf.Key, err = encryptAPIToken(&dstate.GuildState{ID: guild.ID, OwnerID: guild.OwnerID}, formData.Key)
		if web.CheckErr(tmpl, err, "Failed encrypting your API token to save", web.CtxLogger(ctx).Error) {
			return tmpl
		}
	}

	keyChanged := newConf.Key != conf.Key && newConf.Key != ""
	saveNewKey = provider.KeyRequired() && (formData.ResetToken || keyChanged)
	if !saveNewKey {
		newConf.Key = conf.Key
	}

	err = newConf.Save(guild.ID)
	if web.CheckErr(tmpl, err, "Failed saving config :'(", web.CtxLogger(ctx).Error) {
		return tmpl
	}

	err = featureflags.UpdatePluginFeatureFlags(guild.ID, &Plugin{})
	if err != nil {
		web.CtxLogger(ctx).WithError(err).Error("failed updating feature flags")
	}

	err = pubsub.Publish("update_genai", guild.ID, nil)
	if err != nil {
		web.CtxLogger(ctx).WithError(err).Error("Failed sending update genai event")
	}

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
	if config.BaseCmdEnabled {
		baseCmdStr = "enabled"
	}
	templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, status, provider.String(), config.Model, baseCmdStr, web.Indicator(config.BaseCmdEnabled)))

	return templateData, nil
}
