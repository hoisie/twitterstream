package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/araddon/httpstream"
	"log"
)

var pwd *string = flag.String("pwd", "password", "Twitter Password")
var user *string = flag.String("user", "username", "Twitter username")
var track *string = flag.String("track", "", "Twitter terms to track")
var logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")


func HandleLine(th int, line []byte) {
	switch {
	case bytes.HasPrefix(line, []byte(`{"event":`)):
		var event httpstream.Event
		json.Unmarshal(line, &event)
	case bytes.HasPrefix(line, []byte(`{"friends":`)):
		var friends httpstream.FriendList
		json.Unmarshal(line, &friends)
	default:
		tweet := httpstream.Tweet{}
		json.Unmarshal(line, &tweet)
		if tweet.User != nil {
			println(th, " ", tweet.User.Screen_name, ": ", tweet.Text)
		}
	}
}

type Msg struct {
	Line []byte
}

func main() {

	flag.Parse()
	// set the logger and log level
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	stream := make(chan []byte)
	client := httpstream.NewBasicAuthClient(*user, *pwd, func(line []byte) {
		stream <- line
	})
	//err := client.Track([]string{"bieber,iphone,mac,android,ios,lady gaga,dancing,sick,game,when,why,where,how,who"}, stream)
	// this opens a go routine that is effectively thread 1
	err := client.Sample()
	if err != nil {
		println(err.Error())
	}
	// 2nd thread
	go func() {
		for {
			line := <-stream
			println()
			HandleLine(1, line)
		}
	}()
	// 3rd thread
	for {
		line := <-stream
		HandleLine(2, line)
	}
}
