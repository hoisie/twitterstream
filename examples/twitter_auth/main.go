package main

import (
	"flag"
	"fmt"
	oauth "github.com/araddon/goauth"
	u "github.com/araddon/gou"
	"io/ioutil"
)

/*
Usage:

go clean && go run

./twitter_auth  --ck=MY_APP_CONSUMER_KEY --cs=MY_APP_SECRET

*/
var (
	ck        *string = flag.String("ck", "", "Consumer Key")
	cs        *string = flag.String("cs", "", "Consumer Secret")
	goauthcon *oauth.OAuthConsumer
)

func main() {
	flag.Parse()
	goauthcon = &oauth.OAuthConsumer{
		Service:          "twitter",
		RequestTokenURL:  "https://api.twitter.com/oauth/request_token",
		AccessTokenURL:   "https://api.twitter.com/oauth/access_token",
		AuthorizationURL: "https://api.twitter.com/oauth/authorize",
		ConsumerKey:      *ck,
		ConsumerSecret:   *cs,
		CallBackURL:      "oob",
	}

	s, rt, err := goauthcon.GetRequestAuthorizationURL()
	if err != nil {
		fmt.Println(err)
		return
	}
	var pin string

	fmt.Printf("Open %s In your browser.\n Allow access and then enter the PIN number\n", s)
	fmt.Printf("PIN Number: ")
	fmt.Scanln(&pin)

	at := goauthcon.GetAccessToken(rt.Token, pin)
	fmt.Printf("\n\n\ttoken=%s  secret=%s \n\n\n", at.Token, at.Secret)
	r, err := goauthcon.Get("https://api.twitter.com/1.1/account/verify_credentials.json", nil, at)

	if err != nil {
		fmt.Println(err)
		return
	}
	if r != nil && r.Body != nil && r.StatusCode == 200 {
		userData, _ := ioutil.ReadAll(r.Body)
		if len(userData) > 0 {
			jh := u.NewJsonHelper(userData)
			fmt.Println(string(jh.PrettyJson()))
		}
	} else {
		fmt.Println(r)
		fmt.Println("Twitter verify failed")
	}

}
