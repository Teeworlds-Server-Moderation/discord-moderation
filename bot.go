package main

import (
	"fmt"
	"strings"

	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Bot struct {
	Ctx *bot.Context
}

func (b *Bot) Modules(msg *gateway.MessageCreateEvent) (string, error) {
	return strings.Join(config.EnabledModules(), "\n"), nil
}

// Link established a link between a discord channel and a target econ address that can or cannot be correct
// in case th econ address matches any known addresses, its events are then logged to that specific discord channel.
func (b *Bot) Link(msg *gateway.MessageCreateEvent, econAddr string) (string, error) {
	if err := config.Modules().ErrIfDiscordLoggingDisabled(); err != nil {
		return "", err
	}
	err := config.Discord().AddLink(econAddr, msg.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to establish link between this channel and %s: %s", econAddr, err)
	}
	return fmt.Sprintf("established connection between this channel and %s", econAddr), nil
}

// Unlink removes the link between the current channel and its connected econ address.
// no more messages from that address are received anymore.
func (b *Bot) Unlink(msg *gateway.MessageCreateEvent) (string, error) {
	if err := config.Modules().ErrIfDiscordLoggingDisabled(); err != nil {
		return "", err
	}
	addr, err := config.Discord().RemoveChannelLink(msg.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to unlink channel: %s", err)
	}

	return fmt.Sprintf("unlinked this channel from address %s", addr), nil
}
