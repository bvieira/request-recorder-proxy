package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Trace   = load(ioutil.Discard, "TRACE")
	Info    = load(os.Stdout, "INFO")
	Warning = load(os.Stdout, "WARNING")
	Error   = load(os.Stderr, "ERROR")
)

func load(traceHandle io.Writer, name string) *log.Logger {
	return log.New(traceHandle, name+": ", log.Ldate|log.Ltime|log.Lshortfile)
}
