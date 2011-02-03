package twitterstream

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "http"
    "io"
    "json"
    "os"
    "net"
    "strconv"
    "time"
)


var followUrl, _ = http.ParseURL("https://stream.twitter.com/1/statuses/filter.json")
var trackUrl, _ = http.ParseURL("http://stream.twitter.com/1/statuses/filter.json")
var sampleUrl, _ = http.ParseURL("http://stream.twitter.com/1/statuses/sample.json")
var userUrl, _ = http.ParseURL("http://userstream.twitter.com/2/user.json")
var siteStreamUrl, _ = http.ParseURL("https://betastream.twitter.com/2b/site.json")

var retryTimeout int64 = 5e9

type streamConn struct {
    clientConn   *http.ClientConn
    url          *http.URL
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

func (conn *streamConn) connect() (*http.Response, os.Error) {
    if conn.stale {
        return nil, os.NewError("Stale connection")
    }
    var tcpConn net.Conn
    var err os.Error
    if proxy := os.Getenv("HTTP_PROXY"); len(proxy) > 0 {
        proxy_url, _ := http.ParseURL(proxy)
        tcpConn, err = net.Dial("tcp", "", proxy_url.Host)
    } else {
        tcpConn, err = net.Dial("tcp", "", conn.url.Host+":80")
    }
    if err != nil {
        return nil, err
    }
    conn.clientConn = http.NewClientConn(tcpConn, nil)

    var req http.Request
    req.URL = conn.url
    req.Method = "GET"
    req.Header = map[string]string{}
    req.Header["Authorization"] = "Basic " + conn.authData

    if conn.postData != "" {
        req.Method = "POST"
        req.Body = nopCloser{bytes.NewBufferString(conn.postData)}
        req.ContentLength = int64(len(conn.postData))
        req.Header["Content-Type"] = "application/x-www-form-urlencoded"
    }

    err = conn.clientConn.Write(&req)
    if err != nil {
        return nil, err
    }

    resp, err := conn.clientConn.Read()
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
            tcpConn, _ := conn.clientConn.Close()
            if tcpConn != nil {
                tcpConn.Close()
            }
            break
        }

        line, err := reader.ReadBytes('\n')
        if err != nil {
            if conn.stale {
                continue
            }
            resp, err := conn.connect()
            if err != nil {
                println(err.String())
                time.Sleep(retryTimeout)
                continue
            }

            if resp.StatusCode != 200 {
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

func (nopCloser) Close() os.Error { return nil }


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

func (c *Client) connect(url *http.URL, body string) (err os.Error) {
    if c.Username == "" || c.Password == "" {
        return os.NewError("The username or password is invalid")
    }

    var resp *http.Response
    //initialize the new stream
    var sc streamConn

    sc.authData = encodedAuth(c.Username, c.Password)
    sc.postData = body
    sc.url = url
    resp, err = sc.connect()
    if err != nil {
        goto Return
    }

    if resp.StatusCode != 200 {
        err = os.NewError("Twitterstream HTTP Error: " + resp.Status +
            "\n" + url.Path)
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
func (c *Client) Follow(ids []int64, stream chan *Tweet) os.Error {
    c.stream = stream
    var body bytes.Buffer
    body.WriteString("follow=")
    for i, id := range ids {
        body.WriteString(strconv.Itoa64(id))
        if i != len(ids)-1 {
            body.WriteString(",")
        }
    }
    return c.connect(followUrl, body.String())
}

// Track a list of topics
func (c *Client) Track(topics []string, stream chan *Tweet) os.Error {
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
func (c *Client) Sample(stream chan *Tweet) os.Error {
    c.stream = stream
    return c.connect(sampleUrl, "")
}

// Track User tweets and events
func (c *Client) User(stream chan *Tweet, eventStream chan *Event, friendStream chan *FriendList) os.Error {
    c.stream = stream
    c.eventStream = eventStream
    c.friendStream = friendStream
    return c.connect(userUrl, "")
}

// Close the client
func (c *Client) Close() {
    //has it already been closed?
    if c.conn.stale {
        return
    }
    c.conn.Close()
}
