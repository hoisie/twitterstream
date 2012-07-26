/*

https://stream.flowdock.com/flows?filter=example/main,example/flow

*/

package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	. "github.com/araddon/httpstream"
)

var (
	token    *string = flag.String("token", "password", "Password")
	flow     *string = flag.String("flow", "yourorg/flow", "Flowdock url path:  organization/flow")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

	flag.Parse()

	// make a go channel for msgs
	stream := make(chan []byte, 200)
	done := make(chan bool)

	flowUrl, _ := url.Parse("https://" + *token + "@stream.flowdock.com/flows?filter=" + *flow)

	// set the logger and log level
	SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	// the stream listener effectively operates in one "thread"
	client := NewClient(func(line []byte) {
		stream <- line
	})
	_ = client.Connect(flowUrl, "", done)

	for {
		evt := <-stream
		Debug(string(evt))
	}
}
