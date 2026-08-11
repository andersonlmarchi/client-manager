package events

import (
	"encoding/json"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/shared"
)

// DomainEvent is the cross-service event envelope.
type DomainEvent struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Aggregate   string          `json:"aggregate,omitempty"`
	AggregateID string          `json:"aggregate_id,omitempty"`
	Producer    string          `json:"producer"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Payload     json.RawMessage `json:"payload"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// NewDomainEvent builds an envelope with a new UUID and UTC timestamp.
func NewDomainEvent(eventType, producer string, payload any) (DomainEvent, error) {
	if eventType == "" {
		return DomainEvent{}, shared.NewError(shared.CodeInvalid, "event type is required")
	}
	if producer == "" {
		return DomainEvent{}, shared.NewError(shared.CodeInvalid, "event producer is required")
	}

	raw, err := MarshalPayload(payload)
	if err != nil {
		return DomainEvent{}, err
	}

	id, err := shared.NewUUID()
	if err != nil {
		return DomainEvent{}, shared.Wrap(shared.CodeInternal, "event id", err)
	}

	return DomainEvent{
		ID:         id,
		Type:       eventType,
		Producer:   producer,
		OccurredAt: time.Now().UTC(),
		Payload:    raw,
	}, nil
}

// Marshal encodes the envelope as JSON.
func (e DomainEvent) Marshal() ([]byte, error) {
	if e.Type == "" {
		return nil, shared.NewError(shared.CodeInvalid, "event type is required")
	}
	if e.ID == "" {
		return nil, shared.NewError(shared.CodeInvalid, "event id is required")
	}
	data, err := json.Marshal(e)
	if err != nil {
		return nil, shared.Wrap(shared.CodeInternal, "marshal domain event", err)
	}
	return data, nil
}

// UnmarshalDomainEvent decodes a JSON envelope.
func UnmarshalDomainEvent(data []byte) (DomainEvent, error) {
	var e DomainEvent
	if err := json.Unmarshal(data, &e); err != nil {
		return DomainEvent{}, shared.Wrap(shared.CodeInvalid, "unmarshal domain event", err)
	}
	if e.ID == "" || e.Type == "" || e.Producer == "" {
		return DomainEvent{}, shared.NewError(shared.CodeInvalid, "event id, type and producer are required")
	}
	return e, nil
}

// MarshalPayload encodes a payload value as JSON raw message.
func MarshalPayload(v any) (json.RawMessage, error) {
	if v == nil {
		return json.RawMessage("null"), nil
	}
	if raw, ok := v.(json.RawMessage); ok {
		if !json.Valid(raw) {
			return nil, shared.NewError(shared.CodeInvalid, "payload is not valid JSON")
		}
		return raw, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, shared.Wrap(shared.CodeInvalid, "marshal payload", err)
	}
	return json.RawMessage(data), nil
}

// UnmarshalPayload decodes the event payload into dest.
func (e DomainEvent) UnmarshalPayload(dest any) error {
	if len(e.Payload) == 0 || string(e.Payload) == "null" {
		return shared.NewError(shared.CodeInvalid, "payload is empty")
	}
	if err := json.Unmarshal(e.Payload, dest); err != nil {
		return shared.Wrap(shared.CodeInvalid, "unmarshal payload", err)
	}
	return nil
}
