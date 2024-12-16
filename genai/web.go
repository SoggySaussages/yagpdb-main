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

const (
	ConextKeyConfig ConextKey = iota
)

var panelLogKey = cplogs.RegisterActionFormat(&cplogs.ActionFormat{Key: "genai_settings_updated", FormatString: "Updated genai settings"})

func (p *Plugin) InitWeb() {
	web.AddHTMLTemplate("genai/assets/genai.html", PageHTML)
	web.AddSidebarItem(web.SidebarCategoryGenAI, &web.SidebarItem{
		Name: "General",
		URL:  "genai",
		Icon: "fas fa-video",
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

	genaiMux.Handle(pat.Post(""), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), Config{}))
	genaiMux.Handle(pat.Post("/"), web.FormParserMW(web.RenderHandler(HandlePostGenAI, "cp_genai"), Config{}))
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
	newConf := ctx.Value(common.ContextKeyParsedForm).(*Config)
	newConfNoKey := *newConf
	newConfNoKey.Key = ""

	tmpl["GenAIConfig"] = newConfNoKey

	if !ok {
		return tmpl
	}

	var err error
	provider := GenAIProviderFromID(newConf.Provider)
	if provider.KeyRequired() {
		newConf.Key, err = encryptAPIToken(&dstate.GuildState{ID: guild.ID, OwnerID: guild.OwnerID}, newConf.Key)
		if web.CheckErr(tmpl, err, "Failed encrypting your API token to save", web.CtxLogger(ctx).Error) {
			return tmpl
		}
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

	templateData["WidgetTitle"] = "GenAI"
	templateData["SettingsPath"] = "/genai"

	config, err := GetConfig(ag.ID)
	if err != nil {
		return templateData, err
	}

	format := `<ul>
	<li>GenAI status: %s</li>
	<li>GenAI provider: <code>%s</code>%s</li>
	<li>GenAI model: <code>#%s</code>%s</li>
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
	templateData["WidgetBody"] = template.HTML(fmt.Sprintf(format, status, provider.String(), web.Indicator(true), config.Model, web.Indicator(true), baseCmdStr, web.Indicator(config.BaseCmdEnabled)))

	return templateData, nil
}
