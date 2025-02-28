package dcallvoice

import (
	"fmt"

	"github.com/SoggySaussages/sgpdb/bot"
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/lib/discordgo"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "dcallvoice",
	Description:          "Disconnects from all the voice channels the bot is in. Bot Owner Only",
	HideFromHelp:         true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {

		vcs := make([]*discordgo.VoiceState, 0)

		processShards := bot.ReadyTracker.GetProcessShards()
		for _, shard := range processShards {
			guilds := bot.State.GetShardGuilds(int64(shard))
			for _, g := range guilds {
				vc := g.GetVoiceState(common.BotUser.ID)
				if vc != nil {
					vcs = append(vcs, vc)
					go bot.ShardManager.SessionForGuild(g.ID).GatewayManager.ChannelVoiceLeave(g.ID)
				}
			}
		}

		return fmt.Sprintf("Leaving %d voice channels...", len(vcs)), nil
	}),
}
