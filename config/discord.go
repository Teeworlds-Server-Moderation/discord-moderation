package config

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
	"github.com/jxsl13/simple-configo/unparsers"
)

var (
	addrRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):(\d|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`)
)

func newDiscordConfig() *discordConfig {
	config := &discordConfig{
		addressToChannelStr: make(map[string]string),
		addressToChannel:    make(map[string]discord.ChannelID),
		channelToAddress:    make(map[discord.ChannelID]string),
	}

	return config
}

type discordConfig struct {
	Token string

	addressToChannelStr map[string]string

	// addressToChannel maps econ addresses to a Discord channel ID
	addressToChannel map[string]discord.ChannelID
	// channelToAddress maps a discord channel to th ecorresponding econ address
	channelToAddress map[discord.ChannelID]string

	pairDelimiter     string
	keyValueDelimiter string

	skipJoinLeaveMessages bool
	skipWhisperMessages   bool

	sync.RWMutex
}

func (dlc *discordConfig) PostParse() error {
	dlc.Lock()
	defer dlc.Unlock()

	for addr, ch := range dlc.addressToChannelStr {
		value, err := strconv.ParseUint(ch, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid channel ID: %d for %s", value, addr)
		}
		dlc.addressToChannel[addr] = discord.ChannelID(value)
		dlc.channelToAddress[discord.ChannelID(value)] = addr
	}
	return nil
}

func (dlc *discordConfig) Close() error {
	return nil
}

func (dlc *discordConfig) GetChannel(econAddr string) (discord.ChannelID, error) {
	if !addrRegex.MatchString(econAddr) {
		return 0, fmt.Errorf("invalid address: %s", econAddr)
	}
	dlc.RLock()
	channelID, found := dlc.addressToChannel[econAddr]
	dlc.RUnlock()
	if !found {
		return 0, fmt.Errorf("unknown econ address: %s", econAddr)
	}
	return channelID, nil
}
func (dlc *discordConfig) GetEconAddr(channelID discord.ChannelID) (string, error) {

	dlc.RLock()
	addr, found := dlc.channelToAddress[channelID]
	dlc.RUnlock()
	if !found {
		return "", fmt.Errorf("unknown channel id: %d", channelID)
	}
	return addr, nil
}

func (dlc *discordConfig) AddLink(econAddr string, channelID discord.ChannelID) error {
	if !addrRegex.MatchString(econAddr) {
		return fmt.Errorf("invalid address: %s", econAddr)
	}

	dlc.Lock()
	defer dlc.Unlock()

	_, found := dlc.addressToChannel[econAddr]
	if found {
		return fmt.Errorf("address link already exists %s", econAddr)
	}

	_, found = dlc.channelToAddress[channelID]
	if found {
		return fmt.Errorf("channel link already exists: %d", channelID)
	}

	// neither exist -> create a link
	dlc.addressToChannel[econAddr] = channelID
	dlc.addressToChannelStr[econAddr] = strconv.FormatUint(uint64(channelID), 10)
	dlc.channelToAddress[channelID] = econAddr
	return nil
}

// RemoveAddressLink removes the link via its econ address key value
func (dlc *discordConfig) RemoveAddressLink(econAddr string) (discord.ChannelID, error) {
	if !addrRegex.MatchString(econAddr) {
		return 0, fmt.Errorf("invalid address: %s", econAddr)
	}

	dlc.Lock()
	defer dlc.Unlock()

	channel, found := dlc.addressToChannel[econAddr]
	if !found {
		return 0, fmt.Errorf("address not found %s", econAddr)
	}

	delete(dlc.addressToChannelStr, econAddr)
	delete(dlc.addressToChannel, econAddr)
	delete(dlc.channelToAddress, channel)
	return channel, nil
}

// RemoveChannelLink removes the channel via its channelID key
func (dlc *discordConfig) RemoveChannelLink(channelID discord.ChannelID) (string, error) {

	dlc.Lock()
	defer dlc.Unlock()

	addr, found := dlc.channelToAddress[channelID]
	if !found {
		return "", fmt.Errorf("channel not found %d", channelID)
	}

	delete(dlc.channelToAddress, channelID)
	delete(dlc.addressToChannel, addr)
	delete(dlc.addressToChannelStr, addr)
	return addr, nil
}

func (dlc *discordConfig) GetSkipJoinLeaveMessages() bool {
	dlc.RLock()
	defer dlc.RUnlock()
	return dlc.skipJoinLeaveMessages
}

func (dlc *discordConfig) SetSkipJoinLeaveMessages(value bool) {
	dlc.Lock()
	defer dlc.Unlock()
	dlc.skipJoinLeaveMessages = value
}

func (dlc *discordConfig) GetSkipWhisperMessages() bool {
	dlc.RLock()
	defer dlc.RUnlock()
	return dlc.skipWhisperMessages
}

func (dlc *discordConfig) SetSkipWhisperMessages(value bool) {
	dlc.Lock()
	defer dlc.Unlock()
	dlc.skipWhisperMessages = value
}

func (dlc *discordConfig) Name() string {
	return "discord"
}

func (dlc *discordConfig) Options() configo.Options {
	options := configo.Options{
		{
			Key:             "DISCORD_TOKEN",
			Description:     "Create a Discord app at https://discord.com/developers/applications -> Bot -> Token",
			Mandatory:       true,
			ParseFunction:   parsers.String(&dlc.Token),
			UnparseFunction: unparsers.String(&dlc.Token),
		},
		{
			Key:             "PAIR_DELIMITER",
			Description:     "address->channel<delimiter>address2->channel2<delimiter>...",
			DefaultValue:    ",",
			ParseFunction:   parsers.String(&dlc.pairDelimiter),
			UnparseFunction: unparsers.String(&dlc.pairDelimiter),
		},
		{
			Key:             "KEY_VALUE_DELIMITER",
			Description:     "address<delimiter>channel;address2<delimiter>channel2;...",
			DefaultValue:    "->",
			ParseFunction:   parsers.String(&dlc.keyValueDelimiter),
			UnparseFunction: unparsers.String(&dlc.keyValueDelimiter),
		},
		{
			Key:             "ADDRESS_CHANNEL_MAPPING",
			Description:     "ip:econ_port->discord_channel_id,ip:econ_port2->",
			Mandatory:       true,
			ParseFunction:   parsers.Map(&dlc.addressToChannelStr, &dlc.pairDelimiter, &dlc.keyValueDelimiter),
			UnparseFunction: unparsers.Map(&dlc.addressToChannelStr, &dlc.pairDelimiter, &dlc.keyValueDelimiter),
		},
		{
			Key:             "LOGS_SKIP_JOIN_LEAVE",
			DefaultValue:    "true",
			Description:     "Whether to skip logging joining and leaving player messages in discord. (default: true)",
			ParseFunction:   parsers.Bool(&dlc.skipJoinLeaveMessages),
			UnparseFunction: unparsers.Bool(&dlc.skipJoinLeaveMessages),
		},
		{
			Key:             "LOGS_SKIP_WHISPER",
			DefaultValue:    "true",
			Description:     "Whether to skip logging of whisper messages in discord. (default: true)",
			ParseFunction:   parsers.Bool(&dlc.skipWhisperMessages),
			UnparseFunction: unparsers.Bool(&dlc.skipWhisperMessages),
		},
	}
	return options
}
