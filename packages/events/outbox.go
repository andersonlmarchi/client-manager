package events

import (
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
)

// OutboxStatus is the processing state of an outbox row.
type OutboxStatus string

const (
	OutboxPending   OutboxStatus = "pending"
	OutboxProcessed OutboxStatus = "processed"
	OutboxFailed    OutboxStatus = "failed"
)

// OutboxRecord is the durable unit stored by producers and read by consumers.
type OutboxRecord struct {
	ID          string       `json:"id"`
	EventID     string       `json:"event_id"`
	EventType   string       `json:"event_type"`
	Producer    string       `json:"producer"`
	Payload     []byte       `json:"payload"`
	Status      OutboxStatus `json:"status"`
	Attempts    int          `json:"attempts"`
	LastError   string       `json:"last_error,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	ProcessedAt *time.Time   `json:"processed_at,omitempty"`
}

// NewOutboxRecord creates a pending outbox row from a domain event.
func NewOutboxRecord(event DomainEvent) (OutboxRecord, error) {
	data, err := event.Marshal()
	if err != nil {
		return OutboxRecord{}, err
	}
	id, err := shared.NewUUID()
	if err != nil {
		return OutboxRecord{}, shared.Wrap(shared.CodeInternal, "outbox id", err)
	}
	return OutboxRecord{
		ID:        id,
		EventID:   event.ID,
		EventType: event.Type,
		Producer:  event.Producer,
		Payload:   data,
		Status:    OutboxPending,
		Attempts:  0,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// MarkProcessed sets status to processed and records the time.
func (r *OutboxRecord) MarkProcessed(at time.Time) error {
	if r == nil {
		return shared.NewError(shared.CodeInvalid, "outbox record is nil")
	}
	ts := at.UTC()
	r.Status = OutboxProcessed
	r.ProcessedAt = &ts
	r.LastError = ""
	return nil
}

// MarkFailed increments attempts and stores the error message.
func (r *OutboxRecord) MarkFailed(err error) error {
	if r == nil {
		return shared.NewError(shared.CodeInvalid, "outbox record is nil")
	}
	r.Status = OutboxFailed
	r.Attempts++
	if err != nil {
		r.LastError = err.Error()
	}
	return nil
}

// Event reconstructs the DomainEvent from the outbox payload.
func (r OutboxRecord) Event() (DomainEvent, error) {
	return UnmarshalDomainEvent(r.Payload)
}
