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
const DefaultTimeFormat string = "15:04:05 MST 2006/01/02"

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

func NewConsoleLogWriter() ConsoleLogWriter {
	return NewColorConsoleLogWriter(true, DefaultTimeFormat)
}

// This creates a new ConsoleLogWriter
func NewColorConsoleLogWriter(color bool, timeformat string) ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout, color, timeformat)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer, color bool, timeformat string) {
	var timestr string
	var timestrAt int64

	for rec := range w {
		if color {
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
		}
		if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
			timestr, timestrAt = rec.Created.Format(timeformat), at
		}
		fmt.Fprint(out, "[", timestr, "] [", levelStrings[rec.Level], "] [", rec.Source, "] ", rec.Message)
		if color {
			ct.ResetColor()
		}
		fmt.Fprint(out, "\n")
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
	for i := 10; i > 0 && len(w) > 0; i-- {
		time.Sleep(100 * time.Millisecond)
	}
}

