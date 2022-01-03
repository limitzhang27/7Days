package consistenthash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	h := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// 2，4，6，12，14，16，22，24，26
	h.Add("6", "4", "2")

	testCase := map[string]string{
		"2":  "2",
		"13": "4",
		"21": "2",
		"27": "2",
	}

	for k, v := range testCase {
		if h.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 8, 18, 28
	h.Add("8")

	testCase["27"] = "8"
	for k, v := range testCase {
		if h.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}
