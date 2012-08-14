package httpstream

import (
	"encoding/json"
	//"github.com/bsdf/twitter"
	"log"
	"os"
	"strconv"
	"testing"
)

func init() {
	SetLogger(log.New(os.Stdout, "", log.Ltime|log.Lshortfile), "debug")
}

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
}`, `{
  "contributors": null,
  "coordinates": null,
  "created_at": "Fri Aug 10 18:43:15 +0000 2012",
  "entities": {
    "hashtags": null,
    "media": [
      {
        "display_url": "pic.twitter.com/3ndh87QB",
        "expanded_url": "http://twitter.com/jgmcomics/status/233954362341343232/photo/1",
        "id": 233954362362314752,
        "id_str": "233954362362314752",
        "indices": [
          109,
          129
        ],
        "media_url": "http://p.twimg.com/Az8sQMZCQAAdYDr.jpg",
        "media_url_https": "https://p.twimg.com/Az8sQMZCQAAdYDr.jpg",
        "sizes": {
          "large": {
            "h": 763,
            "resize": "fit",
            "w": 1024
          },
          "medium": {
            "h": 447,
            "resize": "fit",
            "w": 600
          },
          "small": {
            "h": 254,
            "resize": "fit",
            "w": 340
          },
          "thumb": {
            "h": 150,
            "resize": "crop",
            "w": 150
          }
        },
        "source_status_id": 233954362341343232,
        "source_status_id_str": "233954362341343232",
        "type": "photo",
        "url": "http://t.co/3ndh87QB"
      }
    ],
    "urls": [
      {
        "display_url": "kickstarter.com/projects/jgmco…",
        "expanded_url": "http://www.kickstarter.com/projects/jgmcomics/the-mighty-titan",
        "indices": [
          88,
          108
        ],
        "url": "http://t.co/BRJihBG9"
      }
    ],
    "user_mentions": [
      {
        "id": 80970494,
        "id_str": "80970494",
        "indices": [
          3,
          13
        ],
        "name": "Joe Martino",
        "screen_name": "jgmcomics"
      },
      {
        "id": 89914089,
        "id_str": "89914089",
        "indices": [
          54,
          70
        ],
        "name": "Chris Giarrusso",
        "screen_name": "Chris_Giarrusso"
      },
      {
        "id": 339473364,
        "id_str": "339473364",
        "indices": [
          75,
          87
        ],
        "name": "Jerry Ordway",
        "screen_name": "JerryOrdway"
      }
    ]
  },
  "favorited": false,
  "geo": null,
  "id": 233996952881221632,
  "id_str": "233996952881221632",
  "in_reply_to_screen_name": null,
  "in_reply_to_status_id": null,
  "in_reply_to_status_id_str": null,
  "in_reply_to_user_id": null,
  "in_reply_to_user_id_str": null,
  "place": null,
  "possibly_sensitive": false,
  "possibly_sensitive_editable": true,
  "retweet_count": 1,
  "retweeted": false,
  "retweeted_status": {
    "contributors": null,
    "coordinates": null,
    "created_at": "Fri Aug 10 15:54:03 +0000 2012",
    "entities": {
      "hashtags": null,
      "media": [
        {
          "display_url": "pic.twitter.com/3ndh87QB",
          "expanded_url": "http://twitter.com/jgmcomics/status/233954362341343232/photo/1",
          "id": 233954362362314752,
          "id_str": "233954362362314752",
          "indices": [
            94,
            114
          ],
          "media_url": "http://p.twimg.com/Az8sQMZCQAAdYDr.jpg",
          "media_url_https": "https://p.twimg.com/Az8sQMZCQAAdYDr.jpg",
          "sizes": {
            "large": {
              "h": 763,
              "resize": "fit",
              "w": 1024
            },
            "medium": {
              "h": 447,
              "resize": "fit",
              "w": 600
            },
            "small": {
              "h": 254,
              "resize": "fit",
              "w": 340
            },
            "thumb": {
              "h": 150,
              "resize": "crop",
              "w": 150
            }
          },
          "type": "photo",
          "url": "http://t.co/3ndh87QB"
        }
      ],
      "urls": [
        {
          "display_url": "kickstarter.com/projects/jgmco…",
          "expanded_url": "http://www.kickstarter.com/projects/jgmcomics/the-mighty-titan",
          "indices": [
            73,
            93
          ],
          "url": "http://t.co/BRJihBG9"
        }
      ],
      "user_mentions": [
        {
          "id": 8.9914089e+07,
          "id_str": "89914089",
          "indices": [
            39,
            55
          ],
          "name": "Chris Giarrusso",
          "screen_name": "Chris_Giarrusso"
        },
        {
          "id": 339473364,
          "id_str": "339473364",
          "indices": [
            60,
            72
          ],
          "name": "Jerry Ordway",
          "screen_name": "JerryOrdway"
        }
      ]
    },
    "favorited": false,
    "geo": null,
    "id": 233954362341343232,
    "id_str": "233954362341343232",
    "in_reply_to_screen_name": null,
    "in_reply_to_status_id": null,
    "in_reply_to_status_id_str": null,
    "in_reply_to_user_id": null,
    "in_reply_to_user_id_str": null,
    "place": null,
    "possibly_sensitive": false,
    "possibly_sensitive_editable": true,
    "retweet_count": 1,
    "retweeted": false,
    "source": "web",
    "text": "2 great tastes that go great together! @Chris_Giarrusso and @JerryOrdway http://t.co/BRJihBG9 http://t.co/3ndh87QB",
    "truncated": false,
    "user": {
      "contributors_enabled": false,
      "created_at": "Thu Oct 08 23:21:23 +0000 2009",
      "default_profile": false,
      "default_profile_image": false,
      "description": "Tattooed IT guy by day, two time cancer survivor, father of 2 great kids. I am the creator of the comicbook heroes Shadowflame, Ripperman and #TheMightyTitan!",
      "favourites_count": 134,
      "follow_request_sent": null,
      "followers_count": 2173,
      "following": null,
      "friends_count": 2391,
      "geo_enabled": false,
      "id": 80970494,
      "id_str": "80970494",
      "is_translator": false,
      "lang": "en",
      "listed_count": 48,
      "location": "The land of New Jersey. ",
      "name": "Joe Martino",
      "notifications": null,
      "profile_background_color": "1A1B1F",
      "profile_background_image_url": "http://a0.twimg.com/profile_background_images/202682843/bg.jpg",
      "profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/202682843/bg.jpg",
      "profile_background_tile": true,
      "profile_image_url": "http://a0.twimg.com/profile_images/2436119900/zgkr8o1fwx9tmjikbai2_normal.jpeg",
      "profile_image_url_https": "https://si0.twimg.com/profile_images/2436119900/zgkr8o1fwx9tmjikbai2_normal.jpeg",
      "profile_link_color": "FA0000",
      "profile_sidebar_border_color": "000000",
      "profile_sidebar_fill_color": "000000",
      "profile_text_color": "FFF7FF",
      "profile_use_background_image": true,
      "protected": false,
      "screen_name": "jgmcomics",
      "show_all_inline_media": true,
      "statuses_count": 11669,
      "time_zone": "Bogota",
      "url": "http://www.kickstarter.com/projects/jgmcomics/the-mighty-titan",
      "utc_offset": -18000,
      "verified": false
    }
  },
  "source": "\u003ca href=\"http://twitter.com/download/iphone\" rel=\"nofollow\"\u003eTwitter for iPhone\u003c/a\u003e",
  "text": "RT @jgmcomics: 2 great tastes that go great together! @Chris_Giarrusso and @JerryOrdway http://t.co/BRJihBG9 http://t.co/3ndh87QB",
  "truncated": false,
  "user": {
    "contributors_enabled": false,
    "created_at": "Sun Apr 10 15:15:00 +0000 2011",
    "default_profile": false,
    "default_profile_image": false,
    "description": "Writer. Father. Husband. Brother. Dreamer. Pirate. Rock star. Robot. Kung Fu Ninja. http://www.kickstarter.com/projects/mdreynolds/max-spencer-space-inventor",
    "favourites_count": 4,
    "follow_request_sent": null,
    "followers_count": 414,
    "following": null,
    "friends_count": 745,
    "geo_enabled": false,
    "id": 280063014,
    "id_str": "280063014",
    "is_translator": false,
    "lang": "en",
    "listed_count": 19,
    "location": "Lewisburg, TN",
    "name": "Michael Reynolds",
    "notifications": null,
    "profile_background_color": "001329",
    "profile_background_image_url": "http://a0.twimg.com/profile_background_images/263169842/x7a2107434b05a4870c3b9f8fd628562.jpg",
    "profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/263169842/x7a2107434b05a4870c3b9f8fd628562.jpg",
    "profile_background_tile": false,
    "profile_image_url": "http://a0.twimg.com/profile_images/1306772734/5720_1089193364042_1652727140_30291055_1025638_n_normal.jpg",
    "profile_image_url_https": "https://si0.twimg.com/profile_images/1306772734/5720_1089193364042_1652727140_30291055_1025638_n_normal.jpg",
    "profile_link_color": "448668",
    "profile_sidebar_border_color": "F7B565",
    "profile_sidebar_fill_color": "000B17",
    "profile_text_color": "004358",
    "profile_use_background_image": true,
    "protected": false,
    "screen_name": "m_d_reynolds",
    "show_all_inline_media": false,
    "statuses_count": 883,
    "time_zone": "Central Time (US & Canada)",
    "url": null,
    "utc_offset": -21600,
    "verified": false
  }
}`, `{
  "contributors": null,
  "coordinates": null,
  "created_at": "Sat Aug 11 20:54:36 +0000 2012",
  "entities": {
    "hashtags": null,
    "urls": [
      {
        "expanded_url": null,
        "indices": [
          124,
          136
        ],
        "url": "http://t.co/"
      }
    ],
    "user_mentions": [
      {
        "id": 64794579,
        "id_str": "64794579",
        "indices": [
          3,
          17
        ],
        "name": "Inward Facing Girl",
        "screen_name": "melaniebiehle"
      },
      {
        "id": 16186995,
        "id_str": "16186995",
        "indices": [
          64,
          76
        ],
        "name": "Kickstarter",
        "screen_name": "kickstarter"
      },
      {
        "id": 15248224,
        "id_str": "15248224",
        "indices": [
          88,
          95
        ],
        "name": "Natalie Baack",
        "screen_name": "nab717"
      }
    ]
  },
  "favorited": false,
  "geo": null,
  "id": 234392395985334272,
  "id_str": "234392395985334272",
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
  "retweeted_status": {
    "contributors": null,
    "coordinates": null,
    "created_at": "Sat Aug 11 20:53:22 +0000 2012",
    "entities": {
      "hashtags": [
        {
          "indices": [
            126,
            135
          ],
          "text": "vivarium"
        }
      ],
      "urls": [
        {
          "display_url": "kck.st/LVTCGk",
          "expanded_url": "http://kck.st/LVTCGk",
          "indices": [
            105,
            125
          ],
          "url": "http://t.co/SsbgFpSu"
        }
      ],
      "user_mentions": [
        {
          "id": 16186995,
          "id_str": "16186995",
          "indices": [
            45,
            57
          ],
          "name": "Kickstarter",
          "screen_name": "kickstarter"
        },
        {
          "id": 15248224,
          "id_str": "15248224",
          "indices": [
            69,
            76
          ],
          "name": "Natalie Baack",
          "screen_name": "nab717"
        }
      ]
    },
    "favorited": false,
    "geo": null,
    "id": 234392085426483200,
    "id_str": "234392085426483200",
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
    "text": "I just backed Vivarium: A Sci-Fi Thriller on @Kickstarter. My friend @nab717 is producing. Can you help? http://t.co/SsbgFpSu #vivarium",
    "truncated": false,
    "user": {
      "contributors_enabled": false,
      "created_at": "Tue Aug 11 19:18:01 +0000 2009",
      "default_profile": false,
      "default_profile_image": false,
      "description": "Blogger, photographer, graphic designer, social media and content marketing consultant at Inward Facing Girl. ",
      "favourites_count": 39,
      "follow_request_sent": null,
      "followers_count": 1153,
      "following": null,
      "friends_count": 1635,
      "geo_enabled": true,
      "id": 64794579,
      "id_str": "64794579",
      "is_translator": false,
      "lang": "en",
      "listed_count": 66,
      "location": "Seattle, WA",
      "name": "Inward Facing Girl",
      "notifications": null,
      "profile_background_color": "FFFFFF",
      "profile_background_image_url": "http://a0.twimg.com/profile_background_images/538372796/ifg-logo-4-28-2012-twitter-background.png",
      "profile_background_image_url_https": "https://si0.twimg.com/profile_background_images/538372796/ifg-logo-4-28-2012-twitter-background.png",
      "profile_background_tile": false,
      "profile_image_url": "http://a0.twimg.com/profile_images/2277637761/hxddiot07gun6n3yynb4_normal.jpeg",
      "profile_image_url_https": "https://si0.twimg.com/profile_images/2277637761/hxddiot07gun6n3yynb4_normal.jpeg",
      "profile_link_color": "ED4624",
      "profile_sidebar_border_color": "CCD0D0",
      "profile_sidebar_fill_color": "FFFFFF",
      "profile_text_color": "333333",
      "profile_use_background_image": true,
      "protected": false,
      "screen_name": "melaniebiehle",
      "show_all_inline_media": true,
      "statuses_count": 8962,
      "time_zone": "Pacific Time (US & Canada)",
      "url": "http://www.inwardfacinggirl.com",
      "utc_offset": -28800,
      "verified": false
    }
  },
  "source": "web",
  "text": "RT @melaniebiehle: I just backed Vivarium: A Sci-Fi Thriller on @Kickstarter. My friend @nab717 is producing. Can you help? http://t.co/ ...",
  "truncated": true,
  "user": {
    "contributors_enabled": false,
    "created_at": "Thu Jun 26 21:04:00 +0000 2008",
    "default_profile": false,
    "default_profile_image": false,
    "description": "Canadian born. Movie market researcher and analyst. Indie film producer (Vivarium currently live on Kickstarter - entervivarium.com). Echo Park dweller. ",
    "favourites_count": 2,
    "follow_request_sent": null,
    "followers_count": 247,
    "following": null,
    "friends_count": 558,
    "geo_enabled": true,
    "id": 15248224,
    "id_str": "15248224",
    "is_translator": false,
    "lang": "en",
    "listed_count": 1,
    "location": "Los Angeles",
    "name": "Natalie Baack",
    "notifications": null,
    "profile_background_color": "EDECE9",
    "profile_background_image_url": "http://a0.twimg.com/images/themes/theme3/bg.gif",
    "profile_background_image_url_https": "https://si0.twimg.com/images/themes/theme3/bg.gif",
    "profile_background_tile": false,
    "profile_image_url": "http://a0.twimg.com/profile_images/2486263694/eujyqzw08o5hlfvbz1bd_normal.jpeg",
    "profile_image_url_https": "https://si0.twimg.com/profile_images/2486263694/eujyqzw08o5hlfvbz1bd_normal.jpeg",
    "profile_link_color": "088253",
    "profile_sidebar_border_color": "D3D2CF",
    "profile_sidebar_fill_color": "E3E2DE",
    "profile_text_color": "634047",
    "profile_use_background_image": true,
    "protected": false,
    "screen_name": "nab717",
    "show_all_inline_media": false,
    "statuses_count": 781,
    "time_zone": "Pacific Time (US & Canada)",
    "url": "http://entervivarium.com",
    "utc_offset": -28800,
    "verified": false
  }
}`}
	//tweet3 = `{"contributors":null,"coordinates":null,"in_reply_to_user_id":null,"truncated":true,"text":"RT @melaniebiehle: I just backed Vivarium: A Sci-Fi Thriller on @Kickstarter. My friend @nab717 is producing. Can you help? http:\/\/t.co\/ ...","entities":{"hashtags":[],"urls":[{"indices":[124,136],"expanded_url":null,"url":"http:\/\/t.co\/"}],"user_mentions":[{"indices":[3,17],"id_str":"64794579","name":"Inward Facing Girl","screen_name":"melaniebiehle","id":64794579},{"indices":[64,76],"id_str":"16186995","name":"Kickstarter","screen_name":"kickstarter","id":16186995},{"indices":[88,95],"id_str":"15248224","name":"Natalie Baack","screen_name":"nab717","id":15248224}]},"place":null,"possibly_sensitive_editable":true,"id_str":"234392395985334272","retweeted_status":{"contributors":null,"coordinates":null,"in_reply_to_user_id":null,"truncated":false,"text":"I just backed Vivarium: A Sci-Fi Thriller on @Kickstarter. My friend @nab717 is producing. Can you help? http:\/\/t.co\/SsbgFpSu #vivarium","entities":{"hashtags":[{"text":"vivarium","indices":[126,135]}],"urls":[{"display_url":"kck.st\/LVTCGk","indices":[105,125],"expanded_url":"http:\/\/kck.st\/LVTCGk","url":"http:\/\/t.co\/SsbgFpSu"}],"user_mentions":[{"indices":[45,57],"id_str":"16186995","name":"Kickstarter","screen_name":"kickstarter","id":16186995},{"indices":[69,76],"id_str":"15248224","name":"Natalie Baack","screen_name":"nab717","id":15248224}]},"place":null,"possibly_sensitive_editable":true,"id_str":"234392085426483200","favorited":false,"geo":null,"source":"\u003Ca href=\"http:\/\/twitter.com\/tweetbutton\" rel=\"nofollow\"\u003ETweet Button\u003C\/a\u003E","possibly_sensitive":false,"retweet_count":0,"in_reply_to_status_id_str":null,"in_reply_to_screen_name":null,"created_at":"Sat Aug 11 20:53:22 +0000 2012","in_reply_to_user_id_str":null,"user":{"show_all_inline_media":true,"lang":"en","friends_count":1635,"profile_sidebar_border_color":"CCD0D0","location":"Seattle, WA","profile_background_image_url_https":"https:\/\/si0.twimg.com\/profile_background_images\/538372796\/ifg-logo-4-28-2012-twitter-background.png","id_str":"64794579","listed_count":66,"profile_use_background_image":true,"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2277637761\/hxddiot07gun6n3yynb4_normal.jpeg","description":"Blogger, photographer, graphic designer, social media and content marketing consultant at Inward Facing Girl. ","follow_request_sent":null,"following":null,"profile_text_color":"333333","default_profile":false,"profile_background_image_url":"http:\/\/a0.twimg.com\/profile_background_images\/538372796\/ifg-logo-4-28-2012-twitter-background.png","followers_count":1153,"is_translator":false,"time_zone":"Pacific Time (US & Canada)","profile_link_color":"ED4624","protected":false,"created_at":"Tue Aug 11 19:18:01 +0000 2009","profile_background_color":"FFFFFF","name":"Inward Facing Girl","default_profile_image":false,"contributors_enabled":false,"statuses_count":8962,"geo_enabled":true,"notifications":null,"profile_background_tile":false,"url":"http:\/\/www.inwardfacinggirl.com","profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2277637761\/hxddiot07gun6n3yynb4_normal.jpeg","screen_name":"melaniebiehle","id":64794579,"verified":false,"utc_offset":-28800,"favourites_count":39,"profile_sidebar_fill_color":"FFFFFF"},"retweeted":false,"in_reply_to_status_id":null,"id":234392085426483200},"favorited":false,"geo":null,"source":"web","possibly_sensitive":false,"retweet_count":0,"in_reply_to_status_id_str":null,"in_reply_to_screen_name":null,"created_at":"Sat Aug 11 20:54:36 +0000 2012","in_reply_to_user_id_str":null,"user":{"show_all_inline_media":false,"lang":"en","friends_count":558,"profile_sidebar_border_color":"D3D2CF","location":"Los Angeles","profile_background_image_url_https":"https:\/\/si0.twimg.com\/images\/themes\/theme3\/bg.gif","id_str":"15248224","listed_count":1,"profile_use_background_image":true,"profile_image_url_https":"https:\/\/si0.twimg.com\/profile_images\/2486263694\/eujyqzw08o5hlfvbz1bd_normal.jpeg","description":"Canadian born. Movie market researcher and analyst. Indie film producer (Vivarium currently live on Kickstarter - entervivarium.com). Echo Park dweller. ","follow_request_sent":null,"following":null,"profile_text_color":"634047","default_profile":false,"profile_background_image_url":"http:\/\/a0.twimg.com\/images\/themes\/theme3\/bg.gif","followers_count":247,"is_translator":false,"time_zone":"Pacific Time (US & Canada)","profile_link_color":"088253","protected":false,"created_at":"Thu Jun 26 21:04:00 +0000 2008","profile_background_color":"EDECE9","name":"Natalie Baack","default_profile_image":false,"contributors_enabled":false,"statuses_count":781,"geo_enabled":true,"notifications":null,"profile_background_tile":false,"url":"http:\/\/entervivarium.com","profile_image_url":"http:\/\/a0.twimg.com\/profile_images\/2486263694\/eujyqzw08o5hlfvbz1bd_normal.jpeg","screen_name":"nab717","id":15248224,"verified":false,"utc_offset":-28800,"favourites_count":2,"profile_sidebar_fill_color":"E3E2DE"},"retweeted":false,"in_reply_to_status_id":null,"id":234392395985334272}`
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
		twlist = append(twlist, tw)
	}
	twx := twlist[1]
	for _, url := range twx.Urls() {
		Debug(url)
	}
	twx = twlist[1]
	u := twx.Entities.Urls[0]
	log.Println(twx.Urls())
	log.Println(u.Expanded_url)

	//prettyJson(tweet3)
	//tw2 := twitter.Tweet{}
	//err = json.Unmarshal([]byte(tweet2), &tw2)
	//log.Println(err)
}
