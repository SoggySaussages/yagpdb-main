package info

import (
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
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
