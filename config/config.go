package config

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"

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
		addressToChannelStr: make(map[string]string),
		addressToChannel:    make(map[string]discord.ChannelID),
		channelToAddress:    make(map[discord.ChannelID]string),
	}
	unparse, err := configo.ParseWithUnparse(cfg, configo.GetEnv())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cfg.init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// handle unparsing at process termination
	go func(onClose func() error) {
		notify := make(chan os.Signal, 1)
		signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
		onClose()
	}(unparse)
	return cfg
}

type config struct {
	DiscordToken string

	BrokerAddress  string
	BrokerUsername string
	BrokerPassword string

	addressToChannelStr map[string]string

	// addressToChannel maps econ addresses to a Discord channel ID
	addressToChannel map[string]discord.ChannelID
	// channelToAddress maps a discord channel to th ecorresponding econ address
	channelToAddress map[discord.ChannelID]string

	pairDelimiter     string
	keyValueDelimiter string

	sync.RWMutex
}

func (c *config) init() error {
	c.Lock()
	defer c.Unlock()

	for addr, ch := range c.addressToChannelStr {
		value, err := strconv.ParseUint(ch, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid channel ID: %d for %s", value, addr)
		}
		c.addressToChannel[addr] = discord.ChannelID(value)
		c.channelToAddress[discord.ChannelID(value)] = addr
	}
	return nil
}

var addrRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):(\d|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`)

func (c *config) GetChannel(econAddr string) (discord.ChannelID, error) {
	if !addrRegex.MatchString(econAddr) {
		return 0, fmt.Errorf("invalid address: %s", econAddr)
	}
	c.RLock()
	channelID, found := c.addressToChannel[econAddr]
	c.RUnlock()
	if !found {
		return 0, fmt.Errorf("unknown econ address: %s", econAddr)
	}
	return channelID, nil
}
func (c *config) GetEconAddr(channelID discord.ChannelID) (string, error) {

	c.RLock()
	addr, found := c.channelToAddress[channelID]
	c.RUnlock()
	if !found {
		return "", fmt.Errorf("unknown channel id: %d", channelID)
	}
	return addr, nil
}

func (c *config) AddLink(econAddr string, channelID discord.ChannelID) error {
	if !addrRegex.MatchString(econAddr) {
		return fmt.Errorf("invalid address: %s", econAddr)
	}

	c.Lock()
	defer c.Unlock()

	_, found := c.addressToChannel[econAddr]
	if found {
		return fmt.Errorf("address link already exists %s", econAddr)
	}

	_, found = c.channelToAddress[channelID]
	if found {
		return fmt.Errorf("channel link already exists: %d", channelID)
	}

	// neither exist -> create a link
	c.addressToChannel[econAddr] = channelID
	c.channelToAddress[channelID] = econAddr
	return nil
}

// RemoveAddressLink removes the link via its econ address key value
func (c *config) RemoveAddressLink(econAddr string) error {
	if !addrRegex.MatchString(econAddr) {
		return fmt.Errorf("invalid address: %s", econAddr)
	}

	c.Lock()
	defer c.Unlock()

	channel, found := c.addressToChannel[econAddr]
	if !found {
		return fmt.Errorf("address not found %s", econAddr)
	}

	delete(c.addressToChannel, econAddr)
	delete(c.channelToAddress, channel)
	return nil
}

// RemoveChannelLink removes the channel via its channelID key
func (c *config) RemoveChannelLink(channelID discord.ChannelID) error {

	c.Lock()
	defer c.Unlock()

	addr, found := c.channelToAddress[channelID]
	if !found {
		return fmt.Errorf("channel not found %d", channelID)
	}

	delete(c.channelToAddress, channelID)
	delete(c.addressToChannel, addr)
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
						ParseFunction: parsers.Map(&c.addressToChannelStr, &c.pairDelimiter, &c.keyValueDelimiter),
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
