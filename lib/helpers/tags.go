package helpers

import (
	"sort"
	"strings"
)

var abc = map[rune]bool{}

func init() {
	for _, r := range "abcdefghijklmnopqrstuvwxyz0123456789_йцукенгшщзхъфывапролджэёячсмитьбю" {
		abc[r] = true
	}
}

func PrepareTag(tag string) string {
	s := strings.TrimSpace(strings.ToLower(tag))
	s = strings.ReplaceAll(s, "-", "_")

	var r []rune
	for _, c := range s {
		if abc[c] {
			r = append(r, c)
		}
	}

	return string(r)
}

func PrepareTags(tags string, delim string) []string {
	s := strings.Split(tags, delim)
	m := make(map[string]bool, len(s))
	for _, tag := range s {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			m[tag] = true
		}
	}

	if len(m) == 0 {
		return nil
	}

	s = s[:0]
	for tag := range m {
		s = append(s, tag)
	}

	sort.Strings(s)

	return s
}
