package main

import (
	"flag"
	"log"
	"os"

	zmq "github.com/alecthomas/gozmq"
	. "github.com/araddon/httpstream"
)

var (
	pwd      *string = flag.String("pwd", "password", "Password")
	user     *string = flag.String("user", "username", "username")
	sink     *string = flag.String("sink", "tcp://*:8888", "Zmq address to connect to:   ")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

	flag.Parse()
	// set the logger and log level
	SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), *logLevel)

	// make a go channel for msgs
	stream := make(chan []byte, 200)
	done := make(chan bool)

	// the stream listener effectively operates in one "thread"
	client := NewBasicAuthClient(*user, *pwd, func(line []byte) {
		stream <- line
	})

	go func() {
		_ = client.Sample(done)

		context, _ := zmq.NewContext()
		socket, _ := context.NewSocket(zmq.PUB) // PUSH?
		socket.Bind(*sink)
		defer func() {
			context.Close()
			socket.Close()
		}()
		for {
			evt := <-stream
			socket.Send(evt, 0)
		}
	}()

	_ = <-done

}
