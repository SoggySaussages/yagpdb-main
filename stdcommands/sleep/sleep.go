package sleep

import (
	"time"

	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/lib/dcmd"
	"github.com/SoggySaussages/syzygy/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "sleep",
	Description:          "Maintenance command, used to test command queueing. Bot Owner Only",
	HideFromHelp:         true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		time.Sleep(time.Second * 5)
		return "Slept, Done", nil
	}),
}
