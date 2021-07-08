package main

import (
	"log"

	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/processors/dclog"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/processors/vpn"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/service"
	"github.com/diamondburned/arikawa/v2/bot"
)

func main() {
	defer config.Close()

	bot.Run(config.Discord().Token, &Bot{},
		func(ctx *bot.Context) error {
			ctx.HasPrefix = bot.NewPrefix("!")
			// log to discord

			if config.Modules().ErrIfDiscordLoggingDisabled() == nil {
				log.Println("enabled discord logging module")
				service.AddEventProcessor(dclog.DiscordLog)
			}

			if config.Modules().ErrIfVPNDetectionDisabled() == nil {
				log.Println("enabled vpn detection module")
				service.AddEventProcessor(vpn.Detect)
			}

			return service.Start(ctx)
		},
	)
}
