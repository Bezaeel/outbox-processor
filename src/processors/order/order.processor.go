package order

import (
	"context"
	"fmt"
	"outbox-processor/src/processors"

	"github.com/google/uuid"
)

type OrderProcessor struct {
	processors.BaseOutboxProcessor[Model]
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{}
}

func (o *OrderProcessor) CanProcess(message processors.OutboxMessageEntity) bool {
	return o.BaseOutboxProcessor.CanProcess(message)
}

func (o *OrderProcessor) ProcessAsync(ctx context.Context, message processors.OutboxMessageEntity) error {
	// publish
	fmt.Println("processing order...")

	return o.BaseOutboxProcessor.ProcessAsync(ctx, message, o.processFunc)
}

func (o *OrderProcessor) processFunc(ctx context.Context, messageId uuid.UUID, payload interface{}) error {
	fmt.Println("executing order...")
	fmt.Println("completing order...")
	return nil
}
