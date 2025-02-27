package rsvp

//go:generate sqlboiler --no-hooks psql

import (
	"sync"

	"github.com/SoggySaussages/syzygy/common"
	"github.com/SoggySaussages/syzygy/lib/when"
	"github.com/SoggySaussages/syzygy/lib/when/rules"
	wcommon "github.com/SoggySaussages/syzygy/lib/when/rules/common"
	"github.com/SoggySaussages/syzygy/lib/when/rules/en"
	"github.com/SoggySaussages/syzygy/timezonecompanion/trules"
)

var (
	logger = common.GetPluginLogger(&Plugin{})

	dateParser *when.Parser
)

const (
	EventAccepted  = "event_accepted"
	EventRejected  = "event_rejected"
	EventWaitlist  = "event_waitlist"
	EventUndecided = "event_undecided"
)

func init() {
	dateParser = when.New(&rules.Options{
		Distance:     10,
		MatchByOrder: true})

	dateParser.Add(
		en.Weekday(rules.Override),
		en.CasualDate(rules.Override),
		en.CasualTime(rules.Override),
		trules.Hour(rules.Override),
		trules.HourMinute(rules.Override),
		en.Deadline(rules.Override),
		en.ExactMonthDate(rules.Override),
	)
	dateParser.Add(wcommon.All...)
}

type Plugin struct {
	setupSessions   []*SetupSession
	setupSessionsMU sync.Mutex
}

func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "RSVP",
		SysName:  "rsvp",
		Category: common.PluginCategoryMisc,
	}
}

func RegisterPlugin() {
	p := &Plugin{}

	common.InitSchemas("rsvp", DBSchemas...)
	common.RegisterPlugin(p)
}
