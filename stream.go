/*
A Go http streaming client. http-streaming is most-associated with the twitter stream api.
This client works with twitter, but has also been tested against the data-sift stream and
flowdock stream api's

httpstream was forked from https://github.com/hoisie/twitterstream

*/
package httpstream

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mrjones/oauth"
)

var (
	filterUrl, _                   = url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
	sampleUrl, _                   = url.Parse("https://stream.twitter.com/1.1/statuses/sample.json")
	userUrl, _                     = url.Parse("https://userstream.twitter.com/2/user.json")
	siteStreamUrl, _               = url.Parse("https://sitestream.twitter.com/2b/site.json")
	retryTimeout     time.Duration = time.Second * 10
	OauthCon         *oauth.Consumer
)

type streamConn struct {
	client   *http.Client
	resp     *http.Response
	url      *url.URL
	at       *oauth.AccessToken
	authData string
	postData string
	stale    bool
	closed   bool
	mu       sync.Mutex
	// wait time before trying to reconnect, this will be
	// exponentially moved up until reaching maxWait, when
	// it will exit
	wait     int
	maxWait  int
	insecure bool
	connect  func() (*http.Response, error)
}

func NewStreamConn(max int, insecure bool) streamConn {
	return streamConn{wait: 1, maxWait: max, insecure: insecure}
}

//type StreamHandler func([]byte)

func (conn *streamConn) Close() {
	// Just mark the connection as stale, and let the connect() handler close after a read
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.stale = true
	conn.closed = true
	if conn.resp != nil {
		conn.resp.Body.Close()
	}
}

func basicauthConnect(conn *streamConn) (*http.Response, error) {
	if conn.isStale() {
		return nil, errors.New("Stale connection")
	}

	conn.client = &http.Client{}
	httpTransport := &http.Transport{}
	if conn.insecure {
		httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	httpTransport.TLSHandshakeTimeout = 30 * time.Second
	httpTransport.ResponseHeaderTimeout = 30 * time.Second
	conn.client.Transport = httpTransport

	var req http.Request
	req.URL = conn.url
	req.Method = "GET"
	req.Header = http.Header{}
	if conn.authData != "" {
		req.Header.Set("Authorization", conn.authData)
	}

	if conn.postData != "" {
		req.Method = "POST"
		req.Body = nopCloser{bytes.NewBufferString(conn.postData)}
		req.ContentLength = int64(len(conn.postData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	Debug(req.Header)
	Debug(conn.postData)
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

func oauthConnect(conn *streamConn, params map[string]string) (*http.Response, error) {
	if conn.isStale() {
		return nil, errors.New("Stale connection")
	}

	resp, err := OauthCon.Post(conn.url.String(), params, conn.at)

	if err != nil {
		if resp != nil && resp.Body != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			Log(ERROR, err, " ", string(data))
			resp.Body.Close()
		} else {
			Log(ERROR, err)
		}

	} else {
		Debugf("connected to %s \n\thttp status = %v", conn.url, resp.Status)
		Debug(resp.Header)
		for n, v := range resp.Header {
			Debug(n, v[0])
		}
	}

	return resp, nil
}

func formString(params map[string]string) string {
	var body bytes.Buffer
	for k, v := range params {
		body.WriteString(URLEscape(k))
		body.WriteString("=")
		body.WriteString(URLEscape(v))
		body.WriteString("&")
	}
	return body.String()
}

func (conn *streamConn) isStale() bool {
	conn.mu.Lock()
	r := conn.stale
	conn.mu.Unlock()
	return r
}

func (conn *streamConn) readStream(resp *http.Response, handler func([]byte), uniqueId string, done chan bool) {

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)
	conn.resp = resp

	for {
		//we've been closed
		if conn.isStale() {
			conn.Close()
			Debug("Connection closed, shutting down ")
			done <- true
			break
		}

		line, err := reader.ReadBytes('\n')

		if err != nil {

			if conn.isStale() {
				Debug("conn stale, continue")
				continue
			}
			time.Sleep(time.Second * time.Duration(conn.wait))
			//try reconnecting, but exponentially back off until MaxWait is reached then exit?
			resp, err := conn.connect()
			if err != nil || resp == nil {
				Log(ERROR, " Could not reconnect to source? sleeping and will retry ", err)
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
		}
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
	Insecure bool
	at       *oauth.AccessToken
	Handler  func([]byte)
}

func NewClient(handler func([]byte)) *Client {
	return &Client{
		Handler: handler,
		MaxWait: 300,
	}
}

func NewOAuthClient(at *oauth.AccessToken, handler func([]byte)) *Client {
	return &Client{
		at:      at,
		Handler: handler,
		MaxWait: 300,
	}
}

func NewBasicAuthClient(username, password string, handler func([]byte)) *Client {
	return &Client{
		Username: username,
		Password: password,
		Handler:  handler,
		MaxWait:  300,
	}
}

// Create a new basic Auth Channel Handler
func NewChannelClient(username, password string, bc chan []byte) *Client {
	return &Client{
		Username: username,
		Password: password,
		Handler:  func(b []byte) { bc <- b },
		MaxWait:  300,
	}
}

func (c *Client) SetMaxWait(max int) {
	c.MaxWait = max
	if c.conn != nil {
		c.conn.maxWait = c.MaxWait
	}
}

func (c *Client) SetInsecure(insecure bool) {
	c.Insecure = insecure
	if c.conn != nil {
		c.conn.insecure = c.Insecure
	}
}

// Connect to an http stream
// @url = http address
// @params = http params to be added
func (c *Client) Connect(url_ *url.URL, params map[string]string, done chan bool) (err error) {

	var resp *http.Response
	sc := NewStreamConn(c.MaxWait, c.Insecure)

	sc.url = url_
	// if http basic auth
	if c.Username != "" && c.Password != "" {
		sc.postData = formString(params)
		sc.authData = "Basic " + encodedAuth(c.Username, c.Password)
		sc.connect = func() (*http.Response, error) {
			return basicauthConnect(&sc)
		}

	} else {
		sc.at = c.at
		sc.connect = func() (*http.Response, error) {
			return oauthConnect(&sc, params)
		}

	}
	resp, err = sc.connect()
	if err != nil {
		Log(ERROR, " error ", err)
		goto Return
	} else if resp == nil {
		err = errors.New("No response on connection, invalid connect")
		Log(ERROR, err.Error())
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

// Filter, look for users, topics.   See doc: https://dev.twitter.com/docs/streaming-api/methods
// @userids list of twitter userids to follow (up to 5000).
// @topics list of words, up to 500
// @languages:  list of languages to filter for
// @locations:  optional list of locations
// @done channel to end on ::
//
//		cl.Filter([]int64{1,2,3,4},nil, nil, nil, false, done )
//
//		cl.Filter([]int64{1,2,3,4},[]string{"golang"},[]string{"en"}, nil, false, done )
//
func (c *Client) Filter(userids []int64, topics []string, languages []string, locations []string, watchStalls bool, done chan bool) error {

	params := make(map[string]string)
	params["stall_warnings"] = "true"
	if userids != nil && len(userids) > 0 {
		users := make([]string, 0)
		for _, id := range userids {
			users = append(users, strconv.FormatInt(id, 10))
		}
		params["follow"] = strings.Join(users, ",")
	}

	if topics != nil && len(topics) > 0 {
		params["track"] = strings.Join(topics, ",")
	}

	if languages != nil && len(languages) > 0 {
		params["language"] = strings.Join(languages, ",")
	}

	if locations != nil && len(locations) > 0 {
		params["locations"] = strings.Join(locations, ",")
	}

	if watchStalls {
		c.Handler = StallWatcher(c.Handler)
	}

	return c.Connect(filterUrl, params, done)
}

// a handler wrapper to watch for twitter stall warnings
func StallWatcher(handler func([]byte)) func([]byte) {
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
	return c.Connect(sampleUrl, nil, done)
}

// Track User tweets and events, uses passed username/pwd
func (c *Client) User(done chan bool) error {
	return c.Connect(userUrl, nil, done)
}

// Close the client
func (c *Client) Close() {
	//has it already been closed?
	if c.conn == nil || c.conn.isStale() {
		return
	}
	c.conn.Close()
}
