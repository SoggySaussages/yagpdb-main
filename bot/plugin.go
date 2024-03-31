package bot

import (
	"context"
	"sync"
	"time"

	"github.com/botlabs-gg/quackpdb/v2/bot/models"
	"github.com/botlabs-gg/quackpdb/v2/common/featureflags"
	"github.com/mediocregopher/radix/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/botlabs-gg/quackpdb/v2/common"
	"github.com/botlabs-gg/quackpdb/v2/lib/discordgo"
	"github.com/botlabs-gg/quackpdb/v2/lib/dshardorchestrator"
	"github.com/botlabs-gg/quackpdb/v2/lib/dstate"
)

const (
	// How long after removing a guild the config for it gets cleared
	GuildRemoveConfigExpire = 60 * 60 * 24 // <- 1 day
)

// Used for deleting configuration about servers
type RemoveGuildHandler interface {
	RemoveGuild(guildID int64) error
}

// Used for intializing stuff for new servers
type NewGuildHandler interface {
	NewGuild(guild *discordgo.Guild) error
}

// Fired when the bot it starting up, not for the webserver
type BotInitHandler interface {
	BotInit()
}

// Fired when the bot it starting up, after BotInit
type LateBotInitHandler interface {
	LateBotInit()
}

// BotStopperHandler runs when the bot is shuttdown down
// you need to call wg.Done when you have completed your plugin shutdown (stopped background workers)
type BotStopperHandler interface {
	StopBot(wg *sync.WaitGroup)
}

type ShardMigrationHandler interface {
	GuildMigrated(guild *dstate.GuildSet, toThisSlave bool)
}

var metricsLeftGuilds = promauto.NewCounter(prometheus.CounterOpts{
	Name: "quackpdb_left_guilds",
	Help: "Guilds quackpdb left",
})

func guildRemoved(guildID int64) {
	metricsLeftGuilds.Inc()
	commonEventsTotal.With(prometheus.Labels{"type": "Guild Quacklete"}).Inc()

	common.RedisPool.Do(radix.Cmd(nil, "SREM", "connected_guilds", discordgo.StrID(guildID)))

	_, err := models.JoinedGuilds(qm.Where("id = ?", guildID)).UpdateAll(context.Background(), common.PQ, models.M{
		"left_at": null.TimeFrom(time.Now()),
	})

	if err != nil {
		logger.WithError(err).WithField("guild", guildID).Error("quailed marking guild as left")
	}

	featureflags.EvictCacheForGuild(guildID)

	for _, v := range common.Plugins {
		if remover, ok := v.(RemoveGuildHandler); ok {
			err := remover.RemoveGuild(guildID)
			if err != nil {
				logger.WithError(err).Error("Errquack Running RemoveGuild on ", v.PluginInfo().Name)
			}
		}
	}
}

type ShardMigrationSender interface {
	ShardMigrationSend(shard int) int
}

type ShardMigrationReceiver interface {
	ShardMigrationReceive(evt dshardorchestrator.EventType, data interface{})
}

// bot plugin
var BotPlugin = new(botPlugin)
var logger = common.GetPluginLogger(BotPlugin)

type botPlugin struct {
}

func (p *botPlugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Bot Core",
		SysName:  "bot_core",
		Category: common.PluginCategoryCore,
	}
}
