TARG=twitterstream
GOFILES=\
	oauth.go\
	escape.go\
	twitterstream.go\
	types.go\

twitterstream: $(GOFILES)
	go build $(GOFILES)

format:
	gofmt -tabs=false -tabwidth=4 -w oauth.go
	gofmt -tabs=false -tabwidth=4 -w escape.go
	gofmt -tabs=false -tabwidth=4 -w twitterstream.go
	gofmt -tabs=false -tabwidth=4 -w types.go
	gofmt -tabs=false -tabwidth=4 -w examples/example.go

