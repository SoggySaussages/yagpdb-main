package memstats

import (
	"bytes"
	"encoding/json"
	"runtime"

	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/lib/dcmd"
	"github.com/SoggySaussages/sgpdb/lib/discordgo"
	"github.com/SoggySaussages/sgpdb/stdcommands/util"
)

var Command = &commands.YAGCommand{
	Cooldown:             2,
	CmdCategory:          commands.CategoryDebug,
	HideFromCommandsPage: true,
	Name:                 "memstats",
	Description:          "Full memory statistics. Bot Owner Only",
	HideFromHelp:         true,
	RunFunc: util.RequireOwner(func(data *dcmd.Data) (interface{}, error) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		buf, _ := json.Marshal(m)

		send := &discordgo.MessageSend{
			Content: "Memory stats",
			File: &discordgo.File{
				ContentType: "application/json",
				Name:        "memory_stats.json",
				Reader:      bytes.NewReader(buf),
			},
		}

		_, err := common.BotSession.ChannelMessageSendComplex(data.ChannelID, send)

		return nil, err
	}),
}
