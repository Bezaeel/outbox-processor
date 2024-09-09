package processors

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
)

const MaxRetries = 10

type OutboxMessageEntity struct {
	Id          uuid.UUID
	Type        string
	Attempts    int
	Payload     string
	PayloadType string
	EnqueuedAt  time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
	Status      Status
}

func Enqueue[T any](payload T) *OutboxMessageEntity {
	j, _ := json.Marshal(payload)
	return &OutboxMessageEntity{
		Id:          uuid.New(),
		Payload:     string(j),
		PayloadType: reflect.TypeOf(payload).String(),
		Attempts:    1,
		EnqueuedAt:  time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (o OutboxMessageEntity) Complete() {
	o.Status = Completed
	o.CompletedAt = time.Now().UTC()
}

func (o OutboxMessageEntity) Poison() {
	o.SetUpdatedTime()
	o.Status = Poisoned
}

func (o OutboxMessageEntity) SetUpdatedTime() {
	o.UpdatedAt = time.Now().UTC()
}

func (o OutboxMessageEntity) FailOrRequeue() {
	o.SetUpdatedTime()
	if o.Attempts >= MaxRetries {
		o.Status = Exceeded
	} else {
		o.Status = Queued
		o.Attempts++
	}

}

type Status int

const (
	Queued Status = iota + 1
	Poisoned
	Completed
	Exceeded
)

func (s Status) ToString() string {
	return [...]string{"QUEUED", "COMPLETED", "POISONED", "EXCEEEDED"}[s-1]
}
