twitterstream is a client for the Twitter Streaming API and Twitter User Stream API.

For more information about this api, visit the [documentation page](http://apiwiki.twitter.com/Streaming-API-Documentation) or [user stream documentation](http://dev.twitter.com/pages/user_streams). 
The Stream connection only supports the `statuses/filter` and the `statuses/sample` method. Other methods require partner accounts at Twitter. 
The User Stream connection support a single function User which populates a FriendList channel and an Event channel.

Here is some example of getting a sample twitter stream and user stream. You need to include a valid twitter screen name and password:

    package main
    
    import (
    	"twitterstream"
    	"fmt"	
    )
    
    func main() {
    	client := twitterstream.NewClient("username", "password")
    	err := client.User()
    	if err != nil {
    	    println(err.String())
    	}
    	for {
    	     select {
    	        case f := <- client.FriendStream:
    	            println("Friends", len(f.Friends))
    	        case e := <- client.EventStream:
    	            println("Event", e.Event)
    	        case tweet := <- client.Stream:
    	           println("Tweet", tweet.User.Screen_name)
    	     }
    	}
    }
    
