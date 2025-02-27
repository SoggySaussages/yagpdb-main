package info

import (
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/lib/dcmd"
)

var Command = &commands.YAGCommand{
	CmdCategory: commands.CategoryGeneral,
	Name:        "Info",
	Description: "Responds with bot information",
	RunInDM:     true,
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		info := `SYZYGY - A **S**oggy **G**eneral **P**urpose **D**iscord **B**ot
This is a fork of [SYZYGY](https://github.com/SoggySaussages/syzygy) hosted by SoggySaussages.`

		return info, nil
	},
}
