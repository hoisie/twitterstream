package main

// twitter oauth

import (
	"flag"
	oauth "github.com/akrennmair/goauth"
	"github.com/araddon/httpstream"
	"log"
	"os"
	"strings"
)

var (
	user     *string = flag.String("user", "", "twitter username")
	ck       *string = flag.String("ck", "", "Consumer Key")
	cs       *string = flag.String("cs", "", "Consumer Secret")
	ot       *string = flag.String("ot", "", "Oauth Token")
	osec     *string = flag.String("os", "", "OAuthTokenSecret")
	logLevel *string = flag.String("logging", "debug", "Which log level: [debug,info,warn,error,fatal]")
)

func main() {

	flag.Parse()
	httpstream.SetLogger(log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile), *logLevel)

	// make a go channel for sending from listener to processor
	// we buffer it, to help ensure we aren't backing up twitter or else they cut us off
	stream := make(chan []byte, 1000)
	done := make(chan bool)

	httpstream.OauthCon = &oauth.OAuthConsumer{
		Service:          "twitter",
		RequestTokenURL:  "http://twitter.com/oauth/request_token",
		AccessTokenURL:   "http://twitter.com/oauth/access_token",
		AuthorizationURL: "http://twitter.com/oauth/authorize",
		ConsumerKey:      *ck,
		ConsumerSecret:   *cs,
		CallBackURL:      "oob",
		UserAgent:        "go/httpstream",
	}

	//at := goauthcon.GetAccessToken(rt.Token, pin)
	at := oauth.AccessToken{Id: "",
		Token:    *ot,
		Secret:   *osec,
		UserRef:  *user,
		Verifier: "",
		Service:  "twitter",
	}
	// the stream listener effectively operates in one "thread"/goroutine
	// as the httpstream Client processes inside a go routine it opens
	// That includes the handler func we pass in here
	client := httpstream.NewOAuthClient(&at, httpstream.OnlyTweetsFilter(func(line []byte) {
		stream <- line
		// although you can do heavy lifting here, it means you are doing all
		// your work in the same thread as the http streaming/listener
		// by using a go channel, you can send the work to a 
		// different thread/goroutine
	}))

	keywords := strings.Split("android,golang,zeromq,javascript", ",")
	err := client.Filter([]int64{14230524, 783214}, keywords, false, done)
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
					done <- true
				}
			}
		}()
		_ = <-done
	}

}
