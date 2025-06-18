package shared

import "sync"

// SafeSet represents a collection of unique elements.
type SafeSet[T comparable] struct {
	m map[T]struct{}
	sync.RWMutex
}

// NewSafeSet creates a new empty SafeSet.
func NewSafeSet[T comparable]() *SafeSet[T] {
	return &SafeSet[T]{m: make(map[T]struct{})}
}

// Add adds an element to the set.
func (s *SafeSet[T]) Add(value T) {
	s.Lock()
	defer s.Unlock()
	s.m[value] = struct{}{}
}

// Remove removes an element from the set.
func (s *SafeSet[T]) Remove(value T) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, value)
}

// Contains checks if the set contains an element.
func (s *SafeSet[T]) Contains(value T) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[value]
	return ok
}

// Size returns the number of elements in the set.
func (s *SafeSet[T]) Size() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

// Intersection returns the intersection of two sets.
func (s *SafeSet[T]) Intersection(other *SafeSet[T]) *SafeSet[T] {
	result := NewSafeSet[T]()
	for k := range s.m {
		if _, ok := other.m[k]; ok {
			result.Add(k)
		}
	}

	return result
}

// Union returns the union of two sets.
func (s *SafeSet[T]) Union(other *SafeSet[T]) *SafeSet[T] {
	s.RLock()
	defer s.RUnlock()
	other.RLock()
	defer other.RUnlock()
	result := NewSafeSet[T]()
	for k := range s.m {
		result.Add(k)
	}

	for k := range other.m {
		if _, ok := s.m[k]; !ok {
			result.Add(k)
		}
	}

	return result
}

func (s *SafeSet[T]) Reset() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[T]struct{})
}

func (s *SafeSet[T]) Range(f func(value T) bool) {
	s.RLock()
	defer s.RUnlock()
	for k := range s.m {
		if !f(k) {
			break
		}
	}
}
