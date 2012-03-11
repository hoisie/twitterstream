httpstream was forked from https://github.com/hoisie/twitterstream

It is now (in forked version) gearing towards generic http streaming (datasift, foursquare, etc).  Processing of any kind is moved out of the thread that is doing the http connection.  

Changes from TwitterStream:

    * Biggest difference is how it handles each stream item:  https://github.com/araddon/httpstream/blob/c11036066b49469f835905370dbd7d39b5aa3c69/stream.go#L157   all handling a passed in handler now

    * new metrics:  This is necessary because twitter stream limits out, so you need to start new clients if you start to approach the (undocumented) max/minute.

    * remove httplib.go, use generic request

    * planning on support for datasift and foursquare



This is an example of using the `Twitter Sample` method:

    package main

    import "github.com/araddon/httpstream"

    func main() {
        stream := make(chan []byte)
        client := httpstream.NewClient(*user, *pwd, func(line []byte) {
            stream <- line
        })
        err := client.Sample()
        if err != nil {
            println(err.Error())
        } else {

            // while this operates in a different "thread(s)" (if more than one proc)
            for {
                line := <-stream
                println(string(line))
                // heavy lifting like json serialization, etc

            }
        }
    }



For more information about this API, visit the [twitter documentation page](https://dev.twitter.com/docs/streaming-api/methods). 
