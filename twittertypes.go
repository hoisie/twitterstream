package httpstream

import (
	"bytes"
	//"encoding/json"
	"net/url"
)

type User struct {
	Id                           *Int64Nullable
	Id_str                       StringNullable `json:"id_str"` // "id_str":"608729011",
	Name                         string
	ScreenName                   string         `json:"screen_name"`
	ContributorsEnabled          bool           `json:"contributors_enabled"`
	CreatedAt                    string         `json:"created_at"`
	Description                  StringNullable `json:"description"`
	FavouritesCount              int            `json:"favourites_count"`
	Followerscount               int            `json:"followers_count"`
	Following                    *BoolNullable  // "following":null,
	Friendscount                 int            `json:"friends_count"`
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
	Utc_offset                   *IntNullable   // "utc_offset":null,
	Verified                     bool
	ShowAllInlineMedia           *BoolNullable `json:"show_all_inline_media"`
	RawBytes                     []byte
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
	Coordinates             *Coordinate
	In_reply_to_screen_name StringNullable
	In_reply_to_status_id   *Int64Nullable
	In_reply_to_user_id     *Int64Nullable
	Id                      *Int64Nullable
	Id_str                  string
	Created_at              string
	Retweet_Count           int32
	Retweeted               *BoolNullable
	Possibly_Sensitive      *BoolNullable
	User                    *User
	RawBytes                []byte
	Truncated               *BoolNullable
	Place                   *Place // "place":null,
	//Geo                     string   // deprecated
	//RetweetedStatus         Tweet `json:"retweeted_status"`
}

func (t *Tweet) Urls() []string {
	if len(t.Entities.Urls) > 0 {
		urls := make([]string, 0)
		for _, u := range t.Entities.Urls {
			if len(string(u.Expanded_url)) > 0 {
				if eu, err := url.QueryUnescape(string(u.Expanded_url)); err == nil {
					urls = append(urls, eu)
				}
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

// Return a list of usernames found in the tweet entity mentions
func (t *Tweet) Mentions() []string {
	if len(t.Entities.User_mentions) > 0 {
		users := make([]string, 0)
		for _, m := range t.Entities.User_mentions {
			users = append(users, m.Screen_name)
		}
		return users
	}
	return nil
}

// Create a nullable coordinates, as the data comes across like so:
//    "coordinates":null,
type Coordinate struct {
	Coordinates []float64
	Type        string
}

type Place struct {
	Attributes  interface{}
	Bounding    BoundingBox `json:"bounding_box"`
	Country     string      `json:"country"`
	CountryCode string      `json:"country_code"`
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	PlaceType   string      `json:"place_type"`
	Url         string      `json:"url"`
}

// Location bounding box of coordinates
type BoundingBox struct {
	Coordinates [][][]float64
	Type        string // "Polygon"
}

/*
func (c *Coordinate) UnmarshalJSON(data []byte) error {
	// do we need this, can't we just use pointer?
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
*/
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

// A twitter url
//  "urls":[{"indices":[123,136],"url":"http:\/\/t.co\/a","display_url":null,"expanded_url":null}]
type TwitterUrl struct {
	Url          string
	Expanded_url StringNullable // may be null
	Display_url  StringNullable // may be null if it gets chopped off after t.co because of shortenring
	Indices      []int
}
type Mention struct {
	Screen_name string
	Name        StringNullable // No idea why this could be null, if a username gets mentioned that doesn't exist?
	Id          *Int64Nullable
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
