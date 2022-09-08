package main

type ValidNeighbours struct {
	top, right, bottom, left Set
}

type WaveCollapseRules map[int]ValidNeighbours

func NewWaveCollapseRules(numTiles int) WaveCollapseRules {
	return make(map[int]ValidNeighbours, numTiles)
}

type Set map[int]struct{}

func NewSet() Set {
	return make(Set)
}

func (s Set) Values() []int {
	vv := make([]int, 0, len(s))
	for value := range s {
		vv = append(vv, value)
	}
	return vv
}
func (s Set) Contains(value int) bool {
	_, ok := s[value]
	return ok
}
func (s Set) Copy() Set {
	n := make(Set, len(s))
	for value := range s {
		n.Add(value)
	}
	return n
}
func (s Set) Add(value int) {
	s[value] = struct{}{}
}
func (s Set) Remove(value int) {
	delete(s, value)
}
func (s Set) Intersection(other Set) Set {
	n := NewSet()
	for value := range s {
		if other.Contains(value) {
			n.Add(value)
		}
	}
	return n
}
