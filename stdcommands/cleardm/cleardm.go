package cleardm

import (
	"fmt"

	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/lib/discordgo"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
	"github.com/sirupsen/logrus"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "cleardm",
	Description:          "clears the DM chat with a user, bot owner only command.",
	HideFromHelp:         true,
	RequiredArgs:         1,
	Arguments: []*dcmd.ArgDef{
		{Name: "Target", Type: dcmd.User},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		var target = data.Args[0].Value.(*discordgo.User)
		dm, err := common.BotSession.UserChannelCreate(target.ID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to get DM channel for user %d", target.ID)
			return nil, err
		}
		messages, err := common.BotSession.ChannelMessages(dm.ID, 100, 0, 0, 0)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to get DM messages for user %d", target.ID)
			return nil, err
		}
		if len(messages) == 0 {
			return "No messages found", nil
		}
		count := 0
		for _, message := range messages {
			if message.Author.ID == common.BotUser.ID {
				err = common.BotSession.ChannelMessageDelete(dm.ID, message.ID)
				if err != nil {
					logrus.WithError(err).Errorf("Failed to delete DM messages for user %d", target.ID)
				} else {
					count++
				}
			}
		}
		return fmt.Sprintf("Deleted %d messages for %s in DMs", count, target.Mention()), nil
	}),
}
