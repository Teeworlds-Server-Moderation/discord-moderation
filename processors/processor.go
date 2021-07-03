package processors

import (
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	a "github.com/streadway/amqp"
)

// EventProcessor is a function that can process events
type EventProcessor func(ctx *bot.Context, channelID discord.ChannelID, eventType string, msg a.Delivery) error
