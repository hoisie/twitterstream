package main

import "github.com/hoisie/twitterstream"

func main() {
    stream := make(chan *twitterstream.Tweet)
    client := twitterstream.NewClient("username", "password")
    err := client.Track([]string{"miley"}, stream)
    if err != nil {
        println(err.Error())
    }
    for {
        tw := <-stream
        println(tw.User.Screen_name, ": ", tw.Text)
    }
}
