package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/config"
	"github.com/Teeworlds-Server-Moderation/discord-moderation/processors"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

var (
	// commands are put in here
	commandChan chan gateway.MessageCreateEvent

	// notification when application is closed
	notify chan os.Signal

	initialized = false
)

var (
	QueueName       = "discord-moderation"
	eventProcessors []processors.EventProcessor
)

func init() {
	commandChan = make(chan gateway.MessageCreateEvent, 1024)

	// cleanup upon application closure
	notify = make(chan os.Signal, 1)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-notify
		time.Sleep(1 * time.Second)

		// graceful shutdown
		close(commandChan)
	}()
	eventProcessors = make([]processors.EventProcessor, 0, 2)
}

// Starts the srvice
func Start(ctx *bot.Context) (err error) {
	if initialized {
		return nil
	}

	brokerSub := config.Broker().Subscriber()
	brokerPub := config.Broker().Publisher()

	initQueuesAndExchanges(brokerSub)
	go eventProcessor(ctx, brokerSub, QueueName)
	go commandProcessor(ctx, brokerPub, commandChan)

	initialized = true
	return nil
}

func AddEventProcessor(processor processors.EventProcessor) {
	eventProcessors = append(eventProcessors, processor)
}

func initQueuesAndExchanges(qcb QueueCreateBinder) {
	createQueueAndBindToExchanges(
		qcb,
		QueueName,
		events.TypeChat,
		events.TypeChatTeam,
		events.TypeChatWhisper,
		events.TypeVoteKickStarted,
		events.TypeVoteSpecStarted,
		events.TypeVoteOptionStarted,
		events.TypeMapChanged,
		events.TypePlayerJoined,
		events.TypePlayerLeft,
		topics.Broadcast,
	)
}
