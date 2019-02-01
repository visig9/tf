package readline_test

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/visig9/tf/readline"
)

var testlines = []string{
	"apple",
	"orange",
	"juice is great!",
}

func TestIsPipe(t *testing.T) {
	if readline.IsPipe(os.Stdin) == true {
		t.Error("isPipe(os.Stdin) == true, want: false")
	}
}

func TestLine(t *testing.T) {
	content := strings.Join(testlines, "\n")
	reader := bufio.NewReader(strings.NewReader(content))

	for _, expected := range testlines {
		if line, err := readline.Line(reader); err == nil {
			if line != expected {
				t.Errorf("%q != %q", line, expected)
			}
		}
	}
}

func TestChannel(t *testing.T) {
	var testlines = []string{
		"apple",
		"orange",
		"juice is great!",
	}

	content := strings.Join(testlines, "\n")

	tmpfile, err := ioutil.TempFile("", "readline-test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	ch := readline.Channel(tmpfile)

	for _, expected := range testlines {
		if line, ok := <-ch; ok {
			if line != expected {
				t.Errorf("%q != %q", line, expected)
			}
		}
	}
}
