package leaveserver

import (
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "leaveserver",
	Description:          "Causes SGPDB to leave the specified server. The bot may still be invited back with full functionality restored. Bot Owner Only",
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
