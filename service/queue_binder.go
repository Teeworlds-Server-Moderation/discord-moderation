package service

import (
	"fmt"
	"log"
	"os"
)

// QueueCreateBinder creates queues and binds them to exchanges
type QueueCreateBinder interface {
	CreateQueue(queue string) error
	BindQueue(queue, exchange string) error
}

func createQueueAndBindToExchanges(qcb QueueCreateBinder, queue string, exchanges ...string) {
	if err := qcb.CreateQueue(queue); err != nil {
		log.Fatalf("Failed to create queue '%s'\n", queue)
	}

	for _, exchange := range exchanges {
		if err := qcb.BindQueue(queue, exchange); err != nil {
			fmt.Printf("Failed to bind queue '%s' to exchange '%s'\n", queue, exchange)
			os.Exit(1)
		}
	}

}
