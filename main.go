package main

import (
	"github.com/Teeworlds-Server-Moderation/discord-log/config"
	"github.com/Teeworlds-Server-Moderation/discord-log/processors"
	"github.com/Teeworlds-Server-Moderation/discord-log/service"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

func main() {
	bot.Run(config.Get().DiscordToken, &Bot{},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")

			// log to discord
			service.AddEventProcessor(processors.DiscordLog)

			return service.Start(ctx)
		},
	)
}

type Bot struct {
	Ctx *bot.Context
}

func (b *Bot) Ping(*gateway.MessageCreateEvent) (string, error) {
	return "Pong!", nil
}
