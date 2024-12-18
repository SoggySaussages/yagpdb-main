package genai

import (
	"strings"

	"github.com/botlabs-gg/yagpdb/v2/automod"
	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
)

type GenAIAutomodTriggerData struct {
	Threshold  int
	Categories []string
}

var _ automod.MessageTrigger = (*GenAIAutomodTrigger)(nil)

type GenAIAutomodTrigger struct{}

func (mc *GenAIAutomodTrigger) Kind() automod.RulePartType {
	return automod.RulePartTrigger
}

func (mc *GenAIAutomodTrigger) DataType() interface{} {
	return &GenAIAutomodTriggerData{}
}

func (mc *GenAIAutomodTrigger) Name() string {
	return "Generative AI trigger"
}

func (mc *GenAIAutomodTrigger) Description() string {
	return "Triggers on messages AI marks abusive. Requires GenAI to be enabled."
}

func (mc *GenAIAutomodTrigger) UserSettings() []*automod.SettingDef {
	settings := []*automod.SettingDef{
		{
			Name:    "Certainty Threshold (%)",
			Key:     "Threshold",
			Kind:    automod.SettingTypeInt,
			Default: 60,
			Min:     1,
			Max:     99,
		},
		{
			Name:    "Max Tokens per Message",
			Key:     "Threshold",
			Kind:    automod.SettingTypeInt,
			Default: 512,
			Min:     1,
			Max:     512,
		},
	}
	catSettings := &automod.SettingDef{
		Name: "Categories to Trigger",
		Key:  "Categories",
		Kind: automod.SettingTypeMultiOptionsCustom,
	}
	for i, s := range GenAIModerationCategories {
		catSettings.Options = append(catSettings.Options, automod.SettingTypeOptionsCustomOption{
			Name:  GenAIModerationCategoriesFormatted[i],
			Value: strings.ReplaceAll(s, " ", "-")})
	}
	return append(settings, catSettings)
}

func (mc *GenAIAutomodTrigger) CheckMessage(triggerCtx *automod.TriggerContext, cs *dstate.ChannelState, m *discordgo.Message) (bool, error) {
	dataCast := triggerCtx.Data.(*GenAIAutomodTriggerData)
	config, err := GetConfig(cs.GuildID)
	if err != nil {
		return false, err
	}

	g, err := common.BotSession.Guild(cs.GuildID)
	if err != nil {
		return false, err
	}

	categories, _, err := GenAIProviderFromID(config.Provider).ModerateMessage(dstate.GuildStateFromDgo(g), m.Content)
	for _, c := range dataCast.Categories {
		if (*categories)[c]*100 > float64(dataCast.Threshold) {
			return true, nil
		}
	}
	return false, nil
}
