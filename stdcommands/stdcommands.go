package stdcommands

import (
	"github.com/SoggySaussages/syzygy/bot"
	"github.com/SoggySaussages/syzygy/bot/eventsystem"
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/common"
	"github.com/SoggySaussages/syzygy/stdcommands/advice"
	"github.com/SoggySaussages/syzygy/stdcommands/allocstat"
	"github.com/SoggySaussages/syzygy/stdcommands/banserver"
	"github.com/SoggySaussages/syzygy/stdcommands/calc"
	"github.com/SoggySaussages/syzygy/stdcommands/catfact"
	"github.com/SoggySaussages/syzygy/stdcommands/ccreqs"
	"github.com/SoggySaussages/syzygy/stdcommands/cleardm"
	"github.com/SoggySaussages/syzygy/stdcommands/createinvite"
	"github.com/SoggySaussages/syzygy/stdcommands/currentshard"
	"github.com/SoggySaussages/syzygy/stdcommands/currenttime"
	"github.com/SoggySaussages/syzygy/stdcommands/customembed"
	"github.com/SoggySaussages/syzygy/stdcommands/dadjoke"
	"github.com/SoggySaussages/syzygy/stdcommands/dcallvoice"
	"github.com/SoggySaussages/syzygy/stdcommands/define"
	"github.com/SoggySaussages/syzygy/stdcommands/dictionary"
	"github.com/SoggySaussages/syzygy/stdcommands/dogfact"
	"github.com/SoggySaussages/syzygy/stdcommands/eightball"
	"github.com/SoggySaussages/syzygy/stdcommands/findserver"
	"github.com/SoggySaussages/syzygy/stdcommands/forex"
	"github.com/SoggySaussages/syzygy/stdcommands/globalrl"
	"github.com/SoggySaussages/syzygy/stdcommands/guildunavailable"
	"github.com/SoggySaussages/syzygy/stdcommands/howlongtobeat"
	"github.com/SoggySaussages/syzygy/stdcommands/info"
	"github.com/SoggySaussages/syzygy/stdcommands/inspire"
	"github.com/SoggySaussages/syzygy/stdcommands/invite"
	"github.com/SoggySaussages/syzygy/stdcommands/leaveserver"
	"github.com/SoggySaussages/syzygy/stdcommands/listflags"
	"github.com/SoggySaussages/syzygy/stdcommands/listroles"
	"github.com/SoggySaussages/syzygy/stdcommands/memstats"
	"github.com/SoggySaussages/syzygy/stdcommands/ping"
	"github.com/SoggySaussages/syzygy/stdcommands/poll"
	"github.com/SoggySaussages/syzygy/stdcommands/roast"
	"github.com/SoggySaussages/syzygy/stdcommands/roll"
	"github.com/SoggySaussages/syzygy/stdcommands/setstatus"
	"github.com/SoggySaussages/syzygy/stdcommands/simpleembed"
	"github.com/SoggySaussages/syzygy/stdcommands/sleep"
	"github.com/SoggySaussages/syzygy/stdcommands/statedbg"
	"github.com/SoggySaussages/syzygy/stdcommands/stateinfo"
	"github.com/SoggySaussages/syzygy/stdcommands/throw"
	"github.com/SoggySaussages/syzygy/stdcommands/toggledbg"
	"github.com/SoggySaussages/syzygy/stdcommands/topcommands"
	"github.com/SoggySaussages/syzygy/stdcommands/topevents"
	"github.com/SoggySaussages/syzygy/stdcommands/topgames"
	"github.com/SoggySaussages/syzygy/stdcommands/topic"
	"github.com/SoggySaussages/syzygy/stdcommands/topservers"
	"github.com/SoggySaussages/syzygy/stdcommands/unbanserver"
	"github.com/SoggySaussages/syzygy/stdcommands/undelete"
	"github.com/SoggySaussages/syzygy/stdcommands/viewperms"
	"github.com/SoggySaussages/syzygy/stdcommands/weather"
	"github.com/SoggySaussages/syzygy/stdcommands/wouldyourather"
	"github.com/SoggySaussages/syzygy/stdcommands/xkcd"
	"github.com/SoggySaussages/syzygy/stdcommands/yagstatus"
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
