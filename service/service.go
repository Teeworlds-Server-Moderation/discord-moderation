package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Teeworlds-Server-Moderation/common/amqp"
	. "github.com/Teeworlds-Server-Moderation/common/events"
	"github.com/Teeworlds-Server-Moderation/common/topics"
	"github.com/Teeworlds-Server-Moderation/discord-log/config"
	"github.com/Teeworlds-Server-Moderation/discord-log/processors"
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
)

var (
	// commands are put in here
	commandChan chan gateway.MessageCreateEvent

	// notification when application is closed
	notify chan os.Signal

	// broker connections
	brokerPub *amqp.Publisher
	brokerSub *amqp.Subscriber

	initialized = false
)

var (
	QueueName       = "discord-moderation"
	eventProcessors []processors.EventProcessor
)

func init() {
	commandChan = make(chan gateway.MessageCreateEvent, 1024)

	// cleanup upon application closure
	notify = make(chan os.Signal)
	signal.Notify(notify, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-notify
		time.Sleep(1 * time.Second)

		// graceful shutdown
		close(commandChan)
		if brokerSub != nil {
			brokerSub.Close()
			brokerPub.Close()
		}
	}()
	eventProcessors = make([]processors.EventProcessor, 0, 2)
}

func Start(ctx *bot.Context) (err error) {
	if initialized {
		return nil
	}

	cfg := config.Get()

	brokerSub, err = amqp.NewSubscriber(cfg.BrokerAddress, cfg.BrokerUsername, cfg.BrokerPassword)
	if err != nil {
		return err
	}

	brokerPub, err = amqp.NewPublisher(cfg.BrokerAddress, cfg.BrokerUsername, cfg.BrokerPassword)
	if err != nil {
		return err
	}

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
		TypeChat,
		TypeChatTeam,
		TypeChatWhisper,
		TypeVoteKickStarted,
		TypeVoteSpecStarted,
		TypeVoteOptionStarted,
		TypeMapChanged,
		TypePlayerJoined,
		TypePlayerLeft,
		topics.Broadcast,
	)
}
