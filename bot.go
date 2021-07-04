package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

type Bot struct {
	Ctx *bot.Context
}

func (b *Bot) Ping(msg *gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}

func (b *Bot) Link(msg *gateway.MessageCreateEvent) (string, error) {
	//channelID := msg.ID
	text := msg.Content
	return fmt.Sprintf("'%s'", text), nil
}
