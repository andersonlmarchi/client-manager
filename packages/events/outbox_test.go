package events_test

import (
	"errors"
	"testing"
	"time"

	"github.com/andersonlmarchi/client-manager/packages/events"
)

func TestOutboxRecordLifecycle(t *testing.T) {
	t.Parallel()

	ev, err := events.NewDomainEvent("organizations.org.created", "organizations", map[string]string{
		"org_id": "o-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	ev.Aggregate = "organization"
	ev.AggregateID = "o-1"

	rec, err := events.NewOutboxRecord(ev)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Status != events.OutboxPending || rec.EventID != ev.ID || rec.Attempts != 0 {
		t.Fatalf("unexpected record: %+v", rec)
	}

	restored, err := rec.Event()
	if err != nil {
		t.Fatal(err)
	}
	if restored.Type != ev.Type || restored.AggregateID != "o-1" {
		t.Fatalf("restored = %+v", restored)
	}

	if err := rec.MarkFailed(errors.New("webhook timeout")); err != nil {
		t.Fatal(err)
	}
	if rec.Status != events.OutboxFailed || rec.Attempts != 1 || rec.LastError == "" {
		t.Fatalf("after fail: %+v", rec)
	}

	at := time.Date(2026, 8, 11, 12, 0, 0, 0, time.UTC)
	if err := rec.MarkProcessed(at); err != nil {
		t.Fatal(err)
	}
	if rec.Status != events.OutboxProcessed || rec.ProcessedAt == nil || !rec.ProcessedAt.Equal(at) {
		t.Fatalf("after processed: %+v", rec)
	}
	if rec.LastError != "" {
		t.Fatal("LastError should be cleared on success")
	}
}

func TestOutboxNilRecordGuards(t *testing.T) {
	t.Parallel()
	var rec *events.OutboxRecord
	if err := rec.MarkProcessed(time.Now()); err == nil {
		t.Fatal("expected nil guard on MarkProcessed")
	}
	if err := rec.MarkFailed(errors.New("x")); err == nil {
		t.Fatal("expected nil guard on MarkFailed")
	}
}
