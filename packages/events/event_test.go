package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/events"
)

func TestNewDomainEventAndRoundTrip(t *testing.T) {
	t.Parallel()

	type body struct {
		UserID string `json:"user_id"`
	}

	ev, err := events.NewDomainEvent("identity.user.logged_in", "identity", body{UserID: "u-1"})
	if err != nil {
		t.Fatal(err)
	}
	if ev.ID == "" || ev.Type != "identity.user.logged_in" || ev.Producer != "identity" {
		t.Fatalf("unexpected event: %+v", ev)
	}
	if ev.OccurredAt.Location() != time.UTC {
		t.Fatal("OccurredAt must be UTC")
	}

	raw, err := ev.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := events.UnmarshalDomainEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != ev.ID || got.Type != ev.Type {
		t.Fatalf("round-trip mismatch: %+v vs %+v", got, ev)
	}

	var decoded body
	if err := got.UnmarshalPayload(&decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.UserID != "u-1" {
		t.Fatalf("payload = %+v", decoded)
	}
}

func TestNewDomainEventRequiresTypeAndProducer(t *testing.T) {
	t.Parallel()
	if _, err := events.NewDomainEvent("", "identity", map[string]string{}); err == nil {
		t.Fatal("expected error for empty type")
	}
	if _, err := events.NewDomainEvent("x", "", map[string]string{}); err == nil {
		t.Fatal("expected error for empty producer")
	}
}

func TestUnmarshalDomainEventRejectsIncomplete(t *testing.T) {
	t.Parallel()
	if _, err := events.UnmarshalDomainEvent([]byte(`{"id":"a"}`)); err == nil {
		t.Fatal("expected error")
	}
	if _, err := events.UnmarshalDomainEvent([]byte(`{`)); err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestMarshalPayloadRawMessage(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"ok":true}`)
	got, err := events.MarshalPayload(raw)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"ok":true}` {
		t.Fatalf("got %s", got)
	}
	if _, err := events.MarshalPayload(json.RawMessage(`{`)); err == nil {
		t.Fatal("expected invalid raw json error")
	}
}
