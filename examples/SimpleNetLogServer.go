package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
	"encoding/json"
	"github.com/daviddengcn/go-colortext"
)

import l4g "github.com/ccpaging/log4go"

type LogRecord struct {
	Level   l4g.Level // The log level
	Created time.Time // The time at which the log message was created (nanoseconds)
	Source  string    // The message source
	Message string    // The log message
}

var (
	port = flag.String("p", "12124", "Port number to listen on")
)

func handleListener(listener *net.UDPConn){
	var timestr string
	var timestrAt int64
	var rec LogRecord

    var buffer [4096]byte
	
	// read into a new buffer
	buflen, addr, err := listener.ReadFrom(buffer[0:])
    if err != nil{
		fmt.Println("[Error] [", addr, "] ", err)
        return
    }

	if buflen <= 0{
		fmt.Println("[Error] [", addr, "] ", "Empty packet")
        return
	}

	bufleft := buffer[:buflen]

	err = json.Unmarshal(bufleft, &rec)
	if err != nil {
		// log to standard output
		fmt.Println("[Error] [", addr, "] ", err)
		fmt.Println(string(buffer[0:]))
		return
	}
		
	switch rec.Level {
		case l4g.CRITICAL:
			ct.ChangeColor(ct.Red, true, ct.White, false)
		case l4g.ERROR:
			ct.ChangeColor(ct.Red, false, 0, false)
		case l4g.WARNING:
			ct.ChangeColor(ct.Yellow, false, 0, false)
		case l4g.INFO:
			ct.ChangeColor(ct.Green, false, 0, false)
		case l4g.DEBUG:
			ct.ChangeColor(ct.Cyan, false, 0, false)
		case l4g.TRACE:
			ct.ChangeColor(ct.Blue, false, 0, false)
		default:
			ct.ChangeColor(ct.None, false, 0, false)
	}
	if at := rec.Created.UnixNano() / 1e9; at != timestrAt {
		timestr, timestrAt = rec.Created.Format("01/02/06 15:04:05"), at
	}
	fmt.Print("[", timestr, "] [", l4g.LevelStrings[rec.Level], "] ", rec.Message)
	ct.ResetColor()
	fmt.Println()
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Erroring out: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	// Bind to the port
	bind, err := net.ResolveUDPAddr("udp4", "0.0.0.0:" + *port)
	checkError(err)

	fmt.Printf("Listening to port %s...\n", *port)
	
	// Create listener
	listener, err := net.ListenUDP("udp", bind)
	checkError(err)

	for {
		handleListener(listener)
	}
}
