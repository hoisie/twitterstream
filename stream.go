/*
A Go http streaming client. http-streaming is most-associated with the twitter stream api.  
This client works with twitter, but has also been tested against the data-sift stream 
as well as a stream-server I use internally (mongrel2)

httpstream was forked from https://github.com/hoisie/twitterstream

*/
package httpstream

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var filterUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/filter.json")
var sampleUrl, _ = url.Parse("https://stream.twitter.com/1/statuses/sample.json")
var userUrl, _ = url.Parse("https://userstream.twitter.com/2/user.json")
var siteStreamUrl, _ = url.Parse("https://sitestream.twitter.com/2b/site.json")

var retryTimeout time.Duration = time.Second * 10

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile | log.Lmicroseconds)
}

type streamConn struct {
	client   *http.Client
	url      *url.URL
	authData string
	postData string
	stale    bool
	// wait time before trying to reconnect, this will be 
	// exponentially moved up until reaching maxWait, when
	// it will exit
	wait    int
	maxWait int
}

func NewStreamConn(max int) streamConn {
	return streamConn{wait: 1, maxWait: max}
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
	if conn.authData != "" {
		req.Header.Set("Authorization", "Basic "+conn.authData)
	}

	if conn.postData != "" {
		req.Method = "POST"
		req.Body = nopCloser{bytes.NewBufferString(conn.postData)}
		req.ContentLength = int64(len(conn.postData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := conn.client.Do(&req)

	if err != nil {
		Log(ERROR, "Could not Connect to Stream: ", err)
		return nil, err
	} else {
		Debugf("connected to %s \n\thttp status = %v", conn.url, resp.Status)
		Debug(resp.Header)
		for n, v := range resp.Header {
			Debug(n, v[0])
		}
	}

	return resp, nil
}

func (conn *streamConn) readStream(resp *http.Response, handler func([]byte), uniqueId string, done chan bool) {

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

			if conn.stale {
				Debug("conn stale, continue")
				continue
			}
			time.Sleep(time.Second * time.Duration(conn.wait))
			//try reconnecting, but exponentially back off until MaxWait is reached
			// then exit?  
			resp, err := conn.connect()
			if err != nil {
				Log(ERROR, " Could not reconnect to source? sleeping and will retry ", conn.wait, err.Error())
				if conn.wait < conn.maxWait {
					conn.wait = conn.wait * 2
				} else {
					Log(ERROR, "exiting, max wait reached")
					done <- true
					return
				}
				continue
			}
			if resp.StatusCode != 200 {
				if conn.wait < conn.maxWait {
					conn.wait = conn.wait * 2
				}
				continue
			}

			reader = bufio.NewReader(resp.Body)
			continue
		} else if conn.wait != 1 {
			conn.wait = 1
		}
		line = bytes.TrimSpace(line)

		if len(line) == 0 {
			continue
			Debug("zer len line?  ended?  ")
		}
		// should we look for twitter stall_warnings and then continue?
		// https://dev.twitter.com/docs/streaming-api/methods

		// TODO:  look for status, stall warnings, etc
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

func getNopCloser(buf *bytes.Buffer) nopCloser {
	return nopCloser{buf}
}

// Client for connecting
type Client struct {
	Username string
	Password string
	// unique id for this connection
	Uniqueid string
	conn     *streamConn
	MaxWait  int
	Handler  func([]byte)
}

func NewBasicAuthClient(username, password string, handler func([]byte)) *Client {
	return &Client{
		Username: username,
		Password: password,
		Handler:  handler,
		MaxWait:  300,
	}
}

func (c *Client) SetMaxWait(max int) {
	c.MaxWait = max
	if c.conn != nil {
		c.conn.maxWait = c.MaxWait
	}
}

func (c *Client) Connect(url_ *url.URL, body string, done chan bool) (err error) {

	var resp *http.Response
	sc := NewStreamConn(c.MaxWait)

	// if http basic auth
	if c.Username != "" && c.Password != "" {
		sc.authData = encodedAuth(c.Username, c.Password)
	}

	sc.postData = body
	sc.url = url_
	resp, err = sc.connect()
	if err != nil {
		Log(ERROR, " errror ", err)
		goto Return
	}

	if resp.StatusCode != 200 {
		Debug("not http 200")
		err = errors.New("stream HTTP Error: " + resp.Status + "\n" + url_.Path)
		goto Return
	}

	//close the current connection
	if c.conn != nil {
		c.conn.Close()
	}

	c.conn = &sc

	go sc.readStream(resp, c.Handler, c.Uniqueid, done)

	return
Return:
	Log(ERROR, "exiting ")
	done <- true
	return
}

// Filter, look for users, topics, and do backfills.   See doc: https://dev.twitter.com/docs/streaming-api/methods
// @userids list of twitter userids to follow (up to 5000). 
// @topics list of words, up to 500
// @done channel to end on ::
//		
//		cl.Filter([]int64{1,2,3,4},nil, done )
//
//		cl.Filter([]int64{1,2,3,4},[]string{"golang"}, done )
//
func (c *Client) Filter(userids []int64, topics []string, watchStalls bool, done chan bool) error {
	var body bytes.Buffer

	body.WriteString("stall_warnings=true&")

	if userids != nil && len(userids) > 0 {
		body.WriteString("follow=")
		for i, id := range userids {
			body.WriteString(strconv.FormatInt(id, 10))
			if i != len(userids)-1 {
				body.WriteString(",")
			}
		}
		body.WriteString("&")
	}
	if topics != nil && len(topics) > 0 {
		body.WriteString("track=")
		for i, topic := range topics {
			body.WriteString(topic)
			if i != len(topics)-1 {
				body.WriteString(",")
			}
		}
		body.WriteString("&")
	}
	Debug("TWFILTER ", filterUrl, body.String())
	if watchStalls {
		c.Handler = stallWatcher(c.Handler)
	}

	return c.Connect(filterUrl, body.String(), done)
}

func stallWatcher(handler func([]byte)) func([]byte) {
	/*
		{ "warning":{
			"code":"FALLING_BEHIND",
			"message":"Your connection is falling behind and messages are being queued for delivery to you. Your queue is now over 60% full. You will be disconnected when the queue is full.",
			"percent_full": 60
		  }
		}
	*/
	lookFor := []byte(`"code":"FALLING_BEHIND"`)
	pctFull := []byte(`"percent_full"`)
	return func(line []byte) {
		if bytes.Index(line, lookFor) > 0 {
			idx := bytes.Index(line, pctFull)
			Log(ERROR, "FALLING BEHIND!!!!  ", string(line[idx+1:idx+5]))
		} else {
			handler(line)
		}

	}
}

// twitter sample stream
func (c *Client) Sample(done chan bool) error {
	return c.Connect(sampleUrl, "", done)
}

// Track User tweets and events, uses passed username/pwd
func (c *Client) User(done chan bool) error {
	return c.Connect(userUrl, "", done)
}

// Close the client
func (c *Client) Close() {
	//has it already been closed?
	if c.conn == nil || c.conn.stale {
		return
	}
	c.conn.Close()
}
