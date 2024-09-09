package service

import (
	"context"
	"encoding/json"
	"fmt"
	"outbox-processor/src/processors"
	"outbox-processor/src/processors/order"
	"outbox-processor/src/processors/shipment"
	"reflect"

	"github.com/google/uuid"
)

var records []processors.OutboxMessageEntity
var poisonedMessages []processors.OutboxMessageEntity
var failedMessages []processors.OutboxMessageEntity
var successfulMessages []processors.OutboxMessageEntity

type IOutboxService interface {
	Execute()
}

type OutboxService struct {
	ctx        context.Context
	processors []processors.IOutboxProcessor
}

func NewOutboxService(_ctx context.Context, _processor []processors.IOutboxProcessor) *OutboxService {
	return &OutboxService{
		ctx:        _ctx,
		processors: _processor,
	}
}

func (o *OutboxService) ExecuteAsync() {
	// fetch top queued messages
	createRecords()
	var queuedEntries = fetchTopMessages()

	for _, outboxMessage := range queuedEntries {
		var processor = o.getProcessor(outboxMessage)

		if processor == nil {
			poisonedMessages = append(poisonedMessages, outboxMessage)
			continue
		}

		err := processor.ProcessAsync(o.ctx, outboxMessage)
		if err != nil {
			failedMessages = append(failedMessages, outboxMessage)
		}
		successfulMessages = append(successfulMessages, outboxMessage)
	}
	fmt.Printf("successful: %v \n", len(successfulMessages))
	fmt.Printf("poisoned: %v \n", len(poisonedMessages))
	fmt.Printf("failed: %v \n", len(failedMessages))
}

func (o *OutboxService) getProcessor(outboxMessage processors.OutboxMessageEntity) processors.IOutboxProcessor {
	var processor processors.IOutboxProcessor

	for _, v := range o.processors {
		if v.CanProcess(outboxMessage) {
			processor = v
			return processor
		}
	}

	fmt.Print("cannot find processor for message")
	return nil
}

// repo methods
func createRecords() {
	for i := 0; i <= 10; i++ {
		var model interface{}
		if i%2 == 0 {
			model = order.Model{
				Id:    1,
				Name:  fmt.Sprintf("a%v", i),
				Email: "e@exam.ple",
			}
		} else {
			model = shipment.Model{
				Id: i,
			}
		}

		serializedModel, _ := json.Marshal(model)
		entity := processors.OutboxMessageEntity{
			Id:          uuid.New(),
			Payload:     string(serializedModel),
			PayloadType: reflect.TypeOf(model).String(),
			Status:      processors.Queued,
		}
		records = append(records, entity)
	}
}

// fetch top messages
func fetchTopMessages() []processors.OutboxMessageEntity {
	return records[:5]
}
