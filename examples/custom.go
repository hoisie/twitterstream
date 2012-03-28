
package main

import (
	"flag"
	"github.com/araddon/httpstream"
	"log"
	"net/url"
	"os"
)

var pwd *string = flag.String("pwd", "password", "Password")
var user *string = flag.String("user", "username", "username")
var track *string = flag.String("track", "", "Twitter terms to track")
var logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")

var customUrl, _ = url.Parse("http://localhost:6767/stream")

func main() {

	flag.Parse()

	// make a go channel for 
	stream := make(chan []byte)

	// set the logger and log level
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	// the stream listener effectively operates in one "thread"
	client := httpstream.NewBasicAuthClient("", "", func(line []byte) {
		println(string(line))
	})

	err := client.Connect(customUrl, "")
	if err != nil {
		println(err.Error())
	} 
	for {
		tw := <-stream
		println(tw)
	}
}
