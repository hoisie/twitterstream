package httpstream

import (
	"encoding/json"
	//"github.com/bsdf/twitter"
	"log"
	"strconv"
	"testing"
)

var (
	tweets = []string{`{
		"contributors":null,
		"text":"RT @WTRC_UK: Before And After Dinner, A Portrait of Andr\u00e9 Gregory by Cindy Kleine \u2014 Kickstarter http:\/\/t.co\/XoV97Mob via @kickstarter @B ...",
		"entities":{"urls":[{"indices":[96,116],"display_url":"kck.st\/MrvOEX","expanded_url":"http:\/\/kck.st\/MrvOEX","url":"http:\/\/t.co\/XoV97Mob"}],"hashtags":[],"user_mentions":[{"indices":[3,11],"screen_name":"WTRC_UK","id_str":"377562197","name":"WALK THE RED CARPET","id":377562197},{"indices":[121,133],"screen_name":"kickstarter","id_str":"16186995","name":"Kickstarter","id":16186995},{"indices":[134,136],"screen_name":"b","id_str":"11266532","name":"Brian Griffing","id":11266532}]},
		"possibly_sensitive_editable":true,
		"retweeted_status":{
			"contributors":null,"text":"Before And After Dinner, A Portrait of Andr\u00e9 Gregory by Cindy Kleine \u2014 Kickstarter http:\/\/t.co\/XoV97Mob via @kickstarter @BeforeAfterDin m","entities":{"urls":[{"indices":[83,103],"display_url":"kck.st\/MrvOEX","expanded_url":"http:\/\/kck.st\/MrvOEX","url":"http:\/\/t.co\/XoV97Mob"}],"hashtags":[],"user_mentions":[{"indices":[108,120],"screen_name":"kickstarter","id_str":"16186995","name":"Kickstarter","id":16186995},{"indices":[121,136],"screen_name":"BeforeAfterDin","id_str":"608729011","name":"BeforeAndAfterDinner","id":608729011}]},
			"possibly_sensitive_editable":true,"place":null,"retweeted":false,"in_reply_to_status_id":null,"possibly_sensitive":false,
			"in_reply_to_screen_name":null,"in_reply_to_user_id":null,"truncated":false,"source":"web","id_str":"232142697144676352",
			"in_reply_to_status_id_str":null,"favorited":false,"in_reply_to_user_id_str":null,
			"user":{
				"profile_background_tile":false,"friends_count":336,"show_all_inline_media":false,"lang":"en","verified":false,
				"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme9\/bg.gif","time_zone":"London","profile_sidebar_fill_color":"252429","listed_count":1,"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2219235571\/redcarpet_normal.jpg","location":"LONDON","profile_sidebar_border_color":"181A1E","description":"We go to every London film premiere. We stand outside to greet the stars and get their autograph. We Walk the Red Carpet. ",
				"default_profile":false,"follow_request_sent":null,"statuses_count":861,"following":null,"notifications":null,
				"id_str":"377562197","is_translator":false,"profile_use_background_image":true,"screen_name":"WTRC_UK","profile_text_color":"666666","profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme9\/bg.gif","protected":false,
				"default_profile_image":false,"profile_link_color":"2FC2EF","name":"WALK THE RED CARPET","contributors_enabled":false,"geo_enabled":true,"favourites_count":9,"created_at":"Wed Sep 21 19:32:51 +0000 2011","followers_count":149,
				"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2219235571\/redcarpet_normal.jpg","id":377562197,"utc_offset":0,"profile_background_color":"1A1B1F","url":"http:\/\/www.walktheredcarpet.co"
			},
			"retweet_count":1,"id":232142697144676352,"created_at":"Sun Aug 05 15:55:06 +0000 2012",
			"coordinates":null,"geo":null
		},
		"place":null,
		"retweeted":false,
		"in_reply_to_status_id":null,
		"possibly_sensitive":false,
		"in_reply_to_screen_name":null,
		"in_reply_to_user_id":null,
		"truncated":true,
		"source":"web",
		"id_str":"232143626149433344",
		"in_reply_to_status_id_str":null,
		"favorited":false,
		"in_reply_to_user_id_str":null,
		"user":{
			"profile_background_tile":false,
			"friends_count":973,"show_all_inline_media":false,
			"lang":"en","verified":false,
			"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png",
			"time_zone":null,
			"profile_sidebar_fill_color":"DDEEF6",
			"listed_count":9,
			"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2313643667\/lwkkpvh5shev6be81bq9_normal.jpeg",
			"location":"New York, New York",
			"profile_sidebar_border_color":"C0DEED",
			"description":"Kickstarter campaign for film about the life and work of Andr\u00e9 Gregory, visionary theatre director, storyteller and the My Dinner With Andre Guy ",
			"default_profile":true,
			"follow_request_sent":null,
			"statuses_count":999,
			"following":null,
			"notifications":null,
			"id_str":"608729011",
			"is_translator":false,
			"profile_use_background_image":true,
			"screen_name":"BeforeAfterDin",
			"profile_text_color":"333333",
			"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png",
			"protected":false,
			"default_profile_image":false,
			"profile_link_color":"0084B4",
			"name":"BeforeAndAfterDinner",
			"contributors_enabled":false,
			"geo_enabled":false,
			"favourites_count":61,
			"created_at":"Fri Jun 15 03:10:49 +0000 2012",
			"followers_count":383,
			"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2313643667\/lwkkpvh5shev6be81bq9_normal.jpeg",
			"id":608729011,
			"utc_offset":null,
			"profile_background_color":"C0DEED",
			"url":"http:\/\/kck.st\/LfU1xd"
		},
		"retweet_count":1,
		"id":232143626149433344,
		"created_at":"Sun Aug 05 15:58:48 +0000 2012",
		"coordinates":null,
		"geo":null
	}`,
		`{
		"contributors":null,
		"text":"START HERE: Read Your Way Into 25 Amazing Authors\nby @bookriot http:\/\/t.co\/F55jA7H9",
		"entities":{"urls":[{"indices":[63,83],"display_url":"kickstarter.com\/projects\/bookr\u2026","expanded_url":"http:\/\/www.kickstarter.com\/projects\/bookriot\/start-here-read-your-way-into-25-amazing-authors","url":"http:\/\/t.co\/F55jA7H9"}],"hashtags":[],"user_mentions":[{"indices":[53,62],"screen_name":"BookRiot","id_str":"355321621","name":"Book Riot","id":355321621}]},
		"possibly_sensitive_editable":true,
		"place":null,
		"retweeted":false,
		"in_reply_to_status_id":null,
		"possibly_sensitive":false,
		"in_reply_to_screen_name":null,
		"in_reply_to_user_id":null,
		"truncated":false,
		"source":"web",
		"id_str":"232126970564071424",
		"in_reply_to_status_id_str":null,
		"favorited":false,
		"in_reply_to_user_id_str":null,
		"user":{
			"profile_background_tile":false,
			"friends_count":144,
			"show_all_inline_media":false,
			"lang":"en",
			"verified":false,
			"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme8\/bg.gif",
			"time_zone":"Eastern Time (US & Canada)",
			"profile_sidebar_fill_color":"EADEAA",
			"listed_count":4,
			"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2002027587\/1AA86AC4-2047-4AA0-B4AC-1AF3F69E0C31_normal",
			"location":"",
			"profile_sidebar_border_color":"D9B17E",
			"description":"Floccinaucinihilipilification. I just wasted 29 of my characters with one word. It was worth it.         Comics. News. GSM. Ramblings. ",
			"default_profile":false,
			"follow_request_sent":null,
			"statuses_count":2864,
			"following":null,
			"notifications":null,
			"id_str":"28444416",
			"is_translator":false,
			"profile_use_background_image":true,
			"screen_name":"MunsterDeLag",
			"profile_text_color":"333333",
			"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme8\/bg.gif",
			"protected":false,
			"default_profile_image":false,
			"profile_link_color":"9D582E",
			"name":"Brandon .",
			"contributors_enabled":false,
			"geo_enabled":false,
			"favourites_count":47,
			"created_at":"Thu Apr 02 23:14:55 +0000 2009",
			"followers_count":43,
			"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2002027587\/1AA86AC4-2047-4AA0-B4AC-1AF3F69E0C31_normal",
			"id":28444416,"utc_offset":-18000,
			"profile_background_color":"8B542B",
			"url":null
		},
		"retweet_count":0,
		"id":232126970564071424,
		"created_at":"Sun Aug 05 14:52:37 +0000 2012",
		"coordinates":null,
		"geo":null
	}`, `{
	  "contributors": null,
	  "coordinates": null,
	  "created_at": "Sun Aug 05 16:12:50 +0000 2012",
	  "entities": {
	    "hashtags": null,
	    "urls": [
	      {
	        "display_url": "kck.st/y6a9RV",
	        "expanded_url": "http://kck.st/y6a9RV",
	        "indices": [
	          66,
	          86
	        ],
	        "url": "http://t.co/8YpU1gKl"
	      }
	    ],
	    "user_mentions": [
	      {
	        "id": 16186995,
	        "id_str": "16186995",
	        "indices": [
	          91,
	          103
	        ],
	        "name": "Kickstarter",
	        "screen_name": "kickstarter"
	      }
	    ]
	  },
	  "favorited": false,
	  "geo": null,
	  "id": 232147157837287424,
	  "id_str": "232147157837287424",
	  "in_reply_to_screen_name": null,
	  "in_reply_to_status_id": null,
	  "in_reply_to_status_id_str": null,
	  "in_reply_to_user_id": null,
	  "in_reply_to_user_id_str": null,
	  "place": null,
	  "possibly_sensitive": false,
	  "possibly_sensitive_editable": true,
	  "retweet_count": 0,
	  "retweeted": false,
	  "source": "\u003ca href=\"http://twitter.com/tweetbutton\" rel=\"nofollow\"\u003eTweet Button\u003c/a\u003e",
	  "text": "The Order of the Stick Reprint Drive by Rich Burlew — Kickstarter http://t.co/8YpU1gKl via @kickstarter",
	  "truncated": false,
	  "user": {
	    "contributors_enabled": false,
	    "created_at": "Sun Aug 05 11:36:11 +0000 2012",
	    "default_profile": true,
	    "default_profile_image": true,
	    "description": null,
	    "favourites_count": 0,
	    "follow_request_sent": null,
	    "followers_count": 0,
	    "following": null,
	    "friends_count": 0,
	    "geo_enabled": false,
	    "id": 738441984,
	    "id_str": "738441984",
	    "is_translator": false,
	    "lang": "en",
	    "listed_count": 0,
	    "location": null,
	    "name": "SUNNY",
	    "notifications": null,
	    "profile_background_color": "C0DEED",
	    "profile_background_image_url": "http://a0.twimg.com/images/themes/theme1/bg.png",
	    "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme1/bg.png",
	    "profile_background_tile": false,
	    "profile_image_url": "http://a0.twimg.com/sticky/default_profile_images/default_profile_0_normal.png",
	    "profile_image_url_https": "https://si0.twimg.com/sticky/default_profile_images/default_profile_0_normal.png",
	    "profile_link_color": "0084B4",
	    "profile_sidebar_border_color": "C0DEED",
	    "profile_sidebar_fill_color": "DDEEF6",
	    "profile_text_color": "333333",
	    "profile_use_background_image": true,
	    "protected": false,
	    "screen_name": "sunnykhanna1983",
	    "show_all_inline_media": false,
	    "statuses_count": 4,
	    "time_zone": null,
	    "url": null,
	    "utc_offset": null,
	    "verified": false
	  }
	}`, `{
  "contributors": null,
  "coordinates": null,
  "created_at": "Fri Aug 10 17:59:24 +0000 2012",
  "entities": {
    "hashtags": null,
    "urls": [
      {
        "display_url": "kck.st/OYWmRT",
        "expanded_url": "http://kck.st/OYWmRT",
        "indices": [
          86,
          106
        ],
        "url": "http://t.co/QNcHxqT5"
      }
    ],
    "user_mentions": [
      {
        "id": 16186995,
        "id_str": "16186995",
        "indices": [
          111,
          123
        ],
        "name": "Kickstarter",
        "screen_name": "kickstarter"
      }
    ]
  },
  "favorited": false,
  "geo": null,
  "id": 233985915146608641,
  "id_str": "233985915146608641",
  "in_reply_to_screen_name": null,
  "in_reply_to_status_id": null,
  "in_reply_to_status_id_str": null,
  "in_reply_to_user_id": null,
  "in_reply_to_user_id_str": null,
  "place": null,
  "possibly_sensitive": false,
  "possibly_sensitive_editable": true,
  "retweet_count": 0,
  "retweeted": false,
  "source": "\u003ca href=\"http://twitter.com/tweetbutton\" rel=\"nofollow\"\u003eTweet Button\u003c/a\u003e",
  "text": "On Richmond's Wheel: The History of Cycling in Richmond by Thomas Houff — Kickstarter http://t.co/QNcHxqT5 via @kickstarter",
  "truncated": false,
  "user": {
    "contributors_enabled": false,
    "created_at": "Fri Aug 10 17:45:38 +0000 2012",
    "default_profile": true,
    "default_profile_image": false,
    "description": null,
    "favourites_count": 0,
    "follow_request_sent": null,
    "followers_count": 0,
    "following": null,
    "friends_count": 9,
    "geo_enabled": false,
    "id": 749763571,
    "id_str": "749763571",
    "is_translator": false,
    "lang": "en",
    "listed_count": 0,
    "location": null,
    "name": "On Richmond's Wheel",
    "notifications": null,
    "profile_background_color": "C0DEED",
    "profile_background_image_url": "http://a0.twimg.com/images/themes/theme1/bg.png",
    "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme1/bg.png",
    "profile_background_tile": false,
    "profile_image_url": "http://a0.twimg.com/profile_images/2489483541/photo-full_normal.jpg",
    "profile_image_url_https": "https://si0.twimg.com/profile_images/2489483541/photo-full_normal.jpg",
    "profile_link_color": "0084B4",
    "profile_sidebar_border_color": "C0DEED",
    "profile_sidebar_fill_color": "DDEEF6",
    "profile_text_color": "333333",
    "profile_use_background_image": true,
    "protected": false,
    "screen_name": "OnRichmondWheel",
    "show_all_inline_media": false,
    "statuses_count": 1,
    "time_zone": null,
    "url": null,
    "utc_offset": null,
    "verified": false
  }
}`}
	//tweet3 = `{"contributors":null,"coordinates":null,"in_reply_to_user_id":null,"truncated":false,"text":"On Richmond's Wheel: The History of Cycling in Richmond by Thomas Houff \u2014 Kickstarter http:\/\/t.co\/QNcHxqT5 via @kickstarter","entities":{"hashtags":[],"urls":[{"display_url":"kck.st\/OYWmRT","indices":[86,106],"expanded_url":"http:\/\/kck.st\/OYWmRT","url":"http:\/\/t.co\/QNcHxqT5"}],"user_mentions":[{"indices":[111,123],"id_str":"16186995","name":"Kickstarter","screen_name":"kickstarter","id":16186995}]},"place":null,"possibly_sensitive_editable":true,"id_str":"233985915146608641","favorited":false,"geo":null,"source":"\u003Ca href=\"http:\/\/twitter.com\/tweetbutton\" rel=\"nofollow\"\u003ETweet Button\u003C\/a\u003E","possibly_sensitive":false,"retweet_count":0,"in_reply_to_status_id_str":null,"in_reply_to_screen_name":null,"created_at":"Fri Aug 10 17:59:24 +0000 2012","in_reply_to_user_id_str":null,"user":{"show_all_inline_media":false,"lang":"en","friends_count":9,"profile_sidebar_border_color":"C0DEED","location":null,"profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme1\/bg.png","id_str":"749763571","listed_count":0,"profile_use_background_image":true,"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2489483541\/photo-full_normal.jpg","description":null,"follow_request_sent":null,"following":null,"profile_text_color":"333333","default_profile":true,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme1\/bg.png","followers_count":0,"is_translator":false,"time_zone":null,"profile_link_color":"0084B4","protected":false,"created_at":"Fri Aug 10 17:45:38 +0000 2012","profile_background_color":"C0DEED","name":"On Richmond's Wheel","default_profile_image":false,"contributors_enabled":false,"statuses_count":1,"geo_enabled":false,"notifications":null,"profile_background_tile":false,"url":null,"profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2489483541\/photo-full_normal.jpg","screen_name":"OnRichmondWheel","id":749763571,"verified":false,"utc_offset":null,"favourites_count":0,"profile_sidebar_fill_color":"DDEEF6"},"retweeted":false,"in_reply_to_status_id":null,"id":233985915146608641}`
)

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

func TestDecodeTweet1Test(t *testing.T) {
	iv := int64(1.6186995e+07)
	log.Println(strconv.FormatInt(iv, 10))
	for i := 0; i < len(tweets); i++ {
		//for i := 3; i < 4; i++ {
		tw := Tweet{}
		err := json.Unmarshal([]byte(tweets[i]), &tw)
		if err != nil {
			t.Error(err)
			log.Println(tweets[i][0:100])
		}
		log.Println(i, " ", err)
	}
	//prettyJson(tweet3)
	//tw2 := twitter.Tweet{}
	//err = json.Unmarshal([]byte(tweet2), &tw2)
	//log.Println(err)
}
