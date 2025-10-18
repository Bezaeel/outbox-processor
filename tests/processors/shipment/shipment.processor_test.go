package shipment_test

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
	processor := shipment.NewShipmentProcessor()
	ctx := context.Background()
	messageId := uuid.New()

	// act
	err := processor.ProcessAsync(ctx, processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "shipment",
		Payload:     `{"id":123}`,
		PayloadType: reflect.TypeOf(shipment.Model{}).String(),
		Status:      processors.Queued,
	})

	// assert
	assert.NoError(t, err)
}

func TestProcessAsync_ShouldProcessValidPayload(t *testing.T) {
	// arrange
	processor := shipment.NewShipmentProcessor()
	ctx := context.Background()
	messageId := uuid.New()

	model := shipment.Model{
		Id: 456,
	}
	serializedModel, _ := json.Marshal(model)

	message := processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "shipment",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(shipment.Model{}).String(),
		Status:      processors.Queued,
	}

	// act
	err := processor.ProcessAsync(ctx, message)

	// assert
	assert.NoError(t, err)
}

func TestProcessAsync_ShouldHandleInvalidJSONPayload(t *testing.T) {
	// arrange
	processor := shipment.NewShipmentProcessor()
	ctx := context.Background()

	testCases := []struct {
		name        string
		payload     string
		description string
	}{
		{
			name:        "MissingClosingBrace",
			payload:     `{"id":123`,
			description: "JSON missing closing brace",
		},
		{
			name:        "MissingOpeningBrace",
			payload:     `"id":123}`,
			description: "JSON missing opening brace",
		},
		{
			name:        "InvalidJSONSyntax",
			payload:     `{"id":123,}`,
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
			payload:     `{"id":"not_a_number"}`,
			description: "Wrong data type for id field",
		},
		{
			name:        "MissingRequiredField",
			payload:     `{}`,
			description: "Missing required id field",
		},
		{
			name:        "ExtraFields",
			payload:     `{"id":123,"extra":"field"}`,
			description: "Extra fields in JSON",
		},
		{
			name:        "NullValue",
			payload:     `{"id":null}`,
			description: "Null value for id field",
		},
		{
			name:        "ArrayInsteadOfObject",
			payload:     `[{"id":123}]`,
			description: "Array instead of object",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			messageId := uuid.New()
			message := processors.OutboxMessageEntity{
				Id:          messageId,
				Type:        "shipment",
				Payload:     tc.payload,
				PayloadType: reflect.TypeOf(shipment.Model{}).String(),
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
	processor := shipment.NewShipmentProcessor()
	ctx, cancel := context.WithCancel(context.Background())
	messageId := uuid.New()

	model := shipment.Model{
		Id: 789,
	}
	serializedModel, _ := json.Marshal(model)

	message := processors.OutboxMessageEntity{
		Id:          messageId,
		Type:        "shipment",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(shipment.Model{}).String(),
		Status:      processors.Queued,
	}

	// Cancel the context before processing
	cancel()

	// act
	err := processor.ProcessAsync(ctx, message)

	// assert
	assert.NoError(t, err) // The processor should handle context cancellation gracefully
}

func TestNewShipmentProcessor_ShouldReturnValidInstance(t *testing.T) {
	// act
	processor := shipment.NewShipmentProcessor()

	// assert
	assert.NotNil(t, processor)
	assert.IsType(t, &shipment.ShipmentProcessor{}, processor)
}

func TestCanProcess_ShouldReturnFalse_WhenPayloadTypeMismatch(t *testing.T) {
	// arrange
	processor := shipment.NewShipmentProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "shipment",
		Payload:     `{"id":123}`,
		PayloadType: "order.Model", // Wrong payload type
		Status:      processors.Queued,
	}

	// act
	result := processor.CanProcess(message)

	// assert
	assert.False(t, result)
}

func TestCanProcess_ShouldReturnTrue_WhenPayloadTypeMatches(t *testing.T) {
	// arrange
	processor := shipment.NewShipmentProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "shipment",
		Payload:     `{"id":123}`,
		PayloadType: reflect.TypeOf(shipment.Model{}).String(),
		Status:      processors.Queued,
	}

	// act
	result := processor.CanProcess(message)

	// assert
	assert.True(t, result)
}

func TestCanProcess_ShouldReturnFalse_WhenPayloadIsNotValid(t *testing.T) {
	// arrange
	processor := shipment.NewShipmentProcessor()

	testCases := []struct {
		name        string
		payloadType string
		description string
	}{
		{
			name:        "WrongPayloadType",
			payloadType: reflect.TypeOf(order.Model{}).String(),
			description: "Wrong payload type (order instead of shipment)",
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
				Type:        "shipment",
				Payload:     `{"id":123}`,
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
	model := shipment.Model{
		Id: 1,
	}
	serializedModel, _ := json.Marshal(model)

	processor := shipment.NewShipmentProcessor()

	message := processors.OutboxMessageEntity{
		Id:          uuid.New(),
		Type:        "no idea",
		Payload:     string(serializedModel),
		PayloadType: reflect.TypeOf(model).String(),
		Status:      processors.Poisoned,
	}

	// act
	expected := processor.CanProcess(message)

	// assert
	assert.True(t, expected)
}
