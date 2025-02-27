package globalrl

import (
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/common"
	"github.com/SoggySaussages/syzygy/lib/dcmd"
	"github.com/SoggySaussages/syzygy/lib/discordgo"
	"github.com/SoggySaussages/syzygy/stdcommands/util"
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
