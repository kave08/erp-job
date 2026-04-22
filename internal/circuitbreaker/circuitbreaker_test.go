package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestBreakerStartsClosed(t *testing.T) {
	b := New(3, time.Second)
	if !b.Allow() {
		t.Fatal("expected breaker to allow requests when closed")
	}
}

func TestBreakerOpensAfterFailures(t *testing.T) {
	b := New(2, time.Hour)

	b.RecordFailure()
	b.RecordFailure()

	if b.State() != StateOpen {
		t.Fatalf("expected state open, got %d", b.State())
	}
	if b.Allow() {
		t.Fatal("expected breaker to block requests when open")
	}
}

func TestBreakerHalfOpenAfterTimeout(t *testing.T) {
	b := New(1, 50*time.Millisecond)
	b.RecordFailure()

	if b.State() != StateOpen {
		t.Fatal("expected state open")
	}

	time.Sleep(100 * time.Millisecond)

	if !b.Allow() {
		t.Fatal("expected breaker to allow one request in half-open state")
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected state half-open, got %d", b.State())
	}
}

func TestBreakerClosesOnHalfOpenSuccess(t *testing.T) {
	b := New(1, 50*time.Millisecond)
	b.RecordFailure()

	time.Sleep(100 * time.Millisecond)
	b.Allow()
	b.RecordSuccess()

	if b.State() != StateClosed {
		t.Fatalf("expected state closed, got %d", b.State())
	}
	if !b.Allow() {
		t.Fatal("expected breaker to allow requests after closing")
	}
}

func TestBreakerReOpensOnHalfOpenFailure(t *testing.T) {
	b := New(1, 50*time.Millisecond)
	b.RecordFailure()

	time.Sleep(100 * time.Millisecond)
	b.Allow()
	b.RecordFailure()

	if b.State() != StateOpen {
		t.Fatalf("expected state open, got %d", b.State())
	}
}

func TestBreakerErrOpen(t *testing.T) {
	if !errors.Is(ErrOpen, ErrOpen) {
		t.Fatal("ErrOpen should match itself")
	}
}
