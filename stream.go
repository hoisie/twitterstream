package httpstream

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
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

type StreamHandler func([]byte)

func (conn *streamConn) Close() {
	// Just mark the connection as stale, and let the connect() handler close after a read
	conn.stale = true
}

func (conn *streamConn) connect() (*http.Response, error) {
	if conn.stale {
		return nil, errors.New("Stale connection")
	}
	/*
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
	*/
	//cf := &tls.Config{Rand: rand.Reader, Time: time.Now}

	//ssl := tls.Client(tcpConn, cf)

	//transport := http.Transport{TLSClientConfig:cf}
	//conn.client = ssl
	//&http.Client{Transport:&transport}
	//conn.client = &http.Client{Transport:&transport}
	conn.client = &http.Client{}
	fmt.Printf("%v \n", *conn.client)

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
	log.Print("made it to about to conn.Client.Do")
	resp, err := conn.client.Do(&req)
	println("status = ", resp.Status)
	for n, v := range resp.Header {
		println(n, v[0])
	}

	if err != nil {
		log.Print(err)
		return nil, err
	}

	return resp, nil
}

func (conn *streamConn) readStream(resp *http.Response, handler StreamHandler, uniqueId string) {

	var reader *bufio.Reader
	reader = bufio.NewReader(resp.Body)

	/*
	   // lets poll back to see if this conn is closed?
	   timer := time.NewTicker(time.Second * 5)


	   go func() {
	       for _ = range timer.C {
	           println("status = ", resp.Status)
	           if conn.stale {
	               println("Conn is Stale? ")
	               close(conn.stream)
	               conn.Close() 
	           }
	       }
	   }()
	*/

	for {
		//we've been closed
		if conn.stale {
			conn.Close()
			println("Closing connection?")
			//conn.Transport.CloseIdleConnections()
			break
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			println("trying to reconnect? in error")
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

		//println("In twstream", string(line))

		if len(line) == 0 {
			continue
		}
		// keep some metrics, to support re-balancing
		incrCounter(uniqueId)
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

func NewClient(username, password string, handler StreamHandler) *Client {
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
	log.Print("about to call sc.connect")
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
	log.Print(followUrl, " body = ", body.String())
	return c.connect(followUrl, body.String())
}

// Track a list of topics
func (c *Client) Track(topics []string) error {
	var body bytes.Buffer
	body.WriteString("track=")
	for i, topic := range topics {
		body.WriteString(topic)
		if i != len(topics)-1 {
			body.WriteString(",")
		}
	}
	log.Print("TURL ", trackUrl, " body = ", body.String())
	return c.connect(trackUrl, body.String())
}

// get sample stream
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
