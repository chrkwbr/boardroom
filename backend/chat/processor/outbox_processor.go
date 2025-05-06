package processor

import (
	"backend/chat/domain"
	"backend/infra/pubsub"
	"database/sql"
	"log"
	"strconv"
	"time"
)

type OutboxProcessor struct {
	db                   *sql.DB
	chatOutboxRepository domain.ChatOutboxRepository
	publisher            pubsub.EventPublisher
	pollingInterval      time.Duration
}

func NewOutboxProcessor(
	db *sql.DB,
	chatOutboxRepository domain.ChatOutboxRepository,
	publisher pubsub.EventPublisher,
	pollingInterval time.Duration,
) *OutboxProcessor {
	return &OutboxProcessor{
		db:                   db,
		chatOutboxRepository: chatOutboxRepository,
		publisher:            publisher,
		pollingInterval:      pollingInterval,
	}
}

func (p *OutboxProcessor) Start() {
	ticker := time.NewTicker(p.pollingInterval)
	go func() {
		for range ticker.C {
			if err := p.processOutbox(); err != nil {
				log.Println("Error processing outbox:", err)
			}
		}
	}()

}

func (p *OutboxProcessor) Stop() {
	// Stop the outbox processor
}

func (p *OutboxProcessor) processOutbox() error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println("Error rolling back transaction:", err)
				return
			}
		}
	}()

	unprocessed, err := p.chatOutboxRepository.FetchUnprocessed(10, tx)
	if err != nil {
		return err
	}
	if len(unprocessed) == 0 {
		return tx.Commit()
	}

	for _, event := range unprocessed {
		log.Println("Processing outbox event:", event.EventId)
		if err := p.publisher.Publish("chat_messages", strconv.FormatInt(event.EventId, 10), event.Payload); err != nil {
			return err
		}

		log.Println("Deleting outbox event:", event.EventId)
		if err := p.chatOutboxRepository.Delete(event, tx); err != nil {
			log.Println("Error deleting outbox event:", err)
			return err
		}
	}

	return tx.Commit()
}
