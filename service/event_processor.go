package service

import (
	"fmt"
	"log"
	"os"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	a "github.com/streadway/amqp"
)

func eventProcessor(ctx *bot.Context, sub *amqp.Subscriber, queue string) {
	log.Println("Started event processor subroutine...")
	messageChan, err := sub.Consume(queue)
	if err != nil {
		fmt.Printf("Failed to consume from queue %s", queue)
		os.Exit(1)
	}

	for {
		select {
		case <-notify:
			// must be as the first case
			log.Println("Closing event processor subroutine...")
			return
		case msg := <-messageChan:
			processEvent(ctx, msg)
		}
	}
}

func getChannelID(msg a.Delivery) (id discord.ChannelID, eventType string, err error) {
	event := events.BaseEvent{}
	err = event.Unmarshal(string(msg.Body))
	if err != nil {
		return 0, "", err
	}
	value, err := config.Get().GetChannel(event.EventSource)
	if err != nil {
		return 0, event.Type, err
	}
	return value, event.Type, nil
}

func processEvent(ctx *bot.Context, msg a.Delivery) error {
	channelID, eventType, err := getChannelID(msg)
	if err != nil {
		return err
	}

	for _, proc := range eventProcessors {
		err = proc(ctx, channelID, eventType, msg)
		if err != nil {
			ctx.SendMessage(channelID, fmtError(err), nil)
		}
	}
	return nil
}
