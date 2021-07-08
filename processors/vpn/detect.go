package vpn

import (
	"errors"
	"fmt"
	"log"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/jxsl13/goripr"
	a "github.com/streadway/amqp"
)

func Detect(ctx *bot.Context, channelID discord.ChannelID, eventType string, message a.Delivery) error {
	if eventType != events.TypePlayerJoined {
		return nil
	}
	event := events.NewPlayerJoinedEvent()
	err := event.Unmarshal(string(message.Body))
	if err != nil {
		return fmt.Errorf("unable to unmarshal PlayerJoinedEvent: %s", err)
	}
	ripr := config.DetectVPN().RDB()

	log.Printf("Trying to find: '%s'\n", event.IP)
	reason, err := ripr.Find(event.IP)

	if errors.Is(goripr.ErrIPNotFound, err) {
		log.Printf("[NO VPN]: %s\n", event.IP)
		return nil
	} else if err != nil {
		return fmt.Errorf("unexpected error occurred: %s", err)
	}
	if err := config.DetectVPN().RequestBan(event.Player, reason, event.EventSource); err != nil {
		return err
	}
	log.Printf("[IS VPN]: %s\n", event.IP)
	return nil
}
