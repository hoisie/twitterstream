httpstream was forked from https://github.com/hoisie/twitterstream

It is now (in forked version) geared mostly towards use as a client for Http Streaming API that forwards on messages via messaging to another consumer.   

For the regular streaming API, there's only three methods: `Follow`, `Track`, and `Sample`. 

This is an example of using the `Sample` method:

    package main
    import "github.com/araddon/httpstream"

    func main() {
        stream := make(chan *twitterstream.Tweet)
        client := twitterstream.NewClient("username", "password")
        err := client.Sample(stream)
        if err != nil {
            println(err.String())
        }
        for {
            tw := <- stream
            println(tw.User.Screen_name, ": ", tw.Text)
        }
    }



For more information about this API, visit the [documentation page](http://apiwiki.twitter.com/Streaming-API-Documentation). 
