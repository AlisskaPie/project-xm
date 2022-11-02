package domain

import (
	"context"

	"github.com/google/uuid"
)

// CompanyEvent is data structure to be sent to Kafka or like.
// Spec: On each mutating operation,
// a JSON formatted event must be produced to a service bus (Kafka, RabbitMQ etc.)
type CompanyEvent struct {
	Action EventActionType
	ID     uuid.UUID
	State  Company
}

type EventActionType string

// List of EventActionType values
const (
	InsertEventActionType EventActionType = "insert"
	UpdateEventActionType EventActionType = "update"
	DeleteEventActionType EventActionType = "delete"
)

// CompanyEventSender is an interface for service bus.
type CompanyEventSender interface {
	Send(ctx context.Context, event CompanyEvent) error
}
