package config

import (
	"github.com/Teeworlds-Server-Moderation/common/amqp"
	configo "github.com/jxsl13/simple-configo"
	"github.com/jxsl13/simple-configo/parsers"
)

type brokerConfig struct {
	address  string
	username string
	password string

	publisher  *amqp.Publisher
	subscriber *amqp.Subscriber
}

func (bc *brokerConfig) Publisher() *amqp.Publisher {
	return bc.publisher
}

func (bc *brokerConfig) Subscriber() *amqp.Subscriber {
	return bc.subscriber
}

func (bc *brokerConfig) PostParse() error {
	// initialize publisher ans subscriber
	brokerSub, err := amqp.NewSubscriber(bc.address, bc.username, bc.password)
	if err != nil {
		return err
	}

	brokerPub, err := amqp.NewPublisher(bc.address, bc.username, bc.password)
	if err != nil {
		return err
	}
	bc.publisher, bc.subscriber = brokerPub, brokerSub
	return nil
}

func (bc *brokerConfig) Close() error {
	err := bc.publisher.Close()
	if err != nil {
		return err
	}
	err = bc.subscriber.Close()
	if err != nil {
		return err
	}
	return nil
}

func (bc *brokerConfig) Name() string {
	return "broker"
}

func (bc *brokerConfig) Options() configo.Options {
	return configo.Options{
		{
			Key:           "BROKER_ADDRESS",
			Description:   "The address of your broker in the container is rabbitmq:5672",
			Mandatory:     true,
			ParseFunction: parsers.String(&bc.address),
		},
		{
			Key:           "BROKER_USER",
			Description:   "The user that can access the broker, e.g.: tw-admin",
			Mandatory:     true,
			ParseFunction: parsers.String(&bc.username),
		},
		{
			Key:           "BROKER_PASSWORD",
			Mandatory:     true,
			Description:   "The password to access the broker with the corresonding username.",
			ParseFunction: parsers.String(&bc.password),
		},
	}
}
