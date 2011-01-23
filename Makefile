include $(GOROOT)/src/Make.inc

TARG=twitterstream
GOFILES=\
	oauth.go\
	escape.go\
	twitterstream.go\
	types.go\

include $(GOROOT)/src/Make.pkg

format:
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w oauth.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w escape.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w twitterstream.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w types.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w examples/example.go

