httpstream was forked from https://github.com/hoisie/twitterstream

A Go http streaming client. http-streaming is most-associated with the twitter stream api.  This client works with twitter, but has also been tested against the data-sift stream as well as http://flowdock.com



This is an example of using the `Twitter stream sample` :

    package main

    import "github.com/araddon/httpstream"

    func main() {
        stream := make(chan []byte)
        done := make(chan bool)
        client := httpstream.NewBasicAuthClient("yourusername", "pwd", func(line []byte) {
            stream <- line
        })
        go func() {
            _ := client.Sample(done)
            for line := range stream {
                println(string(line))
                // heavy lifting like json serialization, etc
            }
        }()
        _ = <- done
    }


There are a few more samples in the Examples folder.

Use a channel instead of func :

        stream := make(chan []byte)
        done := make(chan bool)
        client := httpstream.NewChannelClient("yourusername", "pwd", stream)
        go func() {
            for line := range stream {
                println(string(line))
            }
        }()
        client.Sample(done)
        _ = <- done



For more information about streaming apis

- twitter stream api:  https://dev.twitter.com/docs/streaming-api/methods
- flowdock: https://www.flowdock.com/api
- datasift stream api:  http://dev.datasift.com/docs/streaming-api