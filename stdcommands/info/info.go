package info

import (
	"github.com/botlabs-gg/sgpdb/v2/commands"
	"github.com/botlabs-gg/sgpdb/v2/lib/dcmd"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Info",
	Description: "Responds with bot information",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		info := `SGPDB - A **S**oggy **G**eneral **P**urpose **D**iscord **B**ot
This is a fork of [SGPDB](https://github.com/botlabs-gg/sgpdb) hosted by SoggySaussages.`

		return info, nil
	},
}
