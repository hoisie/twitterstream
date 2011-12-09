package main

import (
    "twitterstream"
)

func main() {
    stream := make(chan *twitterstream.Tweet)
    client := twitterstream.NewClient("username", "password")
    //err := client.Track([]string{"bieber,iphone,mac,android,ios,lady gaga,dancing,sick,game,when,why,where,how,who"}, stream)
    err := client.Sample(stream)
    if err != nil {
        println(err.Error())
    }
    for {
        tw := <-stream
        println(tw.User.Screen_name, ": ", tw.Text)
    }
}


