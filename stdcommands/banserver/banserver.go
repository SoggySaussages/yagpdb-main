package banserver

import (
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
	"github.com/mediocregopher/radix/v3"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "banserver",
	Description:          "Bans the specified server from using the bot. SGPDB will leave the server, and leave whenever invited back. Bot Owner Only",
	HideFromHelp:         true,
	RequiredArgs:         1,
	Arguments: []*dcmd.ArgDef{
		{Name: "server", Type: dcmd.BigInt},
	},
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		err := common.BotSession.GuildLeave(data.Args[0].Int64())
		if err == nil {

			common.RedisPool.Do(radix.FlatCmd(nil, "SADD", "banned_servers", data.Args[0].Int64()))

			return "Banned " + data.Args[0].Str(), nil
		}
		return nil, err
	}),
}
