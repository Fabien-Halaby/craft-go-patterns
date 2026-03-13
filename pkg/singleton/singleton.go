package singleton

import (
	"sync"
)

/**
	Instance:
		- represents a singleton instance constructor.
**/
type Instance[T any] struct {
	once	 sync.Once
	instance T
	factory  func() T
}

/**
	NewInstance:
		- creates a new singleton container with the given factory.
**/
func NewInstance[T any](factory func() T) *Instance[T] {
	return &Instance[T]{
		factory: factory,
	}
}

/**
	Get:
		- returns the singleton instance, creating it if necessary.
**/
func (s *Instance[T]) Get() T {
	s.once.Do(func () {
		s.instance = s.factory()
	})

	return s.instance
}

/**
	Reset:
		- clears the instance (useful for testing).
**/
func (s *Instance[T]) Reset() {
    s.once = sync.Once{}
    var zero T
    s.instance = zero
}