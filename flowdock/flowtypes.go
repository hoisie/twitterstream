package flowdock


import (
)

/*
{"event":"message","tags":[],"uuid":"SSlxf45-xxe","id":111,"flow":"example:main",
	"content":"more","sent":1333222117502,"app":"chat","attachments":[],"user":"19256"}


-- email message:
{"event":"mail","tags":[],"uuid":null,"id":122,"flow":"example:main",
	"content":{
		"subject":"msg",
		"replyTo":[],
		"to":[{"name":null,"address":"main@example.flowdock.com"}],
		"contentType":"text\/html; charset=ISO-8859-1",
		"from":[{"name":"aaron aaron","address":"email@email.com"}],
		"bcc":[],
		"content":"from email",
		"sender":null,
		"cc":[]
	},
	"sent":1333222317815,"app":"influx","attachments":[],
	"user":"0"}

There is a periodic heartbeat per user it looks like

{"event":"activity.user","tags":[],"uuid":null,"id":130,
	"flow":"example:main",
	"content":{"last_activity":1333222272444},
	"sent":1333222525743,"app":null,"attachments":[],"user":"19256"}


*/

type struct Event {
	// "mail", "activity.user",
	Event       string   
	Tags        []string
	Uuid        string
	Id          int
	Flow        string
	Content     interface{}
}