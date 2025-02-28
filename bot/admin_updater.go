package bot

import (
	"time"

	"github.com/SoggySaussages/sgpdb/bot/eventsystem"
	"github.com/SoggySaussages/sgpdb/common"
	"github.com/SoggySaussages/sgpdb/common/config"
	"github.com/SoggySaussages/sgpdb/lib/discordgo"
	"github.com/mediocregopher/radix/v3"
)

var (
	mainServer         int64
	adminRole          int64
	readOnlyAccessRole int64

	// Set of redis admins
	RedisKeyAdmins    = "sgpdb_admins"
	tmpRedisKeyAdmins = "sgpdb_admins_tmp"
	// Set of users with read only access
	RedisKeyReadOnlyAccess    = "sgpdb_ro_access"
	tmpRedisKeyReadOnlyAccess = "sgpdb_ro_access_tmp"
)

func IsBotAdmin(userID int64) (isAdmin bool, err error) {
	if common.IsOwner(userID) {
		return true, nil
	}

	err = common.RedisPool.Do(radix.FlatCmd(&isAdmin, "SISMEMBER", RedisKeyAdmins, userID))
	return
}

func HasReadOnlyAccess(userID int64) (hasAccess bool, err error) {
	err = common.RedisPool.Do(radix.FlatCmd(&hasAccess, "SISMEMBER", RedisKeyReadOnlyAccess, userID))
	return
}

var stopRunCheckAdmins = make(chan bool)

var (
	confMainServer         = config.RegisterOption("sgpdb.main.server", "Main server used for various purposes, like assigning people with a certain role as bot admins", int64(0))
	confAdminRole          = config.RegisterOption("sgpdb.admin.role", "People with this role on the main server has bot admin status", int64(0))
	confReadOnlyAccessRole = config.RegisterOption("sgpdb.readonly.access.role", "People with this role on the main server has global read only access to configs", int64(0))
)

func loopCheckAdmins() {
	mainServer = int64(confMainServer.GetInt())
	adminRole = int64(confAdminRole.GetInt())
	readOnlyAccessRole = int64(confReadOnlyAccessRole.GetInt())

	if mainServer == 0 || (adminRole == 0 && readOnlyAccessRole == 0) {
		logger.Info("One of SGPDB_MAIN_SERVER, SGPDB_ADMIN_ROLE or SGPDB_READONLY_ACCESS_ROLE not provided, not running admin checker")
		return
	}
	logger.Info("Admin checker running")

	// always skip rename first iteration, in case the last run had an error
	skipRename := true

	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			if ReadyTracker.IsGuildShardReady(mainServer) {
				err := requestCheckBotAdmins(skipRename, mainServer, adminRole, readOnlyAccessRole)
				if err != nil {
					skipRename = true
					logger.WithError(err).Error("failed updating bot admins")
				} else {
					skipRename = false
				}
			}
		case <-stopRunCheckAdmins:
			return
		}

	}
}

func requestCheckBotAdmins(skipRename bool, mainServer, adminRole, readOnlyRole int64) error {
	logger.Info("checking for admins")

	// Swap the keys updated last round, assuming thats done
	if !skipRename {
		common.RedisPool.Do(radix.Cmd(nil, "RENAME", tmpRedisKeyAdmins, RedisKeyAdmins))
		common.RedisPool.Do(radix.Cmd(nil, "RENAME", tmpRedisKeyReadOnlyAccess, RedisKeyReadOnlyAccess))
	}

	err := BatchMemberJobManager.NewBatchMemberJob(mainServer, func(g int64, members []*discordgo.Member) {
		for _, member := range members {
			if adminRole != 0 && common.ContainsInt64Slice(member.Roles, adminRole) {
				err := common.RedisPool.Do(radix.FlatCmd(nil, "SADD", tmpRedisKeyAdmins, member.User.ID))
				if err != nil {
					logger.WithError(err).Error("failed adding user to admins")
				}
			}

			if readOnlyRole != 0 && common.ContainsInt64Slice(member.Roles, readOnlyRole) {
				err := common.RedisPool.Do(radix.FlatCmd(nil, "SADD", tmpRedisKeyReadOnlyAccess, member.User.ID))
				if err != nil {
					logger.WithError(err).Error("failed adding user to read only access users")
				}
			}
		}
	})

	return err
}

func HandleGuildMembersChunk(data *eventsystem.EventData) {
	go BatchMemberJobManager.handleGuildMemberChunk(data)
	if memberFetcher != nil {
		go memberFetcher.HandleGuildmembersChunk(data)
	}
}
