// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"time"
	"github.com/daviddengcn/go-colortext"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		switch rec.Level {
			case CRITICAL:
				ct.ChangeColor(ct.Red, true, ct.White, false)
			case ERROR:
				ct.ChangeColor(ct.Red, false, 0, false)
			case WARNING:
				ct.ChangeColor(ct.Yellow, false, 0, false)
			case INFO:
				ct.ChangeColor(ct.Green, false, 0, false)
			case DEBUG:
				ct.ChangeColor(ct.Cyan, false, 0, false)
			case TRACE:
				ct.ChangeColor(ct.Blue, false, 0, false)
			default:
		}
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format("15:04:05 MST 2006/01/02"), at
		}
		fmt.Fprint(out, "[", timestr, "] [", levelStrings[rec.Level], "] [", rec.Source, "] ", rec.Message, "\n")
		ct.ResetColor()
		// fmt.Fprint(out, "\n")
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)
	time.Sleep(100 * time.Millisecond)  // Ugly code, but more faithfully than runtime.Gosched()
}
