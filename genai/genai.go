package genai

import (
	"emperror.dev/errors"
	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/common/featureflags"
	"github.com/botlabs-gg/yagpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
)

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "General Generative AI",
		SysName:  "genai",
		Category: common.PluginCategoryMisc,
	}
}

var logger = common.GetPluginLogger(&Plugin{})

func RegisterPlugin() {
	plugin := &Plugin{}
	common.RegisterPlugin(plugin)
}

var _ featureflags.PluginWithFeatureFlags = (*Plugin)(nil)

const (
	featureFlagEnabled         = "genai_enabled"
	featureFlagCommandsEnabled = "genai_commands_enabled"
)

func (p *Plugin) UpdateFeatureFlags(guildID int64) ([]string, error) {
	config, err := GetConfig(guildID)
	if err != nil {
		return nil, errors.WithStackIf(err)
	}

	var flags []string
	if config.Enabled && config.Key != "" {
		flags = append(flags, featureFlagEnabled)
	}

	return flags, nil
}

func (p *Plugin) AllFeatureFlags() []string {
	return []string{
		featureFlagEnabled,         // set if this server uses genai
		featureFlagCommandsEnabled, // set if this server uses simple genai commands
	}
}

var ErrorAPIKeyInvalid = errors.New("Your Generative AI API token has been invalidated due to a change in security (server owner change, bot token reset, etc.) Please reset your API token.")

type GenAIProviderID uint

const (
	GenAIProviderOpenAIID GenAIProviderID = iota
)

const CountSupportedProviders = 1

type GenAIProviderModelMap map[string]string

type GenAIFunctionDefinition struct {
	Name       string
	Definition string
	Arguments  map[string]string
}

type GenAIFunctionResponse struct {
	Name      string
	Arguments map[string]interface{}
}

type GenAIInput struct {
	SystemMessage string
	UserMessage   string
	Functions     *[]GenAIFunctionDefinition
}

type GenAIResponse struct {
	Content   string
	Functions *[]GenAIFunctionResponse
}

type GenAIResponseUsage struct {
	InputTokens  uint64
	OutputTokens uint64
}

type GenAIModerationCategoryProbability map[string]float64

var GenAIModerationCategories = []string{
	"harassment",
	"harassment/threatening",
	"hate",
	"hate/threatening",
	"illicit",
	"illicit/violent",
	"self-harm",
	"self-harm/intent",
	"self-harm/instructions",
	"sexual",
	"sexual/minors",
	"violence",
	"violence/graphic",
}

type GenAIProvider interface {
	ID() GenAIProviderID
	DefaultModel() string
	ModelMap() *GenAIProviderModelMap

	BasicCompletion(gs *dstate.GuildState, systemMsg, userMsg string) (*GenAIResponse, *GenAIResponseUsage, error)
	ComplexCompletion(gs *dstate.GuildState, input *GenAIInput) (*GenAIResponse, *GenAIResponseUsage, error)
	ModerateMessage(gs *dstate.GuildState, message string) (*[]GenAIModerationCategoryProbability, *GenAIResponseUsage, error)
}

var GenAIProviders = make([]*GenAIProvider, CountSupportedProviders)

type Config struct {
	Enabled bool `json:"enabled" schema:"enabled"` // Wether genai is enabled or not

	// The selected provider of genai (such as openai)
	Provider GenAIProviderID `json:"provider" schema:"provider"`

	// The selected model from the provider (such as gpt-4)
	Model string `json:"model" schema:"model"`

	// The encrypted key for the selected API
	Key string `json:"key" schema:"key"`
}

func (c *Config) Save(guildID int64) error {
	return common.SetRedisJson("genai_config:"+discordgo.StrID(guildID), c)
}

var DefaultConfig = &Config{
	Enabled:  false,
	Provider: GenAIProviderOpenAI{}.Type(),
	Model:    GenAIProviderOpenAI{}.DefaultModel(),
}

// Returns he guild's conifg, or the defaul one if not set
func GetConfig(guildID int64) (*Config, error) {
	var config *Config
	err := common.GetRedisJson("genai_config:"+discordgo.StrID(guildID), &config)
	if err == nil && config == nil {
		return DefaultConfig, nil
	}

	return config, err
}
