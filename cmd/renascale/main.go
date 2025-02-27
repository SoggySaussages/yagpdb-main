package main

import (
	"github.com/SoggySaussages/syzygy/analytics"
	"github.com/SoggySaussages/syzygy/antiphishing"
	"github.com/SoggySaussages/syzygy/common/featureflags"
	"github.com/SoggySaussages/syzygy/common/prom"
	"github.com/SoggySaussages/syzygy/common/run"
	"github.com/SoggySaussages/syzygy/genai"
	"github.com/SoggySaussages/syzygy/lib/confusables"
	"github.com/SoggySaussages/syzygy/trivia"
	"github.com/SoggySaussages/syzygy/web/discorddata"

	// Core syzygy packages

	"github.com/SoggySaussages/syzygy/admin"
	"github.com/SoggySaussages/syzygy/bot/paginatedmessages"
	"github.com/SoggySaussages/syzygy/common/internalapi"
	"github.com/SoggySaussages/syzygy/common/scheduledevents2"

	// Plugin imports
	"github.com/SoggySaussages/syzygy/automod"
	"github.com/SoggySaussages/syzygy/automod_legacy"
	"github.com/SoggySaussages/syzygy/autorole"
	"github.com/SoggySaussages/syzygy/cah"
	"github.com/SoggySaussages/syzygy/commands"
	"github.com/SoggySaussages/syzygy/customcommands"
	"github.com/SoggySaussages/syzygy/discordlogger"
	"github.com/SoggySaussages/syzygy/logs"
	"github.com/SoggySaussages/syzygy/moderation"
	"github.com/SoggySaussages/syzygy/notifications"
	"github.com/SoggySaussages/syzygy/premium"
	"github.com/SoggySaussages/syzygy/premium/patreonpremiumsource"
	"github.com/SoggySaussages/syzygy/reddit"
	"github.com/SoggySaussages/syzygy/reminders"
	"github.com/SoggySaussages/syzygy/reputation"
	"github.com/SoggySaussages/syzygy/rolecommands"
	"github.com/SoggySaussages/syzygy/rsvp"
	"github.com/SoggySaussages/syzygy/safebrowsing"
	"github.com/SoggySaussages/syzygy/serverstats"
	"github.com/SoggySaussages/syzygy/soundboard"
	"github.com/SoggySaussages/syzygy/stdcommands"
	"github.com/SoggySaussages/syzygy/streaming"
	"github.com/SoggySaussages/syzygy/tickets"
	"github.com/SoggySaussages/syzygy/timezonecompanion"
	"github.com/SoggySaussages/syzygy/twitter"
	"github.com/SoggySaussages/syzygy/verification"
	"github.com/SoggySaussages/syzygy/youtube"
	// External plugins
)

func main() {

	run.Init()

	//BotSession.LogLevel = discordgo.LogInformational
	paginatedmessages.RegisterPlugin()
	discorddata.RegisterPlugin()

	// Setup plugins
	analytics.RegisterPlugin()
	safebrowsing.RegisterPlugin()
	antiphishing.RegisterPlugin()
	discordlogger.Register()
	commands.RegisterPlugin()
	stdcommands.RegisterPlugin()
	serverstats.RegisterPlugin()
	notifications.RegisterPlugin()
	customcommands.RegisterPlugin()
	reddit.RegisterPlugin()
	moderation.RegisterPlugin()
	reputation.RegisterPlugin()
	streaming.RegisterPlugin()
	automod_legacy.RegisterPlugin()
	automod.RegisterPlugin()
	logs.RegisterPlugin()
	autorole.RegisterPlugin()
	reminders.RegisterPlugin()
	soundboard.RegisterPlugin()
	youtube.RegisterPlugin()
	rolecommands.RegisterPlugin()
	cah.RegisterPlugin()
	tickets.RegisterPlugin()
	verification.RegisterPlugin()
	premium.RegisterPlugin()
	patreonpremiumsource.RegisterPlugin()
	scheduledevents2.RegisterPlugin()
	twitter.RegisterPlugin()
	rsvp.RegisterPlugin()
	timezonecompanion.RegisterPlugin()
	admin.RegisterPlugin()
	internalapi.RegisterPlugin()
	prom.RegisterPlugin()
	featureflags.RegisterPlugin()
	trivia.RegisterPlugin()
	genai.RegisterPlugin()

	// Register confusables replacer
	confusables.Init()

	run.Run()
}
