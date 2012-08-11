package httpstream

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// details about all the nullable types http://code.google.com/p/go/issues/detail?id=2540

type Int64Nullable int64

func (i Int64Nullable) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		//ParseInt(s string, base int, bitSize int) (i int64, err error)
		if in, err := strconv.ParseInt(string(data), 10, 64); err == nil {
			i = Int64Nullable(in)
		}

	}
	return nil
}

type IntNullable int

func (i IntNullable) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		//ParseInt(s string, base int, bitSize int) (i int64, err error)
		if in, err := strconv.ParseInt(string(data), 10, 32); err == nil {
			i = IntNullable(in)
		}

	}
	return nil
}

type StringNullable string

func (s StringNullable) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		s = StringNullable(string(data))
	}
	return nil
}

type BoolNullable bool

func (b BoolNullable) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		//ParseBool(str string) (value bool, err error)
		if bo, err := strconv.ParseBool(string(data)); err == nil {
			b = BoolNullable(bo)
		}
	}
	return nil
}

type User struct {
	Id                           int64
	Name                         string
	ScreenName                   string `json:"screen_name"`
	ContributorsEnabled          bool   `json:"contributors_enabled"`
	CreatedAt                    string `json:"created_at"`
	Description                  StringNullable
	FavouritesCount              int          `json:"favourites_count"`
	Followerscount               int          `json:"followers_count"`
	Following                    BoolNullable // "following":null,
	Friendscount                 int          `json:"friends_count"`
	Geo_enabled                  bool
	Lang                         string
	Location                     StringNullable
	Listed_count                 int            `json:"listed_count"`
	Notifications                StringNullable //"notifications":null,
	Profile_text_color           string
	Profile_link_color           string
	Profile_background_image_url string
	Profile_background_color     string
	Profile_sidebar_fill_color   string
	Profile_image_url            string
	Profile_sidebar_border_color string
	Profile_background_tile      bool
	Protected                    bool
	Statuses_Count               int `json:"statuses_count"`
	Time_zone                    StringNullable
	Url                          StringNullable // "url":null
	Utc_offset                   IntNullable    // "utc_offset":null,
	Verified                     bool
	//"show_all_inline_media":false,
	//"default_profile":false,
	//"follow_request_sent":null,
	//"is_translator":false,
	//"profile_use_background_image":true,
	//"default_profile_image":false,
}

type Tweet struct {
	Text                    string
	Entities                Entity
	Favorited               bool
	Source                  string
	Contributors            []Contributor
	Coordinates             Coordinate
	In_reply_to_screen_name StringNullable
	In_reply_to_status_id   Int64Nullable
	In_reply_to_user_id     Int64Nullable
	Id                      int64
	Id_str                  string
	Created_at              string
	Retweet_Count           int32
	Retweeted               bool
	Possibly_Sensitive      bool
	User                    *User
	RawBytes                []byte
	Truncated               bool
	//Geo                     string   // deprecated
	//Place                  // "place":null,
	//RetweetedStatus         Tweet `json:"retweeted_status"`
}

func (t *Tweet) Urls() []string {
	if len(t.Entities.Urls) > 0 {
		urls := make([]string, 0)
		for _, u := range t.Entities.Urls {
			if len(u.Expanded_url) > 0 {
				urls = append(urls, string(u.Expanded_url))
			}
		}
		return urls
	}
	return nil
}

func (t *Tweet) Hashes() []string {
	if len(t.Entities.Hashtags) > 0 {
		tags := make([]string, 0)
		for _, t := range t.Entities.Hashtags {
			tags = append(tags, t.Text)
		}
		return tags
	}
	return nil
}

// Create a nullable coordinates, as the data comes across like so:  
//    "coordinates":null,
type Coordinate struct {
	Coordinates []float64
	Type        string
}

func (c *Coordinate) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		m := make(map[string]interface{})
		if err := json.Unmarshal(data, m); err == nil {
			if co, ok := m["coordinates"]; ok {
				if cof, ok := co.([]float64); ok {
					c.Coordinates = cof
				}
			}
			if ty, ok := m["type"]; ok {
				if tys, ok := ty.(string); ok {
					c.Type = tys
				}
			}
		}
	}
	return nil
}

type Contributor struct {
	Id          int64
	Id_str      string
	Screen_name string
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
	Expanded_url StringNullable
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
	Sizes           Sizes
}

type Sizes struct {
	Large  Dimensions
	Medium Dimensions
	Small  Dimensions
	Thumb  Dimensions
}

type Dimensions struct {
	W      int
	Resize string
	H      int
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
