package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/arikawa/v2/discord"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
)

var (
	cfg *config
)

func Get() *config {
	if cfg != nil {
		return cfg
	}
	cfg = &config{
		addressToChannel: make(map[string]string),
		AddressToChannel: make(map[string]discord.ChannelID),
		ChannelToAddress: make(map[discord.ChannelID]string),
	}
	err := configo.Parse(cfg, configo.GetEnv())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cfg.init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfg
}

type config struct {
	DiscordToken string

	BrokerAddress  string
	BrokerUsername string
	BrokerPassword string

	addressToChannel map[string]string

	// AddressToChannel maps econ addresses to a Discord channel ID
	AddressToChannel map[string]discord.ChannelID
	// ChannelToAddress maps a discord channel to th ecorresponding econ address
	ChannelToAddress map[discord.ChannelID]string

	pairDelimiter     string
	keyValueDelimiter string
}

func (c *config) init() error {
	for addr, ch := range c.addressToChannel {
		value, err := strconv.ParseUint(ch, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid channel ID: %d for %s", value, addr)
		}
		c.AddressToChannel[addr] = discord.ChannelID(value)
		c.ChannelToAddress[discord.ChannelID(value)] = addr
	}
	return nil
}

func (c *config) Name() string {
	return "discord-moderation"
}

func (c *config) Options() configo.Options {
	return configo.Options{
		{
			Key:          "ENV_FILE",
			Description:  "path to the location of your .env configuration file (default: ./.env)",
			Mandatory:    true,
			DefaultValue: "./.env",
			ParseFunction: parsers.ReadDotEnvFileMulti(
				configo.Options{
					{
						Key:           "DISCORD_TOKEN",
						Description:   "Create a Discord app at https://discord.com/developers/applications -> Bot -> Token",
						Mandatory:     true,
						ParseFunction: parsers.String(&c.DiscordToken),
					},
					{
						Key:           "PAIR_DELIMITER",
						Description:   "address->channel<delimiter>address2->channel2<delimiter>...",
						DefaultValue:  ",",
						ParseFunction: parsers.String(&c.pairDelimiter),
					},
					{
						Key:           "KEY_VALUE_DELIMITER",
						Description:   "address<delimiter>channel;address2<delimiter>channel2;...",
						DefaultValue:  "->",
						ParseFunction: parsers.String(&c.keyValueDelimiter),
					},
					{
						Key:           "ADDRESS_CHANNEL_MAPPING",
						Description:   "ip:econ_port->discord_channel_id,ip:econ_port2->",
						Mandatory:     true,
						ParseFunction: parsers.Map(&c.addressToChannel, &c.pairDelimiter, &c.keyValueDelimiter),
					},
					{
						Key:           "BROKER_ADDRESS",
						Description:   "The address of your broker in the container is rabbitmq:5672",
						DefaultValue:  "localhost:5672",
						ParseFunction: parsers.String(&c.BrokerAddress),
					},
					{
						Key:           "BROKER_USER",
						Description:   "The user that can access the broker, default: tw-admin",
						DefaultValue:  "tw-admin",
						ParseFunction: parsers.String(&c.BrokerUsername),
					},
					{
						Key:           "BROKER_PASSWORD",
						Mandatory:     true,
						Description:   "The password to access the broker with the corresonding username.",
						ParseFunction: parsers.String(&c.BrokerPassword),
					},
				}..., // slice to var args
			),
		},
	}
}
