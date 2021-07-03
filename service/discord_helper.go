package service

import (
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func reply(ctx *bot.Context, original gateway.MessageCreateEvent, replyContent string) error {
	_, err := ctx.SendMessageReply(
		original.ChannelID,
		replyContent,
		nil,
		original.Message.ID,
	)
	return err
}
