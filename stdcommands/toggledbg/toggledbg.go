package toggledbg

import (
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/sirupsen/logrus"

	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "toggledbg",
	Description:          "Toggles Debug Logging. Restarting the bot will always reset debug logging. Bot Owner Only",
	HideFromHelp:         true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			common.SetLoggingLevel(logrus.InfoLevel)
			return "Disabled debug logging", nil
		}

		common.SetLoggingLevel(logrus.DebugLevel)
		return "Enabled debug logging", nil

	}),
}
