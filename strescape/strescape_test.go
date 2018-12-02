package strescape_test

import (
	"testing"

	"gitlab.com/visig/tf/strescape"
)

func TestSingleQuote(t *testing.T) {
	cases := []struct {
		str  string
		want string
	}{
		{`I'm`, `I\'m`},
		{`a'b''`, `a\'b\'\'`},
	}

	for _, c := range cases {
		if ans := strescape.SingleQuote(c.str); ans != c.want {
			t.Errorf("cu.EscapeSingleQuote(%q) != %q",
				c.want, ans)
		}
	}
}
