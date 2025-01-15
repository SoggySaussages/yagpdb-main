package genai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/botlabs-gg/sgpdb/v2/automod"
	"github.com/botlabs-gg/sgpdb/v2/bot"
	"github.com/botlabs-gg/sgpdb/v2/bot/eventsystem"
	"github.com/botlabs-gg/sgpdb/v2/commands"
	"github.com/botlabs-gg/sgpdb/v2/common"
	"github.com/botlabs-gg/sgpdb/v2/genai/models"
	"github.com/botlabs-gg/sgpdb/v2/lib/dcmd"
	"github.com/botlabs-gg/sgpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/sgpdb/v2/lib/dstate"
	"github.com/botlabs-gg/sgpdb/v2/web"
	"golang.org/x/crypto/scrypt"
)

var _ commands.CommandProvider = (*Plugin)(nil)
var _ bot.BotInitHandler = (*Plugin)(nil)

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p, baseCmd)
}

func (p *Plugin) BotInit() {
	generateFormattedModCategoryList()
	genCustomModerateFuncArgs()

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

	eventsystem.AddHandlerAsyncLastLegacy(p, bot.ConcurrentEventHandler(HandleMessageCreate), eventsystem.EventMessageCreate)
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

func HandleMessageCreate(evt *eventsystem.EventData) {
	mc := evt.MessageCreate()
	cs := evt.CSOrThread()

	if !evt.HasFeatureFlag(featureFlagCommandsEnabled) || !evt.HasFeatureFlag(featureFlagEnabled) {
		return
	}

	if mc.GuildID == 0 {
		return
	}

	if cs == nil {
		logger.Warn("Channel not found in state")
		return
	}

	if !bot.IsUserMessage(mc.Message) {
		return
	}

	if mc.Author.Bot {
		return
	}

	if hasPerms, _ := bot.BotHasPermissionGS(evt.GS, cs.ID, discordgo.PermissionSendMessages); !hasPerms {
		return
	}

	member := dstate.MemberStateFromMember(mc.Member)
	member.GuildID = evt.GS.ID

	cmds, err := models.GenaiCommands(
		models.GenaiCommandWhere.GuildID.EQ(evt.GS.ID)).AllG(evt.Context())
	if err != nil {
		logger.Error(err)
		return
	}

	prefix, err := commands.GetCommandPrefixBotEvt(evt)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, cmd := range cmds {
		var cmdRunsInChannel bool
		var found bool
		for _, v := range cmd.Channels {
			if v == cs.ID || v == cs.ParentID {
				cmdRunsInChannel = cmd.ChannelsWhitelistMode
				found = true
				break
			}
		}
		if !found {
			cmdRunsInChannel = !cmd.ChannelsWhitelistMode
		}

		var cmdRunsForUser bool
		found = false
		if len(cmd.Roles) == 0 {
			cmdRunsForUser = !cmd.RolesWhitelistMode
			found = true
		}
		if !found {
			for _, v := range cmd.Roles {
				if common.ContainsInt64Slice(member.Member.Roles, v) {
					cmdRunsForUser = cmd.RolesWhitelistMode
					found = true
					break
				}
			}
		}
		if !found {
			cmdRunsForUser = !cmd.RolesWhitelistMode
		}

		if !cmd.Enabled || !cmdRunsInChannel || !cmdRunsForUser {
			continue
		}

		var triggersSafe []string
		for _, t := range cmd.Triggers {
			triggersSafe = append(triggersSafe, regexp.QuoteMeta(strings.TrimSpace(t)))
		}

		pattern := `(?mi)\A(<@!?` + discordgo.StrID(common.BotUser.ID) + "> ?|" + regexp.QuoteMeta(prefix) + ")(" + strings.Join(triggersSafe, "|") + `)(\z|[[:space:]])`
		re, err := regexp.Compile(pattern)
		if err != nil {
			logger.Error(err)
			continue
		}

		content := mc.Content
		idx := re.FindStringIndex(content)
		if idx == nil {
			continue
		}

		var userMsg string
		if cmd.AllowInput {
			userMember, _ := bot.GetMember(mc.GuildID, mc.Author.ID)
			userName := userMember.Member.Nick
			if userName == "" {
				userName = mc.Author.String()
			}
			userMsg = userName + ": " + content[idx[1]:]
			if mc.ReferencedMessage != nil {
				authorMember, _ := bot.GetMember(mc.GuildID, mc.ReferencedMessage.Author.ID)
				authorName := mc.ReferencedMessage.Author.String() + "'s"
				if authorMember != nil && authorMember.Member.Nick != "" {
					authorName = authorMember.Member.Nick + "'s"
				}
				if mc.ReferencedMessage.Author.ID == common.BotUser.ID {
					authorName = "your"
				}
				refMessageContent := mc.ReferencedMessage.ContentWithMentionsReplaced()
				if len(mc.ReferencedMessage.Embeds) > 0 {
					emb := mc.ReferencedMessage.Embeds[0]
					refMessageContent = emb.Description
					if emb.Author != nil && emb.Author.Name != "" {
						authorName += " (as " + emb.Author.Name + ")"
					}
				}
				userMsg = fmt.Sprintf("%s\n\nPlease note this message is replying to %s previous message: %s", userMsg, authorName, refMessageContent)
			}
		}

		config, err := GetConfig(evt.GS.ID)
		if err != nil {
			logger.Error(err)
			return
		}
		provider := GenAIProviderFromID(config.Provider)

		r, _, err := provider.BasicCompletion(&evt.GS.GuildState, cmd.Prompt, userMsg, int64(cmd.MaxTokens), cs.NSFW || evt.GS.GetChannel(cs.ParentID).NSFW)
		if err != nil {
			r = &GenAIResponse{Content: fmt.Sprintf("Failed executing your custom GenAI command; %v responded with: %s", provider, err.Error())}
		}
		m, err := common.BotSession.ChannelMessageSendReply(cs.ID, r.Content, &discordgo.MessageReference{ChannelID: cs.ID, MessageID: mc.ID})
		if err == nil && (cmd.AutodeleteTrigger || cmd.AutodeleteResponse) {
			if cmd.AutodeleteTrigger {
				go func() {
					time.Sleep(time.Duration(cmd.AutodeleteTriggerDelay) * time.Second)
					common.BotSession.ChannelMessageDelete(cs.ID, mc.ID)
				}()
			}
			if cmd.AutodeleteResponse {
				go func() {
					time.Sleep(time.Duration(cmd.AutodeleteResponseDelay) * time.Second)
					common.BotSession.ChannelMessageDelete(cs.ID, m.ID)
				}()
			}
		}
		if err != nil {
			logger.Error(err)
		}
	}
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
