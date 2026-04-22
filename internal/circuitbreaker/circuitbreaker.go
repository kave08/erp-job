package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker is open")

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type Breaker struct {
	mu               sync.RWMutex
	state            State
	failures         int
	lastFailureTime  time.Time
	failureThreshold int
	resetTimeout     time.Duration
	halfOpenMaxCalls int
	halfOpenCalls    int
}

func New(failureThreshold int, resetTimeout time.Duration) *Breaker {
	return &Breaker{
		failureThreshold: failureThreshold,
		resetTimeout:     resetTimeout,
		halfOpenMaxCalls: 1,
		state:            StateClosed,
	}
}

func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.lastFailureTime) > b.resetTimeout {
			b.state = StateHalfOpen
			b.halfOpenCalls = 0
			return true
		}
		return false
	case StateHalfOpen:
		if b.halfOpenCalls < b.halfOpenMaxCalls {
			b.halfOpenCalls++
			return true
		}
		return false
	}
	return false
}

func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateHalfOpen:
		b.state = StateClosed
		b.failures = 0
		b.halfOpenCalls = 0
	case StateClosed:
		b.failures = 0
	}
}

func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	b.lastFailureTime = time.Now()

	switch b.state {
	case StateHalfOpen:
		b.state = StateOpen
	case StateClosed:
		if b.failures >= b.failureThreshold {
			b.state = StateOpen
		}
	}
}

func (b *Breaker) State() State {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}
