include $(GOROOT)/src/Make.$(GOARCH)

TARG=twitterstream
GOFILES=\
	twitterstream.go\

include $(GOROOT)/src/Make.pkg

format:
	gofmt -spaces=true -tabindent=false -tabwidth=4 -w twitterstream.go

