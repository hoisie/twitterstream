include $(GOROOT)/src/Make.inc

TARG=twitterstream
GOFILES=\
	twitterstream.go\
	types.go\

include $(GOROOT)/src/Make.pkg

format:
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w twitterstream.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w types.go
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w examples/example.go

