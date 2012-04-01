httpstream was forked from https://github.com/hoisie/twitterstream

A Go http streaming client. http-streaming is most-associated with the twitter stream api.  This client works with twitter, but has also been tested against the data-sift stream as well as a stream I am developing for our startup http://lytics.io 


Changes from TwitterStream:

    * No JSON serialization by default. 

    * remove httplib.go, use generic request




This is an example of using the `Twitter stream sample` :

    package main

    import "github.com/araddon/httpstream"

    func main() {
        stream := make(chan []byte)
        client := httpstream.NewClient("yourusername", "pwd", func(line []byte) {
            stream <- line
        })
        _ := client.Sample()  // non-blocking, it starts a go-routine
        for line := range stream {
            println(string(line))
            // heavy lifting like json serialization, etc
        }
    }


There are a few more samples in the Examples folder.



For more information about this API, visit the [twitter documentation page](https://dev.twitter.com/docs/streaming-api/methods). 
