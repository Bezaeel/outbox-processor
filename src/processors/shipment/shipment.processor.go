package shipment

import (
	"context"
	"fmt"

	"outbox-processor/src/processors"

	"github.com/google/uuid"
)

type ShipmentProcessor struct {
	// telemetry provider
	processors.BaseOutboxProcessor[Model]
}

func NewShipmentProcessor() *ShipmentProcessor {
	return &ShipmentProcessor{}
}

func (o *ShipmentProcessor) CanProcess(message processors.OutboxMessageEntity) bool {
	return o.BaseOutboxProcessor.CanProcess(message)
}

func (o *ShipmentProcessor) ProcessAsync(ctx context.Context, message processors.OutboxMessageEntity) error {
	// publish
	fmt.Println("processing shipment...")

	return o.BaseOutboxProcessor.ProcessAsync(ctx, message, o.processFunc)
}

func (o *ShipmentProcessor) processFunc(ctx context.Context, messageId uuid.UUID, payload interface{}) error {
	fmt.Println("executing shipment...")
	fmt.Println("completing shipment...")
	return nil
}
