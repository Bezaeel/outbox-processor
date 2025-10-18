package main

import (
	"context"
	"outbox-processor/src/processors"
	"outbox-processor/src/processors/order"
	"outbox-processor/src/processors/shipment"
	"outbox-processor/src/service"
)

func main() {
	ctx := context.Background()
	orderProcessor := order.NewOrderProcessor()
	shipmentProcessor := shipment.NewShipmentProcessor()

	var p []processors.IOutboxProcessor
	p = append(p, orderProcessor)
	p = append(p, shipmentProcessor)
	ps := service.NewOutboxService(ctx, p)

	ps.ExecuteAsync()
}
