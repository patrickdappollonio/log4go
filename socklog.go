// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

// This log writer sends output to a socket
type SocketLogWriter struct {
	rec chan *LogRecord
	sock net.Conn
}

// This is the SocketLogWriter's output method
func (w SocketLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w SocketLogWriter) Close() {
	close(w.rec)
	for i := 10; i > 0 && len(w.rec) > 0; i-- {
		time.Sleep(100 * time.Millisecond)
	}
}

func NewSocketLogWriter(proto, hostport string) *SocketLogWriter {
	w := &SocketLogWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
		sock:     nil,
	}

	go w.run(proto, hostport)
	return w
}

func (w SocketLogWriter) run(proto, hostport string) {
	var err error
	w.sock, err = net.Dial(proto, hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(%v): %s\n", hostport, err)
	}

	defer func() {
		if w.sock != nil {
			w.sock.Close()
		}
	}()

	for rec := range w.rec {
		// Marshall into JSON
		js, err := json.Marshal(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%q): %s\n", hostport, err)
			return
		}

		_, err = w.sock.Write(js)
		if err == nil {
			continue
		}
		
		fmt.Fprintf(os.Stderr, "SocketLogWriter(%q): %s\n", hostport, err)
		if proto == "tcp" {
			return
		}
		
		// proto == "udp", retry
		w.sock, err = net.Dial(proto, hostport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%q): %s\n", hostport, err)
		}
	}
}
