package store

import (
	"reflect"
    "sync"

    "github.com/Fabien-Halaby/craft-go-patterns/pkg/singleton"
)

/**
	Store:
		- is a generic thread-safe state container.
**/
type Store[T any] struct {
	mu        sync.RWMutex
	state     T
	listeners map[int]func(T)
	nextID    int
}

// map[reflect.Type]*Instance[Store[T]]
var instances = &sync.Map{}

/**
	GetInstance:
		- returns a singleton Store for type T.
**/
func GetInstance[T any]() *Store[T] {
	key := reflect.TypeOf((*T)(nil)).Elem()

	actual, _ := instances.LoadOrStore(key, singleton.NewInstance(func () *Store[T] {
		return &Store[T] {
			listeners: make(map[int]func(T)),
		}
	}))

	return actual.(*singleton.Instance[*Store[T]]).Get()
}

/**
	GetState:
		- returns a copy of current state (thread-safe).
**/
func (s *Store[T]) GetState() T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

/**
	SetState:
		- updates state and notifies listeners atomically.
**/
func (s *Store[T]) SetState(updates func(*T) ) {
	s.mu.Lock()

	updates(&s.state)
	newState := s.state

	listeners := make([]func(T), 0, len(s.listeners))
	for _, fn := range s.listeners {
		listeners = append(listeners, fn)
	}

	s.mu.Unlock()

	for _, fn := range listeners {
		fn(newState)
	}
}

/**
	Subscribe:
		- registers a state change listener. Returns unsubscribe func.
**/
func (s *Store[T]) Subscribe(fn func (T)) func() {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.listeners[id] = fn
	s.mu.Unlock()

	return func () {
		s.mu.Lock()
		delete(s.listeners, id)
		s.mu.Unlock()
	}
}

/**
	ResetForTest resets the singleton for type T (testing only).
**/
func ResetForTest[T any]() {
	key := reflect.TypeOf((*T)(nil)).Elem()
	instances.Delete(key)
}