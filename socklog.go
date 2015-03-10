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
		rec:  	make(chan *LogRecord, LogBufferLength),
	}

	go w.run(proto, hostport)
	return w
}

func NewHostSocketLogWriter(proto, hostport string) *SocketLogWriter {
	w := &SocketLogWriter{
		rec:  	make(chan *LogRecord, LogBufferLength),
	}

	if proto == "tcp" {
		go w.runTCPServer(proto, hostport)
	} else {
		go w.runUDPServer(proto, hostport)
	}
	return w
}

func (w *SocketLogWriter) run(proto, hostport string) {
	var sock net.Conn

	sock = nil

	defer func() {
		if sock != nil && proto == "tcp" {
			sock.Close()
		}
	}()

	for rec := range w.rec {
		// Marshall into JSON
		js, err := json.Marshal(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			return
		}

		if sock == nil {
			sock, err = net.Dial(proto, hostport)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
				continue
			}
		}

		_, err = sock.Write(js)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			continue
		}
	}
}

func (w *SocketLogWriter) runUDPServer(proto, hostport string) {
	var sock net.Conn

	sock = nil

	defer func() {
		if sock != nil && proto == "tcp" {
			sock.Close()
		}
	}()

	for rec := range w.rec {
		// Marshall into JSON
		js, err := json.Marshal(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			return
		}

		if sock == nil {
			sock, err = net.Dial(proto, hostport)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
				continue
			}
		}

		_, err = sock.Write(js)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			continue
		}
	}
}

func (w *SocketLogWriter) runTCPServer(proto, hostport string) {
	var sock net.Conn

	sock = nil

	defer func() {
		if sock != nil && proto == "tcp" {
			sock.Close()
		}
	}()

	for rec := range w.rec {
		// Marshall into JSON
		js, err := json.Marshal(rec)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			return
		}

		if sock == nil {
			sock, err = net.Dial(proto, hostport)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
				continue
			}
		}

		_, err = sock.Write(js)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
			continue
		}
	}
}

