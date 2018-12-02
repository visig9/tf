package readline

import (
	"bufio"
	"os"
)

// IsPipe test target file is a named pipe or not. usually using on os.Stdin.
func IsPipe(f *os.File) bool {
	stat, _ := f.Stat()
	if stat.Mode()&os.ModeNamedPipe != 0 {
		return true
	}

	return false
}

// Line return single line from bufio.Reader.
//
// If error != nil mean reader still have some lines can read.
func Line(br *bufio.Reader) (string, error) {
	line, isPrefix, err := br.ReadLine()
	if err != nil {
		return "", err
	}

	fullline := line
	for isPrefix {
		line, isPrefix, err = br.ReadLine()
		if err != nil {
			return "", err
		}
		fullline = append(fullline, line...)
	}

	return string(fullline), nil
}

// Channel run a goroutine to retrive all line string from target file
// line by line. User can get those lines within returned channel.
func Channel(file *os.File) chan string {
	ch := make(chan string)
	br := bufio.NewReader(file)

	go func() {
		for {
			line, err := Line(br)

			if err == nil {
				ch <- line
			} else {
				close(ch)
				break
			}
		}
	}()

	return ch
}
