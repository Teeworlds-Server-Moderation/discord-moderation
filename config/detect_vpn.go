package config

import (
	"sync"
	"time"

	"github.com/jxsl13/goripr"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
	"github.com/jxsl13/simple-configo/unparsers"
)

var (
	folderRegex  = `^[a-zA-Z0-9-]+$`
	errFolderMsg = "The folder name must only contain alphanumeric characters, no special caracters nor whitespaces and must not be empty."
)

type detectVPNConfig struct {
	// initialization of
	dataPath        string
	blacklistFolder string
	whitelistFolder string
	redisAddress    string
	redisPassword   string
	redisDatabase   int
	rdb             *goripr.Client

	// these below parameters are guarded
	broadcastBans bool
	banCommand    string
	banReason     string
	banDuration   time.Duration
	// only guards the above parameters
	sync.RWMutex
}

func (dvc *detectVPNConfig) RDB() *goripr.Client {
	return dvc.rdb
}

func (dvc *detectVPNConfig) BroadcastBans() bool {
	dvc.RLock()
	defer dvc.RUnlock()
	return dvc.broadcastBans
}

func (dvc *detectVPNConfig) SetBroadcastBans(value bool) {
	dvc.Lock()
	defer dvc.Unlock()
	dvc.broadcastBans = value
}

func (dvc *detectVPNConfig) BanCommand() string {
	dvc.RLock()
	defer dvc.RUnlock()
	return dvc.banCommand
}

func (dvc *detectVPNConfig) SetBanCommand(value string) {
	dvc.Lock()
	defer dvc.Unlock()
	dvc.banCommand = value
}

func (dvc *detectVPNConfig) BanReason() string {
	dvc.RLock()
	defer dvc.RUnlock()
	return dvc.banReason
}

func (dvc *detectVPNConfig) SetBanReason(value string) {
	dvc.Lock()
	defer dvc.Unlock()
	dvc.banReason = value
}

func (dvc *detectVPNConfig) BanDuration() time.Duration {
	dvc.RLock()
	defer dvc.RUnlock()
	return dvc.banDuration
}

func (dvc *detectVPNConfig) SetBanDuration(value string) error {
	dvc.Lock()
	defer dvc.Unlock()
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	dvc.banDuration = duration
	return nil
}

// Initaliztion and closing
func (dvc *detectVPNConfig) PostParse() error {
	err := dvc.initFolderStructure()
	if err != nil {
		return err
	}
	err = dvc.updateRedisDatabase()
	if err != nil {
		return err
	}

	dvc.rdb, err = goripr.NewClient(goripr.Options{
		Addr:     dvc.redisAddress,
		Password: dvc.redisPassword,
		DB:       dvc.redisDatabase,
	})
	if err != nil {
		return err
	}

	return nil
}

func (dvc *detectVPNConfig) Close() error {
	return dvc.rdb.Close()
}

func (dvc *detectVPNConfig) Name() string {
	return "detect-vpn"
}

func (dvc *detectVPNConfig) Options() configo.Options {
	optionsList := configo.Options{
		{
			Key:             "REDIS_ADDRESS",
			Mandatory:       true,
			Description:     "The REDIS_ADDRESS must have the following format: <hostname/ip>:<port>",
			DefaultValue:    "localhost:6379",
			ParseFunction:   parsers.String(&dvc.redisAddress),
			UnparseFunction: unparsers.String(&dvc.redisAddress),
		},
		{
			Key:             "REDIS_PASSWORD",
			Description:     "Pasword used for the redis database, can be left empty.",
			ParseFunction:   parsers.String(&dvc.redisPassword),
			UnparseFunction: unparsers.String(&dvc.redisPassword),
		},
		{
			Key:             "REDIS_DB",
			Description:     "Is one of the 16 [0:15] ditinct databases that redis offers.",
			DefaultValue:    "1",
			ParseFunction:   parsers.RangesInt(&dvc.redisDatabase, 0, 15),
			UnparseFunction: unparsers.Int(&dvc.redisDatabase),
		},
		{
			Key:             "DATA_PATH",
			Description:     "Is the root folder that contains all of the data of this service.",
			DefaultValue:    "./data",
			ParseFunction:   parsers.String(&dvc.dataPath),
			UnparseFunction: unparsers.String(&dvc.dataPath),
		},
		{
			Key:             "BLACKLIST_FOLDER",
			Description:     "This is a folder WITHIN the DATA_PATH that is created and used to store and retrieve blacklists",
			DefaultValue:    "blacklists",
			ParseFunction:   parsers.Regex(&dvc.blacklistFolder, folderRegex, errFolderMsg),
			UnparseFunction: unparsers.String(&dvc.blacklistFolder),
		},
		{
			Key:             "WHITELIST_FOLDER",
			Description:     "This is a folder WITHIN the DATA_PATH that is created and used to whitelists",
			DefaultValue:    "whitelists",
			ParseFunction:   parsers.Regex(&dvc.whitelistFolder, folderRegex, errFolderMsg),
			UnparseFunction: unparsers.String(&dvc.whitelistFolder),
		},
		{
			Key:             "BAN_REASON",
			Description:     "The default reason that is used when a specific ban range does not specify a reason with # comments",
			DefaultValue:    "VPN",
			ParseFunction:   parsers.String(&dvc.banReason),
			UnparseFunction: unparsers.String(&dvc.banReason),
		},
		{
			Key:             "BAN_DURATION",
			Description:     "The duration a VPN IP is banned by default.(e.g. 10s, 5m, 1h, 1h5m10s, 24h)",
			DefaultValue:    "24h",
			ParseFunction:   parsers.Duration(&dvc.banDuration),
			UnparseFunction: unparsers.Duration(&dvc.banDuration),
		},
		{
			Key:             "BROADCAST_BANS",
			Description:     "If a VPN user is detected on one server, you may want to execute the ban command on all servers that are connected to this system.",
			DefaultValue:    "false",
			ParseFunction:   parsers.Bool(&dvc.broadcastBans),
			UnparseFunction: unparsers.Bool(&dvc.broadcastBans),
		},
		{
			Key:             "BAN_COMMAND",
			Description:     "You may use the variables {IP}, {ID}, {DURATION:MINUTES}, {DURATION:SECONDS}, {REASON}",
			DefaultValue:    "ban {IP} {DURATION:MINUTES} {REASON}",
			ParseFunction:   parsers.String(&dvc.banCommand),
			UnparseFunction: unparsers.String(&dvc.banCommand),
		},
	}

	return optionsList
}
