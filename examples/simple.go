package main

import (
	"flag"
	"github.com/araddon/httpstream"
)

var pwd *string = flag.String("pwd", "password", "Password")
var user *string = flag.String("user", "username", "username")
var track *string = flag.String("track", "", "Twitter terms to track")

func main() {

	flag.Parse()
	stream := make(chan []byte)
	//stream := make(chan []byte,1000) make a buffered queue channel instead of a blocking one

	// the stream listener effectively operates in one "thread"
	client := httpstream.NewClient(*user, *pwd, func(line []byte) {
		stream <- line
	})
	//err := client.Track([]string{"bieber,iphone,mac,android,ios,lady gaga,dancing,sick,game,when,why,where,how,who"}, stream)
	err := client.Sample()
	if err != nil {
		println(err.Error())
	} else {

		// while this operates in a different "thread(s)" (if more than one proc)
		for {
			tw := <-stream
			println(string(tw))
			// heavy lifting

		}
	}

}
