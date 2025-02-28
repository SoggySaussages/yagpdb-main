package globalrl

import (
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/lib/discordgo"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	Name:                 "globalrl",
	Description:          "Tests the global ratelimit functionality. Bot Owner Only",
	RequiredArgs:         1,
	HideFromHelp:         true,
	HideFromCommandsPage: true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {

		rlEvt := &discordgo.RateLimit{
			URL: "Wew",
			TooManyRequests: &discordgo.TooManyRequests{
				Bucket:     "wewsss",
				Message:    "Too many!",
				RetryAfter: 5,
			},
		}

		go common.BotSession.HandleEvent("__RATE_LIMIT__", rlEvt)

		return "Done", nil
	}),
}
