package order_test

import (
	"context"
	"encoding/json"
	"outbox-processor/src/processors"
	"outbox-processor/src/processors/order"
	"outbox-processor/src/processors/shipment"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProcessFunc_ShouldProcessPayload(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()
	ctx := context.Background()
	messageId := uuid.New()

	// act
	err := processor.ProcessAsync(ctx, processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "order",
		Payload:     `{"id":123,"name":"Test Order","email":"test@example.com"}`,
		PayloadType: reflect.TypeOf(order.Model{}).String(),
		Status:      processors.Queued,
	})

	// assert
	assert.NoError(t, err)
}

func TestCanProcess_ShouldReturnFalse_WhenPayloadIsNotValid(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()

	testCases := []struct {
		name        string
		payloadType string
		description string
	}{
		{
			name:        "WrongPayloadType",
			payloadType: reflect.TypeOf(shipment.Model{}).String(),
			description: "Wrong payload type (shipment instead of order)",
		},
		{
			name:        "InvalidPayloadType",
			payloadType: "invalid.Model",
			description: "Invalid payload type string",
		},
		{
			name:        "EmptyPayloadType",
			payloadType: "",
			description: "Empty payload type",
		},
		{
			name:        "NilPayloadType",
			payloadType: "<nil>",
			description: "Nil payload type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			message := processors.OutboxMessageEntity{
				Id:          uuid.New(),
				Type:        "order",
				Payload:     `{"id":123,"name":"Test","email":"test@example.com"}`,
				PayloadType: tc.payloadType,
				Status:      processors.Queued,
			}

			// act
			result := processor.CanProcess(message)

			// assert
			assert.False(t, result, "Should return false for %s", tc.description)
		})
	}
}

func TestCanProcess_ShouldReturnTrue_WhenPayloadIsValid(t *testing.T) {
	// arrange

	model := order.Model{
		Id:    1,
		Name:  "test",
		Email: "e@exam.ple",
	}
	serializedModel, _ := json.Marshal(model)

	_order := order.NewOrderProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "no idea",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(model).String(),
		Status:      processors.Poisoned,
	}

	// act
	expected := _order.CanProcess(message)

	// assert
	assert.True(t, expected)
}

func TestProcessAsync_ShouldProcessValidPayload(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()
	ctx := context.Background()
	messageId := uuid.New()

	model := order.Model{
		Id:    456,
		Name:  "Valid Order",
		Email: "valid@example.com",
	}
	serializedModel, _ := json.Marshal(model)

	message := processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "order",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(order.Model{}).String(),
		Status:      processors.Queued,
	}

	// act
	err := processor.ProcessAsync(ctx, message)

	// assert
	assert.NoError(t, err)
}

func TestProcessAsync_ShouldHandleInvalidJSONPayload(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()
	ctx := context.Background()

	testCases := []struct {
		name        string
		payload     string
		description string
	}{
		{
			name:        "MissingClosingBrace",
			payload:     `{"id":123,"name":"Invalid JSON"`,
			description: "JSON missing closing brace",
		},
		{
			name:        "MissingOpeningBrace",
			payload:     `"id":123,"name":"Invalid JSON"}`,
			description: "JSON missing opening brace",
		},
		{
			name:        "InvalidJSONSyntax",
			payload:     `{"id":123,"name":"Invalid JSON",}`,
			description: "JSON with trailing comma",
		},
		{
			name:        "EmptyPayload",
			payload:     ``,
			description: "Empty payload string",
		},
		{
			name:        "NonJSONString",
			payload:     `not a json string`,
			description: "Plain text instead of JSON",
		},
		{
			name:        "WrongDataType",
			payload:     `{"id":"not_a_number","name":"Test","email":"test@example.com"}`,
			description: "Wrong data type for id field",
		},
		{
			name:        "MissingRequiredField",
			payload:     `{"name":"Test","email":"test@example.com"}`,
			description: "Missing required id field",
		},
		{
			name:        "ExtraFields",
			payload:     `{"id":123,"name":"Test","email":"test@example.com","extra":"field"}`,
			description: "Extra fields in JSON",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			messageId := uuid.New()
			message := processors.OutboxMessageEntity{
				Id:          messageId,
				Type:        "order",
				Payload:     tc.payload,
				PayloadType: reflect.TypeOf(order.Model{}).String(),
				Status:      processors.Queued,
			}

			// act
			err := processor.ProcessAsync(ctx, message)

			// assert
			assert.NoError(t, err, "Should handle %s gracefully", tc.description)
		})
	}
}

func TestProcessAsync_ShouldHandleContextCancellation(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()
	ctx, cancel := context.WithCancel(context.Background())
	messageId := uuid.New()

	model := order.Model{
		Id:    789,
		Name:  "Cancelled Order",
		Email: "cancelled@example.com",
	}
	serializedModel, _ := json.Marshal(model)

	message := processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "order",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(order.Model{}).String(),
		Status:      processors.Queued,
	}

	// Cancel the context before processing
	cancel()

	// act
	err := processor.ProcessAsync(ctx, message)

	// assert
	assert.NoError(t, err) // The processor should handle context cancellation gracefully
}

func TestNewOrderProcessor_ShouldReturnValidInstance(t *testing.T) {
	// act
	processor := order.NewOrderProcessor()

	// assert
	assert.NotNil(t, processor)
	assert.IsType(t, &order.OrderProcessor{}, processor)
}

func TestCanProcess_ShouldReturnFalse_WhenPayloadTypeMismatch(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "order",
		Payload:     `{"id":123,"name":"Test","email":"test@example.com"}`,
		PayloadType: "shipment.Model", // Wrong payload type
		Status:      processors.Queued,
	}

	// act
	result := processor.CanProcess(message)

	// assert
	assert.False(t, result)
}

func TestCanProcess_ShouldReturnTrue_WhenPayloadTypeMatches(t *testing.T) {
	// arrange
	processor := order.NewOrderProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "order",
		Payload:     `{"id":123,"name":"Test","email":"test@example.com"}`,
		PayloadType: reflect.TypeOf(order.Model{}).String(),
		Status:      processors.Queued,
	}

	// act
	result := processor.CanProcess(message)

	// assert
	assert.True(t, result)
}
