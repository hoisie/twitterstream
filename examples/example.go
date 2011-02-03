package main

import "twitterstream"

func main() {
    stream := make(chan *twitterstream.Tweet)
    client := twitterstream.NewClient("username", "password")
    err := client.Track([]string{"miley"}, stream)
    if err != nil {
        println(err.String())
    }
    for {
        tw := <-stream
        println(tw.User.Screen_name, ": ", tw.Text)
    }
}
