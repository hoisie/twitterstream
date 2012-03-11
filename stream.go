package httpstream

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var followUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/filter.json")
var trackUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/filter.json")
var sampleUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/sample.json")
var userUrl, _ = url.Parse("https://userstream.twitter.com/2/user.json")
var siteStreamUrl, _ = url.Parse("https://sitestream.twitter.com/2b/site.json")

var retryTimeout time.Duration = 5e9

type streamConn struct {
	client   *http.Client
	url      *url.URL
	authData string
	postData string
	stale    bool
}

//type StreamHandler func([]byte)

func (conn *streamConn) Close() {
	// Just mark the connection as stale, and let the connect() handler close after a read
	conn.stale = true
}

func (conn *streamConn) connect() (*http.Response, error) {
	if conn.stale {
		return nil, errors.New("Stale connection")
	}

	conn.client = &http.Client{}

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
	
	resp, err := conn.client.Do(&req)
	Debugf("connected to %s \n\thttp status = %d", conn.url, resp.Status)
	for n, v := range resp.Header {
		Debug(n, v[0])
	}

	if err != nil {
		Log(ERROR, "Could not Connect to Stream: ", err)
		return nil, err
	}

	return resp, nil
}

func (conn *streamConn) readStream(resp *http.Response, handler func([]byte), uniqueId string) {

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)

	for {
		//we've been closed
		if conn.stale {
			conn.Close()
			Debug("Connection closed, shutting down ")
			//conn.Transport.CloseIdleConnections()
			break
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			Debug("trying to reconnect? in error")
			if conn.stale {
				continue
			}

			//try reconnecting
			resp, err := conn.connect()
			if err != nil {
				Log(ERROR, " Could not reconnect to source? sleeping and will retry ", err.Error())
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
		// should we look for twitter stall_warnings and then continue?
		// https://dev.twitter.com/docs/streaming-api/methods

		// keep some metrics, to support re-balancing
		IncrCounter(uniqueId)

		handler(line)
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

func (nopCloser) Close() error {
	return nil
}

type Client struct {
	Username string
	Password string
	// unique id for this connection
	Uniqueid    string
	conn    *streamConn
	Handler func([]byte)
}

func NewClient(username, password string, handler func([]byte)) *Client {
	return &Client{
		Username: username,
		Password: password,
		Handler:  handler,
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
		err = errors.New("Twitterstream HTTP Error: " + resp.Status + "\n" + url_.Path)
		goto Return
	}

	//close the current connection
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn = &sc
	go sc.readStream(resp, c.Handler, c.Uniqueid)

Return:
	return
}

// Follow a list of user ids
func (c *Client) Follow(ids []int64) error {
	var body bytes.Buffer
	body.WriteString("follow=")
	for i, id := range ids {
		body.WriteString(strconv.FormatInt(id, 10))
		if i != len(ids)-1 {
			body.WriteString(",")
		}
	}
	Debug("TWFOLLOW ", followUrl, body.String())
	return c.connect(followUrl, body.String())
}

// Twitter Track a list of topics
func (c *Client) Track(topics []string) error {
	var body bytes.Buffer
	body.WriteString("track=")
	for i, topic := range topics {
		body.WriteString(topic)
		if i != len(topics)-1 {
			body.WriteString(",")
		}
	}
	Debug("TWTRACK ", trackUrl, " body = ", body.String())
	return c.connect(trackUrl, body.String())
}

// twitter sample stream
func (c *Client) Sample() error {
	return c.connect(sampleUrl, "")
}

// Track User tweets and events, uses passed username/pwd
func (c *Client) User() error {
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
