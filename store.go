package ratelimiter

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Limiters centrally maintain a map of rate limiters.
type Limiters struct {
	store map[string]*rate.Limiter
	mutex sync.Mutex
}

// NewStore initialises a in memory map for storing rate limiters.
func NewStore() *Limiters {
	return &Limiters{
		store: make(map[string]*rate.Limiter),
	}
}

// Register adds a new limiter to the limiter store with the given key.
// It returns an error if given key already exists.
func (l *Limiters) Register(key string, rps int, burst int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if _, ok := l.store[key]; ok {
		return errors.New("rate limiter key already exists")
	}
	l.store[key] = rate.NewLimiter(rate.Every(time.Second/time.Duration(rps)), burst)
	return nil
}

// Get retrieves an existing limiter.
// It returns an error if the given key has not been created.
func (l *Limiters) Get(key string) (*rate.Limiter, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	val, ok := l.store[key]
	if ok {
		return val, nil
	}
	return nil, errors.New("rate limiter key does not exist")
}
