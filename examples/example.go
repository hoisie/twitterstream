package main

import "twitterstream"

func main() {
    client := twitterstream.NewClient("username", "password")
    err := client.Sample()
    if err != nil {
        println(err.String())
    }
    for {
        tw := <-client.Stream
        println(tw.User.Screen_name, ": ", tw.Text)
    }
}
