package service

import (
	"log"
	"strings"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// Execute a specific command
func Command(message gateway.MessageCreateEvent) {
	commandChan <- message
}

func commandProcessor(ctx *bot.Context, pub *amqp.Publisher, commands chan gateway.MessageCreateEvent) {
	log.Println("Starting command processor...")
	for {
		select {
		case <-notify:
			// this must be at first, as it's the most important
			log.Println("Closing command processor subroutine...")
			return
		case commandMsg := <-commands:
			err := processCommand(commandMsg, ctx, pub)
			if err != nil {
				reply(ctx, commandMsg, err.Error())
			}
		}
	}
}

func getEconAddr(command gateway.MessageCreateEvent) (string, error) {
	return config.Discord().GetEconAddr(command.ChannelID)
}

func processCommand(command gateway.MessageCreateEvent, ctx *bot.Context, pub *amqp.Publisher) error {
	econAddr, err := getEconAddr(command)
	if err != nil {
		return err
	}
	cmdExecRequest := events.NewRequestCommandExecEvent()
	cmdExecRequest.Command = strings.Trim(command.Content, " \n\r\t")
	return pub.Publish("", econAddr, cmdExecRequest)
}
