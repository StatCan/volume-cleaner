package main

type Set struct {
	list map[string]struct{}
}

func (s *Set) Has(v string) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set) Add(v string) {
	s.list[v] = struct{}{}
}

func (s *Set) Remove(v string) {
	delete(s.list, v)
}

func (s *Set) Difference(otherSet *Set) *Set {
	newSet := NewSet()

	for v := range s.list {
		if !otherSet.Has(v) {
			newSet.Add(v)
		}
	}

	return newSet
}

func NewSet() *Set {
	s := &Set{}
	s.list = make(map[string]struct{})
	return s
}
