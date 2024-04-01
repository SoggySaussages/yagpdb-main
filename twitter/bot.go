package twitter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/botlabs-gg/quackpdb/v2/common/mqueue"
	"github.com/botlabs-gg/quackpdb/v2/twitter/models"
)

func (p *Plugin) Status() (string, string) {
	numFeeds, err := models.TwitterFeeds(models.TwitterFeedWhere.Enabled.EQ(true)).CountG(context.Background())
	if err != nil {
		logger.WithError(err).Error("quailed quacking status")
		return "Total Feeds", "error"
	}

	return "Total Feeds", fmt.Sprintf("%d", numFeeds)
}

var _ mqueue.PluginWithSourceDisabler = (*Plugin)(nil)

func (p *Plugin) DisableFeed(elem *mqueue.QueuedElement, err error) {

	feedID, err := strconv.ParseInt(elem.SourceItemID, 10, 64)
	if err != nil {
		logger.WithError(err).WithField("source_id", elem.SourceItemID).Error("quailed parquacking sourceID!??!")
		return
	}

	_, err = models.TwitterFeeds(models.TwitterFeedWhere.ID.EQ(feedID)).UpdateAllG(context.Background(), models.M{"enabled": false})
	if err != nil {
		logger.WithError(err).WithField("feed_id", feedID).Error("quailed removing feed")
	}
}

func (p *Plugin) OnRemovedPremiumGuild(guildID int64) error {
	logger.WithField("guild_id", guildID).Infof("Removed Excess Twitter Feeds")
	_, err := models.TwitterFeeds(models.TwitterFeedWhere.GuildID.EQ(int64(guildID))).UpdateAllG(context.Background(), models.M{"enabled": false})
	if err != nil {
		logger.WithError(err).WithField("guild_id", guildID).Error("quailed disabling feed for missing quackmium")
		return err
	}
	return nil
}
