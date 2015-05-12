// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"sync"
	"github.com/daviddengcn/go-colortext"
)

var stdout io.Writer = os.Stdout

// This is the standard writer that prints to standard output.
type ConsoleLogWriter struct {
	rec chan *LogRecord
	closing bool
    wg *sync.WaitGroup
	color bool
	format string
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (c *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	if c.closing {
		fmt.Fprintf(os.Stderr, "ConsoleLogWriter: channel has been closed. Message is [%s]\n", rec.Message)
		return
	}
	c.rec <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (c *ConsoleLogWriter) Close() {
	if c.closing {
		return
	}
	c.closing = true
	close(c.rec)
    c.wg.Wait()
}

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() *ConsoleLogWriter {
	c := &ConsoleLogWriter{
		rec:  	make(chan *LogRecord, LogBufferLength),
		closing: 	false,
        wg: 	&sync.WaitGroup{},	
		color:	true,
		format: "[%T %D] [%L] (%S) %M",
	}

    c.wg.Add(1)
	go c.run(stdout)
	return c
}

// Must be called before the first log message is written.
func (c *ConsoleLogWriter) SetColor(color bool) *ConsoleLogWriter {
	c.color = color
	return c
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (c *ConsoleLogWriter) SetFormat(format string) *ConsoleLogWriter {
	c.format = format
	return c
}

func (c *ConsoleLogWriter) run(out io.Writer) {
    defer c.wg.Done()

	for {
		rec, ok := <-c.rec
		if !ok {
			return
		}
		if c.color {
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
					ct.ChangeColor(ct.Magenta, false, 0, false)
				case TRACE:
					ct.ChangeColor(ct.Cyan, false, 0, false)
				default:
			}
		}
		fmt.Fprint(out, FormatLogRecord(c.format, rec))
		if c.color {
			ct.ResetColor()
		}
	}
}

