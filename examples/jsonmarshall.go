package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/araddon/httpstream"
	"log"
	"os"
	//"reflect"
	"strings"
)

var (
	pwd      *string = flag.String("pwd", "password", "Twitter Password")
	user     *string = flag.String("user", "username", "Twitter username")
	track    *string = flag.String("track", "", "Twitter terms to track")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func printPretty(tweet *httpstream.Tweet) {
	b, err := json.MarshalIndent(tweet, " ", "   ")
	if err == nil && tweet.Place != nil {
		println(string(b))
		//log.Println(tweet.Urls())
	}
}
func printPrettyBytes(line []byte) {
	var tweet map[string]interface{}
	json.Unmarshal(line, &tweet)
	b, err := json.MarshalIndent(tweet, " ", "   ")
	if err == nil {
		if co, ok := tweet["place"]; ok && co != nil {
			println(string(b))
		}
	}
}

// Since there are multiple types besides tweets sent in the http stream
// determine which it is in order to serialize correctly
func HandleLine(th int, line []byte) {
	switch {
	case bytes.HasPrefix(line, []byte(`{"event":`)):
		var event httpstream.Event
		json.Unmarshal(line, &event)
	case bytes.HasPrefix(line, []byte(`{"friends":`)):
		var friends httpstream.FriendList
		json.Unmarshal(line, &friends)
	default:
		//printPrettyBytes(line)
		tweet := httpstream.Tweet{}
		json.Unmarshal(line, &tweet)
		printPretty(&tweet)
		//if tweet.Coordinates != nil {
		//println(th, " ", tweet.User.Screen_name, ": ", tweet.Text)
		//	printPretty(tweet)
		//}
	}
}

type Msg struct {
	Line []byte
}

func main() {

	var err error
	flag.Parse()
	// set the logger and log level
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	stream := make(chan []byte)
	done := make(chan bool)

	client := httpstream.NewBasicAuthClient(*user, *pwd, func(line []byte) {
		stream <- line
	})
	//err := client.Track([]string{"bieber,iphone,mac,android,ios,lady gaga,dancing,sick,game,when,why,where,how,who"}, stream)
	// this opens a go routine that is effectively thread 1
	if len(*track) > 0 {
		err = client.Filter(nil, strings.Split(*track, ","), true, done)
	} else {
		err = client.Sample(done)
	}
	if err != nil {
		println(err.Error())
	}
	// 2nd thread
	go func() {
		for {
			line := <-stream
			HandleLine(1, line)
		}
	}()
	// 3rd thread
	for {
		line := <-stream
		HandleLine(2, line)
	}
}
