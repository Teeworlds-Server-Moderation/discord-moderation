package main

import (
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/processors/dclog"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/service"
	"github.com/diamondburned/arikawa/v2/bot"
)

func main() {
	defer config.Close()

	bot.Run(config.Discord().Token, &Bot{},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")
			// log to discord
			service.AddEventProcessor(dclog.DiscordLog)

			return service.Start(ctx)
		},
	)
}
