twitterstream is a client for the Twitter Streaming API. 

For more information about this api, visit the [documentation page](http://apiwiki.twitter.com/Streaming-API-Documentation). 
It only supports the `statuses/filter` and the `statuses/sample` method. Other methods require partner accounts at Twitter. 

Here is some example of getting a sample twitter stream. You need to include a valid twitter screen name and password:

	package main

	import "twitterstream"

	func main() {
		client := twitterstream.NewClient("username", "password")
		err := client.Sample()
		if err != nil {
			println(err.String())
		}
		for {
			tw := <- client.Stream
			println(tw.User.Screen_name, ": ", tw.Text)
		}
	}
