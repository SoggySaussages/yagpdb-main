package stdcommands

import (
	"github.com/SoggySaussages/sgpdb/bot"
	"github.com/SoggySaussages/sgpdb/bot/eventsystem"
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/stdcommands/advice"
	"github.com/SoggySaussages/sgpdb/stdcommands/allocstat"
	"github.com/SoggySaussages/sgpdb/stdcommands/banserver"
	"github.com/SoggySaussages/sgpdb/stdcommands/calc"
	"github.com/SoggySaussages/sgpdb/stdcommands/catfact"
	"github.com/SoggySaussages/sgpdb/stdcommands/ccreqs"
	"github.com/SoggySaussages/sgpdb/stdcommands/cleardm"
	"github.com/SoggySaussages/sgpdb/stdcommands/createinvite"
	"github.com/SoggySaussages/sgpdb/stdcommands/currentshard"
	"github.com/SoggySaussages/sgpdb/stdcommands/currenttime"
	"github.com/SoggySaussages/sgpdb/stdcommands/customembed"
	"github.com/SoggySaussages/sgpdb/stdcommands/dadjoke"
	"github.com/SoggySaussages/sgpdb/stdcommands/dcallvoice"
	"github.com/SoggySaussages/sgpdb/stdcommands/define"
	"github.com/SoggySaussages/sgpdb/stdcommands/dictionary"
	"github.com/SoggySaussages/sgpdb/stdcommands/dogfact"
	"github.com/SoggySaussages/sgpdb/stdcommands/eightball"
	"github.com/SoggySaussages/sgpdb/stdcommands/findserver"
	"github.com/SoggySaussages/sgpdb/stdcommands/forex"
	"github.com/SoggySaussages/sgpdb/stdcommands/globalrl"
	"github.com/SoggySaussages/sgpdb/stdcommands/guildunavailable"
	"github.com/SoggySaussages/sgpdb/stdcommands/howlongtobeat"
	"github.com/SoggySaussages/sgpdb/stdcommands/info"
	"github.com/SoggySaussages/sgpdb/stdcommands/inspire"
	"github.com/SoggySaussages/sgpdb/stdcommands/invite"
	"github.com/SoggySaussages/sgpdb/stdcommands/leaveserver"
	"github.com/SoggySaussages/sgpdb/stdcommands/listflags"
	"github.com/SoggySaussages/sgpdb/stdcommands/listroles"
	"github.com/SoggySaussages/sgpdb/stdcommands/memstats"
	"github.com/SoggySaussages/sgpdb/stdcommands/ping"
	"github.com/SoggySaussages/sgpdb/stdcommands/poll"
	"github.com/SoggySaussages/sgpdb/stdcommands/roast"
	"github.com/SoggySaussages/sgpdb/stdcommands/roll"
	"github.com/SoggySaussages/sgpdb/stdcommands/say"
	"github.com/SoggySaussages/sgpdb/stdcommands/setstatus"
	"github.com/SoggySaussages/sgpdb/stdcommands/simpleembed"
	"github.com/SoggySaussages/sgpdb/stdcommands/sleep"
	"github.com/SoggySaussages/sgpdb/stdcommands/statedbg"
	"github.com/SoggySaussages/sgpdb/stdcommands/stateinfo"
	"github.com/SoggySaussages/sgpdb/stdcommands/throw"
	"github.com/SoggySaussages/sgpdb/stdcommands/toggledbg"
	"github.com/SoggySaussages/sgpdb/stdcommands/topcommands"
	"github.com/SoggySaussages/sgpdb/stdcommands/topevents"
	"github.com/SoggySaussages/sgpdb/stdcommands/topgames"
	"github.com/SoggySaussages/sgpdb/stdcommands/topic"
	"github.com/SoggySaussages/sgpdb/stdcommands/topservers"
	"github.com/SoggySaussages/sgpdb/stdcommands/unbanserver"
	"github.com/SoggySaussages/sgpdb/stdcommands/undelete"
	"github.com/SoggySaussages/sgpdb/stdcommands/viewperms"
	"github.com/SoggySaussages/sgpdb/stdcommands/weather"
	"github.com/SoggySaussages/sgpdb/stdcommands/wouldyourather"
	"github.com/SoggySaussages/sgpdb/stdcommands/xkcd"
	"github.com/SoggySaussages/sgpdb/stdcommands/yagstatus"
)

var (
	_ bot.BotInitHandler       = (*Plugin)(nil)
	_ commands.CommandProvider = (*Plugin)(nil)
)

type Plugin struct{}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Standard Commands",
		SysName:  "standard_commands",
		Category: common.PluginCategoryCore,
	}
}

func (p *Plugin) AddCommands() {
	commands.AddRootCommands(p,
		// Info
		info.Command,
		invite.Command,

		// Standard
		define.Command,
		weather.Command,
		calc.Command,
		topic.Command,
		catfact.Command,
		dadjoke.Command,
		dogfact.Command,
		advice.Command,
		ping.Command,
		throw.Command,
		roll.Command,
		customembed.Command,
		simpleembed.Command,
		currenttime.Command,
		listroles.Command,
		memstats.Command,
		wouldyourather.Command,
		poll.Command,
		undelete.Command,
		viewperms.Command,
		topgames.Command,
		xkcd.Command,
		howlongtobeat.Command,
		inspire.Command,
		forex.Command,
		roast.Command,
		eightball.Command,
		say.Command,

		// Maintenance
		stateinfo.Command,
		leaveserver.Command,
		banserver.Command,
		cleardm.Command,
		allocstat.Command,
		unbanserver.Command,
		topservers.Command,
		topcommands.Command,
		topevents.Command,
		currentshard.Command,
		guildunavailable.Command,
		yagstatus.Command,
		setstatus.Command,
		createinvite.Command,
		findserver.Command,
		dcallvoice.Command,
		ccreqs.Command,
		sleep.Command,
		toggledbg.Command,
		globalrl.Command,
		listflags.Command,
	)

	statedbg.Commands()
	commands.AddRootCommands(p, dictionary.Command)
}

func (p *Plugin) BotInit() {
	eventsystem.AddHandlerAsyncLastLegacy(p, ping.HandleMessageCreate, eventsystem.EventMessageCreate)
}

func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{})
}
