package main

import (
	"fmt"

	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Bot struct {
	Ctx *bot.Context
}

func (b *Bot) Ping(msg *gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}

func (b *Bot) Link(msg *gateway.MessageCreateEvent, econAddr string) (string, error) {
	err := config.Get().AddLink(econAddr, msg.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to establish link between this channel and %s: %s", econAddr, err)
	}
	return fmt.Sprintf("established connection between this channel and %s", econAddr), nil
}

func (b *Bot) Unlink(msg *gateway.MessageCreateEvent) (string, error) {
	addr, err := config.Get().RemoveChannelLink(msg.ChannelID)
	if err != nil {
		return "", fmt.Errorf("failed to unlink channel: %s", err)
	}

	return fmt.Sprintf("unlinked this channel from address %s", addr), nil
}
