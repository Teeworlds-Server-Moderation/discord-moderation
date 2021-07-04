package service

import (
	"log"
	"strings"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/discord-log/config"
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
		case commandMsg := <-commands:
			err := processCommand(commandMsg, ctx, pub)
			if err != nil {
				reply(ctx, commandMsg, err.Error())
			}
		case <-notify:
			log.Println("Closing command processor subroutine...")
			return
		}
	}
}

func getEconAddr(command gateway.MessageCreateEvent) (string, error) {
	return config.Get().GetEconAddr(command.ChannelID)
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
