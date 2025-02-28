package main

import (
	"github.com/SoggySaussages/sgpdb/analytics"
	"github.com/SoggySaussages/sgpdb/antiphishing"
	"github.com/SoggySaussages/sgpdb/common/featureflags"
	"github.com/SoggySaussages/sgpdb/common/prom"
	"github.com/SoggySaussages/sgpdb/common/run"
	"github.com/SoggySaussages/sgpdb/genai"
	"github.com/SoggySaussages/sgpdb/lib/confusables"
	"github.com/SoggySaussages/sgpdb/trivia"
	"github.com/SoggySaussages/sgpdb/web/discorddata"

	// Core sgpdb packages

	"github.com/SoggySaussages/sgpdb/admin"
	"github.com/SoggySaussages/sgpdb/bot/paginatedmessages"
	"github.com/SoggySaussages/sgpdb/common/internalapi"
	"github.com/SoggySaussages/sgpdb/common/scheduledevents2"

	// Plugin imports
	"github.com/SoggySaussages/sgpdb/automod"
	"github.com/SoggySaussages/sgpdb/automod_legacy"
	"github.com/SoggySaussages/sgpdb/autorole"
	"github.com/SoggySaussages/sgpdb/cah"
	"github.com/SoggySaussages/sgpdb/commands"
	"github.com/SoggySaussages/sgpdb/customcommands"
	"github.com/SoggySaussages/sgpdb/discordlogger"
	"github.com/SoggySaussages/sgpdb/logs"
	"github.com/SoggySaussages/sgpdb/moderation"
	"github.com/SoggySaussages/sgpdb/notifications"
	"github.com/SoggySaussages/sgpdb/premium"
	"github.com/SoggySaussages/sgpdb/premium/patreonpremiumsource"
	"github.com/SoggySaussages/sgpdb/reddit"
	"github.com/SoggySaussages/sgpdb/reminders"
	"github.com/SoggySaussages/sgpdb/reputation"
	"github.com/SoggySaussages/sgpdb/rolecommands"
	"github.com/SoggySaussages/sgpdb/rsvp"
	"github.com/SoggySaussages/sgpdb/safebrowsing"
	"github.com/SoggySaussages/sgpdb/serverstats"
	"github.com/SoggySaussages/sgpdb/soundboard"
	"github.com/SoggySaussages/sgpdb/stdcommands"
	"github.com/SoggySaussages/sgpdb/streaming"
	"github.com/SoggySaussages/sgpdb/tickets"
	"github.com/SoggySaussages/sgpdb/timezonecompanion"
	"github.com/SoggySaussages/sgpdb/twitter"
	"github.com/SoggySaussages/sgpdb/verification"
	"github.com/SoggySaussages/sgpdb/youtube"
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
