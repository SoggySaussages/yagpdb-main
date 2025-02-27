package leaveserver

import (
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/common"
	"github.com/SoggySaussages/syzygy/lib/dcmd"
	"github.com/SoggySaussages/syzygy/stdcommands/util"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "leaveserver",
	Description:          "Causes SYZYGY to leave the specified server. The bot may still be invited back with full functionality restored. Bot Owner Only",
	HideFromHelp:         true,
	RequiredArgs:         1,
	Arguments: []*dcmd.ArgDef{
		{Name: "server", Type: dcmd.BigInt},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		err := common.BotSession.GuildLeave(data.Args[0].Int64())
		if err == nil {
			return "Left " + data.Args[0].Str(), nil
		}
		return err, err
	}),
}
