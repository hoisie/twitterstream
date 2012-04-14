package main

import (
	"flag"
	"github.com/araddon/httpstream"
	"log"
	"os"
)

var (
	pwd      *string = flag.String("pwd", "password", "Password")
	user     *string = flag.String("user", "username", "username")
	track    *string = flag.String("track", "", "Twitter terms to track")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

	flag.Parse()
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	// make a go channel for sending from listener to processor
	// we buffer it, to help ensure we aren't backing up twitter or else they cut us off
	stream := make(chan []byte, 1000)
	done := make(chan bool)

	// the stream listener effectively operates in one "thread"/goroutine
	// as the httpstream Client processes inside a go routine it opens
	// That includes the handler func we pass in here
	client := httpstream.NewBasicAuthClient(*user, *pwd, httpstream.OnlyTweetsFilter(func(line []byte) {
		stream <- line
		// although you can do heavy lifting here, it means you are doing all
		// your work in the same thread as the http streaming/listener
		// by using a go channel, you can send the work to a 
		// different thread/goroutine
	}))

	//err := client.Track([]string{"eat,iphone,mac,android,ios,burger"}, stream)
	err := client.Sample(done)
	if err != nil {
		httpstream.Log(httpstream.ERROR, err.Error())
	} else {

		go func() {
			// while this could be in a different "thread(s)"
			ct := 0
			for tw := range stream {
				println(string(tw))
				// heavy lifting
				ct++
				if ct > 10 {
					os.Exit(0)
				}
			}
		}()
		_ = <-done
	}

}
