package httpstream

import (
	"bytes"
)

type User struct {
	Lang                         string
	Verified                     bool
	Followers_count              int
	Location                     string
	Screen_name                  string
	Following                    bool
	Friends_count                int
	Profile_background_color     string
	Favourites_count             int
	Description                  string
	Notifications                string
	Profile_text_color           string
	Url                          string
	Time_zone                    string
	Statuses_count               int
	Profile_link_color           string
	Geo_enabled                  bool
	Profile_background_image_url string
	Protected                    bool
	Contributors_enabled         bool
	Profile_sidebar_fill_color   string
	Name                         string
	Profile_background_tile      string
	Created_at                   string
	Profile_image_url            string
	Id                           int64
	Utc_offset                   int
	Profile_sidebar_border_color string
}

type Tweet struct {
	Text                    string
	Entities                Entity
	RawBytes                []byte
	Truncated               bool
	Geo                     string
	In_reply_to_screen_name string
	Favorited               bool
	Source                  string
	Contributors            string
	In_reply_to_status_id   int64
	In_reply_to_user_id     int64
	Id                      int64
	Id_str                  string
	Created_at              string
	Retweet_Count           int32
	Retweeted               bool
	Possibly_Sensitive      bool
	User                    *User
}

func (t *Tweet) Urls() []string {
	if len(t.Entities.Urls) > 0 {
		urls := make([]string,0)
		for _, u := range t.Entities.Urls {
			urls = append(urls, u.Expanded_url)
		}
		return urls
	}
	return nil
}

func (t *Tweet) Hashes() []string {
	if len(t.Entities.Hashtags) > 0 {
		tags := make([]string,0)
		for _, t := range t.Entities.Hashtags {
			tags = append(tags, t.Text)
		}
		return tags
	}
	return nil
}


type SiteStreamMessage struct {
	For_user int64
	Message  Tweet
}

type Event struct {
	Target     User
	Source     User
	Created_at string
	Event      string
}

type Entity struct {
	Hashtags      []Hashtag
	Urls          []TwitterUrl
	User_mentions []Mention
	Media         []Media
}

type Hashtag struct {
	Text    string
	Indices []int
}
type TwitterUrl struct {
	Url          string
	Expanded_url string
	Display_url  string
	Indices      []int
}
type Mention struct {
	Screen_name string
	Name        string
	Id          int64
	Id_str      string
	Indices     []int
}

type Media struct {
	Id              int64
	Id_str          string
	Display_url     string
	Expanded_url    string
	Indices         []int
	Media_url       string
	Media_url_https string
	Url             string
	Type            string
	Screen_name     string
	Sizes           []Sizes
}

type Sizes struct {
	large  Dimensions
	medium Dimensions
	small  Dimensions
	thumb  Dimensions
}

type Dimensions struct {
	w      int
	resize string
	h      int
}

type FriendList struct {
	Friends []int64
}

/*
The twitter stream contains non-tweets (deletes)

{"delete":{"status":{"user_id_str":"36484472","id_str":"191029491823423488","user_id":36484472,"id":191029491823423488}}}
{"delete":{"status":{"id_str":"191184618165256194","id":191184618165256194,"user_id":355665960,"user_id_str":"355665960"}}}
{"delete":{"status":{"id_str":"172129790210482176","id":172129790210482176,"user_id_str":"499324766","user_id":499324766}}}
{"delete":{"status":{"user_id_str":"366839894","user_id":366839894,"id_str":"116974717763719168","id":116974717763719168}}}
{"delete":{"status":{"user_id_str":"382739413","id":191184546841112579,"user_id":382739413,"id_str":"191184546841112579"}}}
{"delete":{"status":{"user_id_str":"388738304","id_str":"123723751366987776","id":123723751366987776,"user_id":388738304}}}
{"delete":{"status":{"user_id_str":"156157535","id_str":"190608148829179907","id":190608148829179907,"user_id":156157535}}}

*/
// a function to filter out the delete messages
func OnlyTweetsFilter(handler func([]byte)) func([]byte) {
	delTw := []byte(`{"delete"`)
	return func(line []byte) {
		if !bytes.HasPrefix(line, delTw) {
			handler(line)
		}
	}
}
