// Package ratelimit provides an in-memory sliding-window rate limiter
// keyed by a string key (typically client IP or email).
package ratelimit

import (
	"sync"
	"time"
)

// Service defines the rate limiter operations.
type Service interface {
	// Allow checks if a request identified by key is within the allowed rate.
	// Returns true if the request is allowed, false if rate-limited.
	Allow(key string) bool

	// Remaining returns how many more attempts are allowed for the key
	// within the current window.
	Remaining(key string) int

	// ResetAfter returns the duration until the current window expires
	// and the counter resets for the given key.
	ResetAfter(key string) time.Duration
}

// Config holds the rate limiter parameters.
type Config struct {
	MaxAttempts int           // maximum number of requests allowed within the window
	Window      time.Duration // sliding window duration
}

// entry holds the timestamps of requests for a single key.
type entry struct {
	timestamps []time.Time
}

// impl is the concrete implementation of Service using an in-memory map.
type impl struct {
	mu    sync.Mutex
	cfg   Config
	items map[string]*entry
}

// New creates a new rate limiter service with the given config.
func New(cfg Config) Service {
	return &impl{
		cfg:   cfg,
		items: make(map[string]*entry),
	}
}

// Allow checks if the request identified by key is within the rate limit.
func (s *impl) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cfg.MaxAttempts <= 0 {
		// disabled — always allow
		return true
	}

	now := time.Now()
	windowStart := now.Add(-s.cfg.Window)

	e, ok := s.items[key]
	if !ok {
		// first request — create entry and allow
		s.items[key] = &entry{
			timestamps: []time.Time{now},
		}
		return true
	}

	// prune timestamps outside the current window
	valid := e.timestamps[:0]
	for _, t := range e.timestamps {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	e.timestamps = valid

	if len(valid) >= s.cfg.MaxAttempts {
		return false
	}

	e.timestamps = append(e.timestamps, now)

	return true
}

// Remaining returns the number of remaining attempts for the key.
func (s *impl) Remaining(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cfg.MaxAttempts <= 0 {
		return 0
	}

	now := time.Now()
	windowStart := now.Add(-s.cfg.Window)

	e, ok := s.items[key]
	if !ok {
		return s.cfg.MaxAttempts
	}

	valid := e.timestamps[:0]
	for _, t := range e.timestamps {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	e.timestamps = valid

	remaining := s.cfg.MaxAttempts - len(valid)
	if remaining < 0 {
		remaining = 0
	}

	return remaining
}

// ResetAfter returns the duration until the window expires for the key.
func (s *impl) ResetAfter(key string) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cfg.MaxAttempts <= 0 {
		return 0
	}

	e, ok := s.items[key]
	if !ok {
		return 0
	}

	if len(e.timestamps) == 0 {
		return 0
	}

	// the window expires when the oldest timestamp + window passes
	oldest := e.timestamps[0]
	expireAt := oldest.Add(s.cfg.Window)
	remaining := time.Until(expireAt)
	if remaining < 0 {
		return 0
	}

	return remaining
}
