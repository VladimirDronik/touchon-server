package intset

import (
	"testing"
)

type DecorFunc func(s *IntSet)

func IntSetAdd(x int) DecorFunc {
	return func(s *IntSet) {
		s.Add(x)
	}
}

func IntSetAddAll(values ...int) DecorFunc {
	return func(s *IntSet) {
		s.AddAll(values...)
	}
}

func IntSetUnionWith(t *IntSet) DecorFunc {
	return func(s *IntSet) {
		s.UnionWith(t)
	}
}

func IntSetRemove(x int) DecorFunc {
	return func(s *IntSet) {
		s.Remove(x)
	}
}

func IntSetClear() DecorFunc {
	return func(s *IntSet) {
		s.Clear()
	}
}

func IntSetIntersectWith(t *IntSet) DecorFunc {
	return func(s *IntSet) {
		s.IntersectWith(t)
	}
}

func IntSetDifferenceWith(t *IntSet) DecorFunc {
	return func(s *IntSet) {
		s.DifferenceWith(t)
	}
}

func IntSetSymmetricDifference(t *IntSet) DecorFunc {
	return func(s *IntSet) {
		s.SymmetricDifference(t)
	}
}

func IntSetCopy(t *IntSet) DecorFunc {
	return func(s *IntSet) {
		x := s.Copy()
		*t = *x
	}
}

func NewIntSet(values ...int) *IntSet {
	s := &IntSet{}
	s.AddAll(values...)
	return s
}

func TestIntset(t *testing.T) {
	s := &IntSet{}

	var tests = []struct {
		decorFunc DecorFunc
		result    string
		len       int
		has       map[int]bool
	}{
		{IntSetAddAll(1, 144, 9), "{1 9 144}", 3, map[int]bool{1: true, 0: false}},
		{IntSetCopy(s), "{1 9 144}", 3, map[int]bool{9: true, 144: true}},
		{IntSetUnionWith(NewIntSet(9, 42)), "{1 9 42 144}", 4, map[int]bool{42: true}},
		{IntSetRemove(9), "{1 42 144}", 3, map[int]bool{9: false, 1: true}},
		{IntSetAdd(9), "{1 9 42 144}", 4, map[int]bool{9: true}},
		{IntSetClear(), "{}", 0, map[int]bool{1: false, 9: false, 42: false, 144: false}},
		{IntSetAddAll(1, 2, 3), "{1 2 3}", 3, map[int]bool{2: true, 3: true}},
		{IntSetSymmetricDifference(NewIntSet(1, 2, 4)), "{3 4}", 2, map[int]bool{1: false, 2: false, 3: true, 4: true}},
		{IntSetIntersectWith(NewIntSet(1, 2, 4)), "{4}", 1, map[int]bool{1: false, 2: false, 3: false, 4: true}},
		{IntSetDifferenceWith(NewIntSet(1, 2, 4)), "{}", 0, map[int]bool{4: false}},
	}

	for i, test := range tests {
		test.decorFunc(s)
		if s.String() != test.result {
			t.Errorf("Test №%d: expected %s, got %s\n", i+1, test.result, s.String())
		}
		if s.Len() != test.len {
			t.Errorf("Test №%d: len expected %d, got %d\n", i+1, test.len, s.Len())
		}
		for x, r := range test.has {
			if s.Has(x) != r {
				t.Errorf("Test №%d: has(%d) = %t\n", i+1, x, s.Has(x))
			}
		}
	}
}
