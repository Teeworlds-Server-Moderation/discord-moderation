package main

import (
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/processors"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/service"
	"github.com/diamondburned/arikawa/v2/bot"
)

func main() {
	cfg := config.Get()
	defer config.Close()

	bot.Run(cfg.DiscordToken, &Bot{},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")

			// log to discord
			service.AddEventProcessor(processors.DiscordLog)

			return service.Start(ctx)
		},
	)
}
