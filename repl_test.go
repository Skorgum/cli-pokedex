package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello Skorgum ",
			expected: []string{"hello", "skorgum"},
		},
		{
			input:    "  Go  Lang  ",
			expected: []string{"go", "lang"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("for input %q expected length %d, got %d", c.input, len(c.expected), len(actual))
			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("for input %q at index %d expected %q, got %q", c.input, i, c.expected[i], actual[i])
			}
		}
	}
}
