package twitterstream

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"
)

var followUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/filter.json")
var trackUrl, _ = url.Parse("http://stream.twitter.com/1/statuses/filter.json")
var sampleUrl, _ = url.Parse("http://stream.twitter.com/1/statuses/sample.json")
var userUrl, _ = url.Parse("http://userstream.twitter.com/2/user.json")
var siteStreamUrl, _ = url.Parse("https://betastream.twitter.com/2b/site.json")

var retryTimeout time.Duration = 5e9

type streamConn struct {
	clientConn   *httputil.ClientConn
	url          *url.URL
	stream       chan *Tweet
	eventStream  chan *Event
	friendStream chan *FriendList
	authData     string
	postData     string
	stale        bool
}

func (conn *streamConn) Close() {
	// Just mark the connection as stale, and let the connect() handler close after a read
	conn.stale = true
}

func (conn *streamConn) connect() (*http.Response, error) {
	if conn.stale {
		return nil, errors.New("Stale connection")
	}
	var tcpConn net.Conn
	var err error
	if proxy := os.Getenv("HTTP_PROXY"); len(proxy) > 0 {
		proxy_url, _ := url.Parse(proxy)
		tcpConn, err = net.Dial("tcp", proxy_url.Host)
	} else {
		tcpConn, err = net.Dial("tcp", conn.url.Host+":443")
	}
	if err != nil {
		return nil, err
	}
	cf := &tls.Config{Rand: rand.Reader, Time: time.Now}
	ssl := tls.Client(tcpConn, cf)

	conn.clientConn = httputil.NewClientConn(ssl, nil)

	var req http.Request
	req.URL = conn.url
	req.Method = "GET"
	req.Header = http.Header{}
	req.Header.Set("Authorization", "Basic "+conn.authData)

	if conn.postData != "" {
		req.Method = "POST"
		req.Body = nopCloser{bytes.NewBufferString(conn.postData)}
		req.ContentLength = int64(len(conn.postData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := conn.clientConn.Do(&req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (conn *streamConn) readStream(resp *http.Response) {

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)
	for {
		//we've been closed
		if conn.stale {
			conn.clientConn.Close()
			break
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if conn.stale {
				continue
			}

			//try reconnecting
			resp, err := conn.connect()
			if err != nil {
				println(err.Error())
				time.Sleep(retryTimeout)
				continue
			}
			if resp.StatusCode != 200 {
				time.Sleep(retryTimeout)
				continue
			}

			reader = bufio.NewReader(resp.Body)
			continue
		}
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
		}
		switch {
		case bytes.HasPrefix(line, []byte(`{"event":`)):
			if conn.eventStream != nil {
				var event Event
				json.Unmarshal(line, &event)
				conn.eventStream <- &event
			}
		case bytes.HasPrefix(line, []byte(`{"friends":`)):
			if conn.friendStream != nil {
				var friends FriendList
				json.Unmarshal(line, &friends)
				conn.friendStream <- &friends
			}
		default:
			if conn.stream != nil {
				var tweet Tweet
				json.Unmarshal(line, &tweet)
				if tweet.Id != 0 {
					conn.stream <- &tweet
				}
			}
		}
	}
}

func encodedAuth(user, pwd string) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(user + ":" + pwd))
	encoder.Close()
	return buf.String()
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

type Client struct {
	Username     string
	Password     string
	stream       chan *Tweet
	eventStream  chan *Event
	friendStream chan *FriendList
	conn         *streamConn
}

func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
	}
}

func (c *Client) connect(url_ *url.URL, body string) (err error) {
	if c.Username == "" || c.Password == "" {
		return errors.New("The username or password is invalid")
	}

	var resp *http.Response
	//initialize the new stream
	var sc streamConn
	sc.authData = encodedAuth(c.Username, c.Password)
	sc.postData = body
	sc.url = url_
	resp, err = sc.connect()
	if err != nil {
		goto Return
	}

	if resp.StatusCode != 200 {
		err = errors.New("Twitterstream HTTP Error: " + resp.Status +
			"\n" + url_.Path)
		goto Return
	}

	//close the current connection
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn = &sc
	sc.stream = c.stream
	sc.eventStream = c.eventStream
	sc.friendStream = c.friendStream
	go sc.readStream(resp)

Return:
	return
}

// Follow a list of user ids
func (c *Client) Follow(ids []int64, stream chan *Tweet) error {
	c.stream = stream
	var body bytes.Buffer
	body.WriteString("follow=")
	for i, id := range ids {
		body.WriteString(strconv.FormatInt(id, 10))
		if i != len(ids)-1 {
			body.WriteString(",")
		}
	}
	return c.connect(followUrl, body.String())
}

// Track a list of topics
func (c *Client) Track(topics []string, stream chan *Tweet) error {
	c.stream = stream
	var body bytes.Buffer
	body.WriteString("track=")
	for i, topic := range topics {
		body.WriteString(topic)
		if i != len(topics)-1 {
			body.WriteString(",")
		}
	}
	return c.connect(trackUrl, body.String())
}

// Filter a list of user ids
func (c *Client) Sample(stream chan *Tweet) error {
	c.stream = stream
	return c.connect(sampleUrl, "")
}

// Track User tweets and events
func (c *Client) User(stream chan *Tweet, eventStream chan *Event, friendStream chan *FriendList) error {
	c.stream = stream
	c.eventStream = eventStream
	c.friendStream = friendStream
	return c.connect(userUrl, "")
}

// Close the client
func (c *Client) Close() {
	//has it already been closed?
	if c.conn == nil || c.conn.stale {
		return
	}
	c.conn.Close()
}
