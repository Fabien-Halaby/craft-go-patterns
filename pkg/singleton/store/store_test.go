package store

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestState struct {
	Counter int
	Message string
}

type AnotherState struct {
	Flag bool
}

func TestGenericScore(t *testing.T) {
	ResetForTest[TestState]()

	s := GetInstance[TestState]()
	
	state := s.GetState()
	assert.Equal(t, 0, state.Counter)

	s.SetState(func (st *TestState) {
		st.Counter = 42
		st.Message = "Hello"
	})

	state = s.GetState()
	assert.Equal(t, 42, state.Counter)
	assert.Equal(t, "Hello", state.Message)
}

func TestSingletonPerType(t *testing.T) {
	ResetForTest[TestState]()
	ResetForTest[AnotherState]()

	s1 := GetInstance[TestState]()
	s2 := GetInstance[TestState]()
	s3 := GetInstance[AnotherState]()

	assert.Same(t, s1, s2)

	assert.NotSame(t, s1, s3)
}

func TestConcurrentAccess(t *testing.T) {
	ResetForTest[TestState]()
	
	s := GetInstance[TestState]()
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.SetState(func(st *TestState) {
				st.Counter = n
			})
		}(i)
	}
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.GetState()
		}()
	}
	
	wg.Wait()
}