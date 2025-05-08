package processor

import (
	"backend/chat/command/domain"
	"backend/infra/hub"
	"backend/infra/pubsub"
	"context"
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"
)

type OutboxProcessor struct {
	db                   *sql.DB
	chatOutboxRepository domain.ChatOutboxRepository
	publisher            pubsub.EventPublisher
	hub                  *hub.Hub
	client               *hub.Client
	mu                   sync.Mutex
	processing           bool
}

func NewOutboxProcessor(
	db *sql.DB,
	chatOutboxRepository domain.ChatOutboxRepository,
	publisher pubsub.EventPublisher,
) *OutboxProcessor {
	h, err := hub.GetHubFactory().GetHub(hub.ChatEventOutbox)
	if err != nil {
		log.Panicln("Failed to get hub:", err)
	}

	processor := &OutboxProcessor{
		db:                   db,
		chatOutboxRepository: chatOutboxRepository,
		publisher:            publisher,
		hub:                  h,
		client:               h.CreateAndRegisterClient(256),
	}
	go processor.client.Receive(processor.OutboxReceiver)
	return processor
}

func (p *OutboxProcessor) OutboxReceiver(msg []byte) {
	p.mu.Lock()
	if p.processing {
		log.Println("Already processing outbox, skipping...")
		p.mu.Unlock()
		return
	}
	p.processing = true
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		p.processing = false
		p.mu.Unlock()
	}()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := p.processOutbox(ctx); err != nil {
		log.Println("Error processing outbox:", err)
		return
	}
}

func (p *OutboxProcessor) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client != nil {
		p.hub.UnregisterClient(p.client)
		p.client = nil
	}
}

func (p *OutboxProcessor) processOutbox(ctx context.Context) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			if err := tx.Rollback(); err != nil {
				log.Println("Error rolling back transaction:", err)
			}
		}
	}()

	unprocessed, err := p.chatOutboxRepository.FetchUnprocessed(10, tx)
	if err != nil {
		return err
	}
	if len(unprocessed) == 0 {
		log.Println("No unprocessed outbox events found")
		if err := tx.Commit(); err != nil {
			return err
		}
		committed = true
		return nil
	}

	for _, event := range unprocessed {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := p.publisher.Publish("chat_messages", strconv.FormatInt(event.EventId, 10), event.Payload); err != nil {
				return err
			}
		}

		if err := p.chatOutboxRepository.Delete(event, tx); err != nil {
			return err
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err = tx.Commit(); err != nil {
			return err
		}
		committed = true
		return nil
	}
}
