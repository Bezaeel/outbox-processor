package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type IOutboxProcessor interface {
	CanProcess(message OutboxMessageEntity) bool
	ProcessAsync(ctx context.Context, message OutboxMessageEntity) error
}

type BaseOutboxProcessor[T any] struct {
	Value T
}

func (o *BaseOutboxProcessor[T]) CanProcess(message OutboxMessageEntity) bool {
	return reflect.TypeOf(o.Value).String() == message.PayloadType
}

func (o *BaseOutboxProcessor[T]) ProcessAsync(ctx context.Context,
	message OutboxMessageEntity,
	processFunc func(ctx context.Context, messageId uuid.UUID, payload interface{}) error) error {

	var payload T

	err := json.Unmarshal([]byte(message.Payload), &payload)
	if err != nil {
		fmt.Print("unable to deserialize...")
	}

	return processFunc(ctx, message.Id, payload)
}
