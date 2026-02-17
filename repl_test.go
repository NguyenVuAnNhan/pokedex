package main

import "testing"

func TestCleanInput(t *testing.T) {
    cases := []struct {
		input string
		expected []string
	}{
		{"", []string{}},
		{"   ", []string{}},
		{"hello world", []string{"hello", "world"}},
		{"  hello   world  ", []string{"hello", "world"}},
		{"Hello World", []string{"hello", "world"}},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) == %q, expected %q", c.input, actual, c.expected)
			continue
		}
		for i := range actual {	
			if actual[i] != c.expected[i] {
				t.Errorf("cleanInput(%q) == %q, expected %q", c.input, actual, c.expected)
			}
		}
	}
}