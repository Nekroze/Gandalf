package pathing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexp(t *testing.T) {
	cases := []struct {
		blob     string
		regex    string
		target   string
		err      bool
		expected []string
	}{
		{``, ``, ``, true, []string{}},
		{
			blob:     `https://example.com:9000`,
			regex:    `(?P<scheme>[a-zA-Z]+)://(?P<domain>[.-0-9_a-zA-Z]+)(:(?P<port>[0-9]+))?`,
			target:   `scheme`,
			err:      false,
			expected: []string{`https`},
		},
		{
			blob:     `https://example.com:9000`,
			regex:    `(?P<scheme>[a-zA-Z]+)://(?P<domain>[.-0-9_a-zA-Z]+)(:(?P<port>[0-9]+))?`,
			target:   `domain`,
			err:      false,
			expected: []string{`example.com`},
		},
		{
			blob:     `https://example.com:9000`,
			regex:    `(?P<scheme>[a-zA-Z]+)://(?P<domain>[.-0-9_a-zA-Z]+)(:(?P<port>[0-9]+))?`,
			target:   `port`,
			err:      false,
			expected: []string{`9000`},
		},
		{
			blob:     `https://example.com`,
			regex:    `(?P<scheme>[a-zA-Z]+)://(?P<domain>[.-0-9_a-zA-Z]+)(:(?P<port>[0-9]+))?`,
			target:   `port`,
			err:      false,
			expected: []string{},
		},
		{
			blob:     `https://example.com`,
			regex:    `(?P<scheme>[a-zA-Z]+)://(?P<domain>[.-0-9_a-zA-Z]+):(?P<port>[0-9]+)`,
			target:   `port`,
			err:      true,
			expected: []string{},
		},
	}
	for i, tt := range cases {
		t.Run(fmt.Sprintf("Case %d - %s", i, tt.target), func(st *testing.T) {
			rc := NewRegexpExtractor(tt.regex)
			result, err := rc.Extractor(tt.blob, tt.target)
			assert.Equal(st, tt.expected, result)
			errChecker := assert.NoError
			if tt.err {
				errChecker = assert.Error
			}
			errChecker(t, err)
		})
	}
}
