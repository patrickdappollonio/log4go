// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

// This log writer sends output to a socket
type SocketLogWriter struct {
	rec chan *LogRecord
	closing bool
    wg *sync.WaitGroup
}

// This is the SocketLogWriter's output method
func (w *SocketLogWriter) LogWrite(rec *LogRecord) {
	if w.closing {
		fmt.Fprintf(os.Stderr, "SocketLogWriter: channel has been closed. Message is [%s]\n", rec.Message)
		return
	}
	w.rec <- rec
}

func (w *SocketLogWriter) Close() {
	w.closing = true
	close(w.rec)
    w.wg.Wait()
}

func NewSocketLogWriter(proto, hostport string) *SocketLogWriter {
	w := &SocketLogWriter{
		rec:  	make(chan *LogRecord, LogBufferLength),
		closing: 	false,
        wg: 	&sync.WaitGroup{},	
	}

	go w.run(proto, hostport)
	return w
}

func (w *SocketLogWriter) run(proto, hostport string) {
	var sock net.Conn

	sock = nil

    w.wg.Add(1)
	defer func() {
		if sock != nil && proto == "tcp" {
			sock.Close()
		}
	    w.wg.Done()
	}()

	for {
		rec, ok := <-w.rec
		if !ok {
			return
		}
		
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
				if sock != nil {
					sock.Close()
					sock = nil
				}
				continue
			}
		}

		_, err = sock.Write(js)
		if err == nil {
			continue
		}

		fmt.Fprintf(os.Stderr, "SocketLogWriter(%s): %v\n", hostport, err)
		sock.Close()
		sock = nil
	}
}
