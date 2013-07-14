package main

import (
	"flag"
	"github.com/araddon/httpstream"
	"log"
	"net/url"
	"os"
)

var (
	pwd          *string = flag.String("pwd", "password", "Password")
	user         *string = flag.String("user", "username", "username")
	track        *string = flag.String("track", "", "Twitter terms to track")
	logLevel     *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
	customUrl, _         = url.Parse("http://localhost:6767/stream")
)

func main() {

	flag.Parse()

	// make a go channel for
	stream := make(chan []byte, 200)
	done := make(chan bool)

	// set the logger and log level
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), *logLevel)

	// the stream listener effectively operates in one "thread"
	client := httpstream.NewBasicAuthClient("", "", func(line []byte) {
		//println(string(line))
		stream <- line
	})
	client.MaxWait = 20

	err := client.Connect(customUrl, nil, done)
	if err != nil {
		println(err.Error())
	} else {
		go func() {
			for line := range stream {

				println(string(line))
			}
		}()

		_ = <-done
	}

}
