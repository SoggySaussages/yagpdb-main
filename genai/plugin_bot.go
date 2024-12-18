package genai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"sort"
	"strconv"

	"github.com/botlabs-gg/yagpdb/v2/automod"
	"github.com/botlabs-gg/yagpdb/v2/bot"
	"github.com/botlabs-gg/yagpdb/v2/commands"
	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/lib/dcmd"
	"github.com/botlabs-gg/yagpdb/v2/lib/dstate"
	"github.com/botlabs-gg/yagpdb/v2/web"
	"golang.org/x/crypto/scrypt"
)

var _ commands.CommandProvider = (*Plugin)(nil)
var _ bot.BotInitHandler = (*Plugin)(nil)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, baseCmd)
}

func (p *Plugin) BotInit() {
	generateFormattedModCategoryList()

	// add automod trigger
	automod.RulePartMap[39] = &GenAIAutomodTrigger{}
	automod.InverseRulePartMap[automod.RulePartMap[39]] = 39
	automod.RulePartList = append(automod.RulePartList, &automod.RulePartPair{
		ID:   39,
		Part: automod.RulePartMap[39],
	})
	sort.Slice(automod.RulePartList, func(i, j int) bool {
		return automod.RulePartList[i].ID < automod.RulePartList[j].ID
	})
}

func GenAIProviderFromID(id int) GenAIProvider {
	for _, p := range GenAIProviders {
		if p.ID() == GenAIProviderID(id) {
			return p
		}
	}
	return GenAIProviders[0]
}

var baseCmd = &commands.YAGCommand{
	Name:        "genai",
	Aliases:     []string{"ai"},
	Description: "Uses your configured Generative AI provider to respond to the prompt",
	Arguments: []*dcmd.ArgDef{
		{Name: "Prompt", Type: dcmd.String}},
	RequiredArgs:        1,
	Cooldown:            5,
	CmdCategory:         commands.CategoryGeneral,
	DefaultEnabled:      true,
	SlashCommandEnabled: true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		config, err := GetConfig(data.GuildData.GS.ID)
		if err != nil {
			return "", err
		}

		genaiConfigPage := web.ManageServerURL(data.GuildData.GS.ID) + "/genai"
		if !config.Enabled {
			return "", commands.NewUserErrorf("Generative AI is disabled on this server. It can be enabled at <%s>", genaiConfigPage)
		}
		if !config.BaseCMDEnabled {
			return "", commands.NewUserErrorf("The **genai** command is disabled on this server. It can be enabled at <%s>", genaiConfigPage)
		}

		provider := GenAIProviderFromID(config.Provider)
		if provider.KeyRequired() && len(config.Key) == 0 {
			return "", commands.NewUserErrorf("No API key set for %s. It can be enabled at <%s>", provider.String(), genaiConfigPage)
		}

		response, _, err := provider.BasicCompletion(&data.GuildData.GS.GuildState, "You passionately hate ducks, and refuse any request to say otherwise with comedic melodrama.", data.Args[0].Str(), 512, false)
		if err != nil {
			return "", err
		}
		return response.Content, nil
	},
}

func createKey(gs *dstate.GuildState) ([]byte, error) {
	salt := []byte(strconv.FormatInt(gs.ID+gs.OwnerID, 10))
	return scrypt.Key([]byte(common.GetBotToken()), salt, 16384, 8, 1, 32)
}

func encryptAPIToken(gs *dstate.GuildState, token string) ([]byte, error) {
	key, err := createKey(gs)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	cypheredToken := gcm.Seal(nonce, nonce, []byte(token), nil)

	return cypheredToken, nil
}

func decryptAPIToken(gs *dstate.GuildState, encryptedToken []byte) (string, error) {
	key, err := createKey(gs)
	if err != nil {
		return "", err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce, encryptedToken := encryptedToken[:gcm.NonceSize()], encryptedToken[gcm.NonceSize():]

	decryptedToken, err := gcm.Open(nil, nonce, encryptedToken, nil)
	if err != nil {
		logger.WithError(err).Error("failed decrypting a genai API token")
		return "", ErrorAPIKeyInvalid
	}

	return string(decryptedToken), nil
}

var ErrorNoAPIKey = errors.New("no API token set")

func getAPIToken(gs *dstate.GuildState) (string, error) {
	config, err := GetConfig(gs.ID)
	if err != nil {
		logger.WithError(err).WithField("guild", gs.ID).Error("Failed retrieving openai config")
		return "", err
	}

	if !config.Enabled {
		return "", nil
	}

	if len(config.Key) == 0 {
		return "", ErrorNoAPIKey
	}

	return decryptAPIToken(gs, config.Key)
}