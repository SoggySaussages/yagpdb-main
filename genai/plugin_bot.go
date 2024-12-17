package genai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strconv"

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
}

func GenAIProviderFromID(id GenAIProviderID) GenAIProvider {
	for _, p := range GenAIProviders {
		if p.ID() == id {
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
		&dcmd.ArgDef{Name: "Prompt", Type: dcmd.String}},
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
		if !config.BaseCmdEnabled {
			return "", commands.NewUserErrorf("The **genai** command is disabled on this server. It can be enabled at <%s>", genaiConfigPage)
		}

		provider := GenAIProviderFromID(config.Provider)
		if provider.KeyRequired() && config.Key == "" {
			return "", commands.NewUserErrorf("No API key set for %s. It can be enabled at <%s>", provider.String(), genaiConfigPage)
		}

		response, _, err := provider.BasicCompletion(&data.GuildData.GS.GuildState, "", data.Args[0].Str(), 512, false)
		if err != nil {
			return "", err
		}
		return response.Content, nil
	},
}

func createKey(gs *dstate.GuildState) ([]byte, error) {
	//	if gs.OwnerID == 0 {
	//		gs.OwnerID = bot.State.GetGuild(gs.ID).OwnerID
	//	}
	logger.Infof("%d, %d, %s.", gs.ID, gs.OwnerID, common.GetBotToken())
	salt := []byte(strconv.FormatInt(gs.ID+gs.OwnerID, 10))
	return scrypt.Key([]byte(common.GetBotToken()), salt, 16384, 8, 1, 32)
}

func encryptAPIToken(gs *dstate.GuildState, token string) (string, error) {
	key, err := createKey(gs)
	if err != nil {
		return "", err
	}
	logger.Info(key)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	logger.Info(string(nonce))

	cypheredToken := gcm.Seal(nonce, nonce, []byte(token), nil)

	return string(cypheredToken), nil
}

func decryptAPIToken(gs *dstate.GuildState, encryptedToken string) (string, error) {
	key, err := createKey(gs)
	if err != nil {
		return "", err
	}
	logger.Info(key)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	encryptedTokenBytes := []byte(encryptedToken)
	nonce, encryptedTokenBytes := encryptedTokenBytes[:gcm.NonceSize()], encryptedTokenBytes[gcm.NonceSize():]
	logger.Info(string(nonce))

	decryptedToken, err := gcm.Open(nil, nonce, encryptedTokenBytes, nil)
	logger.Info(string(decryptedToken))
	if err != nil {
		logger.WithError(err).Error("failed decrypting a genai API token")
		return "", ErrorAPIKeyInvalid
	}

	return string(decryptedToken), nil
}

func getAPIToken(gs *dstate.GuildState) (string, error) {
	config, err := GetConfig(gs.ID)
	if err != nil {
		logger.WithError(err).WithField("guild", gs.ID).Error("Failed retrieving openai config")
		return "", err
	}

	if !config.Enabled {
		return "", nil
	}

	return decryptAPIToken(gs, config.Key)
}
