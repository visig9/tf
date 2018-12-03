// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package readline_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gitlab.com/visig/tf/readline"
	"golang.org/x/sys/unix"
)

func TestIsPipeOnFIFO(t *testing.T) {
	tmpdirPath, err := ioutil.TempDir("", "readline-test-dir")
	defer os.RemoveAll(tmpdirPath)
	if err != nil {
		t.Error(err)
		return
	}

	tmpfilePath := filepath.Join(tmpdirPath, "test-fifo")

	if err := unix.Mkfifo(tmpfilePath, 0666); err != nil {
		t.Error(err, tmpfilePath)
		return
	}

	tmpfile, err := os.OpenFile(tmpfilePath, os.O_RDWR, 0666)
	defer tmpfile.Close()
	if err != nil {
		t.Error(err)
		return
	}

	if out := readline.IsPipe(tmpfile); out != true {
		t.Error("isPipe(tmpfile) != true, want: true")
	}
}
