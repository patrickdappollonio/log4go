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
type SocketLogWriter chan *LogRecord

// This is the SocketLogWriter's output method
func (w SocketLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

func (w SocketLogWriter) Close() {
	close(w)
	for i := 10; i > 0 && len(w) > 0; i-- {
		time.Sleep(100 * time.Millisecond)
	}
}

func NewSocketLogWriter(proto, hostport string) SocketLogWriter {
	sock, err := net.Dial(proto, hostport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewSocketLogWriter(%s): %v\n", hostport, err)
		return nil
	}

	w := SocketLogWriter(make(chan *LogRecord, LogBufferLength))

	go func() {
		defer func() {
			if sock != nil && proto == "tcp" {
				sock.Close()
			}
		}()

		for rec := range w {
			// Marshall into JSON
			js, err := json.Marshal(rec)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
				return
			}

			_, err = sock.Write(js)
			if err != nil {
				fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
				return
			}
		}
	}()

	return w
}

