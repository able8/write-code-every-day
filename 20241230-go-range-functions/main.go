package main

import (
	"fmt"
	"iter"
)

func main() {

	set := New[int]()

	set.Add(1)
	set.Add(2)
	set.Add(3)

	PrintAllElementsPush(set)

	PrintAllElementsPull(set)

	PrintAllElements(set)

	next, stop := iter.Pull(set.All())
	for {
		v, ok := next()
		if !ok {
			break
		}
		fmt.Println("inter.Pull: ", v)
		stop()
	}
}

// Set holds a set of elements.
type Set[E comparable] struct {
	m map[E]struct{}
}

// New returns a new [Set].
func New[E comparable]() *Set[E] {
	return &Set[E]{m: make(map[E]struct{})}
}

// Add adds an element to a set.
func (s *Set[E]) Add(v E) {
	s.m[v] = struct{}{}
}

// Contains reports whether an element is in a set.
func (s *Set[E]) Contains(v E) bool {
	_, ok := s.m[v]
	return ok
}

// Union returns the union of two sets.
func Union[E comparable](s1, s2 *Set[E]) *Set[E] {
	r := New[E]()
	// Note for/range over internal Set field m.
	// We are looping over the maps in s1 and s2.
	for v := range s1.m {
		r.Add(v)
	}
	for v := range s2.m {
		r.Add(v)
	}
	return r
}

func (s *Set[E]) Push(f func(E) bool) {
	for v := range s.m {
		if !f(v) {
			return
		}
	}
}

func PrintAllElementsPush[E comparable](s *Set[E]) {
	s.Push(func(v E) bool {
		fmt.Println(v)
		return true
	})
}

// Pull returns a next function that returns each
// element of s with a bool for whether the value
// is valid. The stop function should be called
// when finished calling the next function.
func (s *Set[E]) Pull() (func() (E, bool), func()) {
	ch := make(chan E)
	stopCh := make(chan bool)

	go func() {
		defer close(ch)
		for v := range s.m {
			select {
			case ch <- v:
			case <-stopCh:
				return
			}
		}
	}()

	next := func() (E, bool) {
		v, ok := <-ch
		return v, ok
	}

	stop := func() {
		close(stopCh)
	}

	return next, stop
}

func PrintAllElementsPull[E comparable](s *Set[E]) {
	next, stop := s.Pull()
	defer stop()
	for v, ok := next(); ok; v, ok = next() {
		fmt.Println(v)
	}
}

// All is an iterator over the elements of s.
func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

func PrintAllElements[E comparable](s *Set[E]) {
	for v := range s.All() {
		fmt.Println(v)
	}
}
