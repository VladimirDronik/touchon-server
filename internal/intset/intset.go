package intset

import (
	"strconv"
	"strings"
)

func New(size int) *IntSet {
	return &IntSet{
		words: make([]uint64, 0, size/64),
	}
}

// IntSet представляет собой множество небольших неотрицательных
// целых чисел. Нулевое значение представляет пустое множество.
type IntSet struct {
	words []uint64
}

// Has указывает, содержит ли множество неотрицательное значение x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add добавляет неотрицательное значение x в множество.
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll добавляет неотрицательные значения в множество.
func (s *IntSet) AddAll(values ...int) {
	for _, x := range values {
		s.Add(x)
	}
}

// UnionWith делает множество s равным объединению множеств s и t.
func (s *IntSet) UnionWith(t *IntSet) {
	if t == nil {
		return
	}

	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Len возвращает количество элементов множества.
func (s *IntSet) Len() int {
	var bitsCount = func(i uint64) int {
		i = i - ((i >> 1) & 0x5555555555555555)
		i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
		return int(((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101 >> 56)
	}

	sum := 0
	for _, word := range s.words {
		sum += bitsCount(word)
	}
	return sum
}

// Remove удаляет неотрицательное значение x из множества.
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] &^= 1 << bit
}

// Clear удаляет все элементы множества.
func (s *IntSet) Clear() {
	s.words = s.words[:0]
}

// Copy создает новое множество и копирует в него эл-ты множества.
func (s *IntSet) Copy() *IntSet {
	var x = new(IntSet)
	x.words = make([]uint64, len(s.words))
	copy(x.words, s.words)
	return x
}

// IntersectWith оставляет только те эл-ты, которые есть в обоих множествах
func (s *IntSet) IntersectWith(t *IntSet) {
	if t == nil {
		return
	}

	length := len(s.words)
	if len(t.words) < length {
		length = len(t.words)
	}
	s.words = s.words[:length] // удаляем старшие эл-ты

	for i := 0; i < length; i++ {
		s.words[i] &= t.words[i]
	}
}

// DifferenceWith из множества удаляет все эл-ты,
// имеющиеся во втором множестве (разность множеств).
func (s *IntSet) DifferenceWith(t *IntSet) {
	if t == nil {
		return
	}

	length := len(s.words)
	if len(t.words) < length {
		length = len(t.words)
	}
	s.words = s.words[:length] // удаляем старшие эл-ты

	for i := 0; i < length; i++ {
		s.words[i] &^= t.words[i]
	}
}

// SymmetricDifference - содержит эл-ты, имеющиеся
// в одном из множеств, но не в обоих одновременно.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	if t == nil {
		return
	}

	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String возвращает множество как строку вида "{1 2 3}".
func (s *IntSet) String() string {
	elems := s.Elems()
	sElems := make([]string, 0, len(elems))
	for _, n := range elems {
		sElems = append(sElems, strconv.Itoa(n))
	}
	return "{" + strings.Join(sElems, " ") + "}"
}

func (s *IntSet) Elems() (r []int) {
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				r = append(r, 64*i+j)
			}
		}
	}

	return
}
