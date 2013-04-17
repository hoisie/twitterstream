package main

import (
	"flag"
	"github.com/araddon/httpstream"
	"log"
	"os"
	"strings"
)

var (
	pwd      *string = flag.String("pwd", "password", "Password")
	user     *string = flag.String("user", "username", "username")
	track    *string = flag.String("track", "", "Twitter terms to track")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

	flag.Parse()
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), *logLevel)

	stream := make(chan []byte, 1000)
	done := make(chan bool)

	client := httpstream.NewBasicAuthClient(*user, *pwd, httpstream.OnlyTweetsFilter(func(line []byte) {
		stream <- line
	}))

	keywords := strings.Split("android,golang,zeromq,javascript", ",")
	err := client.Filter([]int64{14230524, 783214}, keywords,[]string{"en"}, false, done)
	if err != nil {
		httpstream.Log(httpstream.ERROR, err.Error())
	} else {

		go func() {
			ct := 0
			for tw := range stream {
				println(string(tw))
				ct++
				if ct > 10 {
					done <- true
				}
			}
		}()
		_ = <-done
	}

}
