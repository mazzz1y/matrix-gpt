package text

import (
	"testing"
)

func TestReplaceIgnoreCase(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		old      string
		new      string
		n        int
		expected string
	}{
		{
			name:     "old substring exists",
			s:        "Hello World",
			old:      "World",
			new:      "Everyone",
			n:        -1,
			expected: "Hello Everyone",
		},
		{
			name:     "old substring does not exist",
			s:        "Hello World",
			old:      "Universe",
			new:      "Everyone",
			n:        -1,
			expected: "Hello World",
		},
		{
			name:     "old substring exists but different case",
			s:        "Hello World",
			old:      "WORLD",
			new:      "Everyone",
			n:        -1,
			expected: "Hello Everyone",
		},
		{
			name:     "limited replacements",
			s:        "Hello World World",
			old:      "World",
			new:      "Everyone",
			n:        1,
			expected: "Hello Everyone World",
		},
		{
			name:     "more replacements than occurrences",
			s:        "Hello World",
			old:      "World",
			new:      "Everyone",
			n:        3,
			expected: "Hello Everyone",
		},
		{
			name:     "empty string",
			s:        "",
			old:      "World",
			new:      "Everyone",
			n:        -1,
			expected: "",
		},
		{
			name:     "empty old substring",
			s:        "Hello World",
			old:      "",
			new:      "Everyone",
			n:        -1,
			expected: "Hello World",
		},
		{
			name:     "empty new substring",
			s:        "Hello World",
			old:      "World",
			new:      "",
			n:        -1,
			expected: "Hello ",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := ReplaceIgnoreCase(tc.s, tc.old, tc.new, tc.n)
			if result != tc.expected {
				t.Errorf("ReplaceIgnoreCase(%q, %q, %q, %d) = %q; want %q", tc.s, tc.old, tc.new, tc.n, result, tc.expected)
			}
		})
	}
}
