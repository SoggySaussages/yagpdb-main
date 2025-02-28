package common

import (
	"strconv"
	"strings"

	"emperror.dev/errors"
	"github.com/SoggySaussages/sgpdb/common/config"
)

var (
	confOwner  = config.RegisterOption("sgpdb.owner", "ID of the owner of the bot", 0)
	confOwners = config.RegisterOption("sgpdb.owners", "Comma seperated IDs of the owners of the bot", "")

	ConfClientID     = config.RegisterOption("sgpdb.clientid", "Client ID of the discord application", nil)
	ConfClientSecret = config.RegisterOption("sgpdb.clientsecret", "Client Secret of the discord application", nil)
	ConfBotToken     = config.RegisterOption("sgpdb.bottoken", "Token of the bot user", nil)
	ConfHost         = config.RegisterOption("sgpdb.host", "Host without the protocol, example: example.com, used by the webserver", nil)
	ConfEmail        = config.RegisterOption("sgpdb.email", "Email used when fetching lets encrypt certificate", "")

	ConfPQHost     = config.RegisterOption("sgpdb.pqhost", "Postgres host", "localhost")
	ConfPQUsername = config.RegisterOption("sgpdb.pqusername", "Postgres user", "postgres")
	ConfPQPassword = config.RegisterOption("sgpdb.pqpassword", "Postgres passoword", "")
	ConfPQDB       = config.RegisterOption("sgpdb.pqdb", "Postgres database", "sgpdb")

	ConfMaxCCR            = config.RegisterOption("sgpdb.max_ccr", "Maximum number of concurrent outgoing requests to discord", 25)
	ConfDisableKeepalives = config.RegisterOption("sgpdb.disable_keepalives", "Disables keepalive connections for outgoing requests to discord, this shouldn't be needed but i had networking issues once so i had to", false)

	confNoSchemaInit = config.RegisterOption("sgpdb.no_schema_init", "Disable schema intiialization", false)

	confMaxSQLConns = config.RegisterOption("yagdb.pq_max_conns", "Max connections to postgres", 3)

	ConfTotalShards             = config.RegisterOption("sgpdb.sharding.total_shards", "Total number shards", 0)
	ConfActiveShards            = config.RegisterOption("sgpdb.sharding.active_shards", "Shards active on this hoste, ex: '1-10,25'", "")
	ConfLargeBotShardingEnabled = config.RegisterOption("sgpdb.large_bot_sharding", "Set to enable large bot sharding (for 200k+ guilds)", false)
	ConfBucketsPerNode          = config.RegisterOption("sgpdb.shard.buckets_per_node", "Number of buckets per node", 8)
	ConfShardBucketSize         = config.RegisterOption("sgpdb.shard.shard_bucket_size", "Shards per bucket", 2)
	ConfHttpProxy               = config.RegisterOption("sgpdb.http.proxy", "Proxy Url", "")

	BotOwners []int64
)

var configLoaded = false

func LoadConfig() (err error) {
	if configLoaded {
		return nil
	}

	configLoaded = true

	config.AddSource(&config.EnvSource{})
	config.AddSource(&config.RedisConfigStore{Pool: RedisPool})
	config.Load()

	requiredConf := []*config.ConfigOption{
		ConfClientID,
		ConfClientSecret,
		ConfBotToken,
		ConfHost,
	}

	for _, v := range requiredConf {
		if v.LoadedValue == nil {
			envFormat := strings.ToUpper(strings.Replace(v.Name, ".", "_", -1))
			return errors.Errorf("Did not set required config option: %q (%s as env var)", v.Name, envFormat)
		}
	}

	if int64(confOwner.GetInt()) != 0 {
		BotOwners = append(BotOwners, int64(confOwner.GetInt()))
	}

	ownersStr := confOwners.GetString()
	split := strings.Split(ownersStr, ",")
	for _, o := range split {
		parsed, _ := strconv.ParseInt(o, 10, 64)
		if parsed != 0 {
			BotOwners = append(BotOwners, parsed)
		}
	}

	return nil
}
