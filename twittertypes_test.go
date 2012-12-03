package httpstream

import (
	"encoding/json"
	//"github.com/bsdf/twitter"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var (
	tweets    = make([]string, 0)
	tweetData []interface{}
)

func init() {
	SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), "debug")
	loadJsonData()
}

func loadJsonData() {
	// load the tweet data
	if jsonb, err := ioutil.ReadFile("data/testdata.json"); err == nil {
		parts := bytes.Split(jsonb, []byte("\n\n"))
		for _, part := range parts {
			tweets = append(tweets, string(bytes.Trim(part, "\n \t\r")))
		}
	}
}
func prettyJson(js string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(js), &m); err == nil {
		if b, er := json.MarshalIndent(m, "", "  "); er == nil {
			log.Println(string(b))
		} else {
			log.Println(er)
		}
	} else {
		log.Println(err)
	}
}

func TestNullableString(t *testing.T) {
	m := make(map[string]StringNullable)
	var valid bool
	js := `{"url":"http:\/\/a0.twimg.com\/images\/themes\/theme14\/bg.gif"}`
	if err := json.Unmarshal([]byte(js), &m); err == nil {
		if string(m["url"]) == "http://a0.twimg.com/images/themes/theme14/bg.gif" {
			valid = true
		}
	}
	//Debug(m)
	if !valid {
		t.Fail()
	}
	m = make(map[string]StringNullable)
	js = `{"url":null}`
	if err := json.Unmarshal([]byte(js), &m); err != nil {
		t.Fail()
	}
}

func TestDecodeTweet1Test(t *testing.T) {
	twlist := make([]Tweet, 0)
	for i := 0; i < len(tweets); i++ {
		//log.Println(tweets[i])
		//for i := 3; i < 4; i++ {
		tw := Tweet{}
		err := json.Unmarshal([]byte(tweets[i]), &tw)
		if err != nil {
			t.Error(err)
			log.Println(tweets[i][0:100])
		}
		log.Println(i, " ", err, tw.Text)
		twlist = append(twlist, tw)
	}
	/*
		twx := twlist[1]
		for _, url := range twx.Urls() {
			Debug(url)
		}
		twx = twlist[1]
		u := twx.Entities.Urls[0]
		log.Println(twx.Urls())
		log.Println(u.Expanded_url)
	*/
	//prettyJson(tweet3)
	//tw2 := twitter.Tweet{}
	//err = json.Unmarshal([]byte(tweet2), &tw2)
	//log.Println(err)
}
