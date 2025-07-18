package structure

// Go still doesn't have sets
// therefore, this was derived from https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
// no external packages are needed for now, but that could change

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

func (s *Set) Length() int {
	return len(s.list)
}

func (s *Set) Clear() {
	s.list = make(map[string]struct{})
}

// returns the set (HACK: the returned set is mutable and breaks encapsulation)
func (s *Set) GetSet() map[string]struct{} {
	return s.list
}

// returns a new set with values in self but not in otherSet
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
