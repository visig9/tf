package logger

import (
	"log"
	"os"
)

// Err is a logger for error message.
var Err = log.New(os.Stderr, os.Args[0]+": ", 0)

// Std is a logger for std output.
var Std = log.New(os.Stdout, "", 0)
