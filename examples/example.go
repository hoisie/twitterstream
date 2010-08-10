package main

import (
	"twitterstream"
	"fmt"	
)

func StreamListen(client *twitterstream.Client) {
	for {
		tw := <- client.TweetStream()
		println("Tweet:", tw.User.Screen_name, ": ", tw.Text)
	}
}

func UserStreamListen(client *twitterstream.UserClient) {
	for {
        ue := <- client.EventStream()
        println("User Stream event:", ue.Event, "| from :", ue.Source.Screen_name)
    }
}

func FriendListStreamListen(client *twitterstream.UserClient) {
	for {
        fl := <- client.FriendListStream()
        fmt.Printf("Friend list: %v\n", fl.Friends)
    }
}

func main() {
		
        userClient := twitterstream.NewClient("username", "password")
        err := userClient.User()
	    if err != nil {
	    	println(err.String())
	    }
	    // Background listening for user events
	    go UserStreamListen(userClient) 
	    // Background listening for friendlist 
	    go FriendListStreamListen(userClient)
	    client := twitterstream.NewClient("username", "password")
        err = client.Sample()
        if err != nil {
                println(err.String())
        }
        StreamListen(client)
}
