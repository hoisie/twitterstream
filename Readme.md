twitterstream is a client for the Twitter Streaming API. It also has partial support for the User Stream API and the Site streams API.

For the regular streaming API, there's only three methods: `Follow`, `Track`, and `Sample`. 

This is an example of using the `Sample` method:

    package main
    import "twitterstream"

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


This is an example of using the `Track` method:

    package main
    import "twitterstream"

    func main() {
        stream := make(chan *twitterstream.Tweet)
        client := twitterstream.NewClient("username", "password")
        err := client.Track([]string{ "miley"}, stream)
        if err != nil {
            println(err.String())
        }
        for {
            tw := <- stream
            println(tw.User.Screen_name, ": ", tw.Text)
        }
    }

For more information about this API, visit the [documentation page](http://apiwiki.twitter.com/Streaming-API-Documentation). 
