package twitterstream

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "fmt"
    "http"
    "io"
    "json"
    "os"
    "net"
    "strconv"
    "strings"
    "sync"
)

var filterApi, _ = http.ParseURL("http://stream.twitter.com/1/statuses/filter.json")

type streamConn struct {
    clientConn *http.ClientConn
    stream     chan Tweet
    authData   string
    postData   string
    stale      bool
}

func (conn *streamConn) Close() {
    println("closing the conn!")
    conn.stale = true
    tcpConn, _ := conn.clientConn.Close()
    tcpConn.Close()
}

func (conn *streamConn) connect() (*http.Response, os.Error) {
    tcpConn, err := net.Dial("tcp", "", filterApi.Host+":80")
    if err != nil {
        return nil, err
    }
    conn.clientConn = http.NewClientConn(tcpConn, nil)

    var req http.Request
    req.URL = filterApi
    req.Method = "POST"
    req.Body = nopCloser{bytes.NewBufferString(conn.postData)}
    req.ContentLength = int64(len(conn.postData))
    req.Header = map[string]string{
        "Content-Type": "application/x-www-form-urlencoded",
    }

    if conn.authData != "" {
        req.Header["Authorization"] = "Basic " + conn.authData
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
    reader := bufio.NewReader(resp.Body)
    fmt.Println("Readstream started!")
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            println(err.String())
            //we've been closed
            if conn.stale {
                return
            }

            //otherwise, reconnect
            resp, err := conn.connect()
            if err != nil {
                println(err.String())
            }

            if resp.StatusCode != 200 {
                println("HTTP Error" + resp.Status)
            }

            continue
        }
        line = strings.TrimSpace(line)

        if len(line) == 0 {
            continue
        }

        var tweet Tweet
        json.Unmarshal(line, &tweet)

        conn.stream <- tweet
    }
}


type FilterStream struct {
    Username string
    Password string
    Stream   chan Tweet
    conn     *streamConn
    connLock *sync.Mutex
}

func NewFilterStream(username, password string) *FilterStream {
    return &FilterStream{username, password, make(chan Tweet), nil, new(sync.Mutex)}
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

// Follow a list of user ids
func (c *FilterStream) Follow(ids []int64) os.Error {
    c.connLock.Lock()
    println("in follow!")
    var body bytes.Buffer
    body.WriteString("follow=")
    for i, id := range ids {
        body.WriteString(strconv.Itoa64(id))
        if i != len(ids)-1 {
            body.WriteString(",")
        }
    }

    var sc streamConn
    sc.authData = encodedAuth(c.Username, c.Password)
    sc.postData = body.String()
    resp, err := sc.connect()
    if err != nil {
        return err
    }

    if resp.StatusCode != 200 {
        return os.NewError("HTTP Error" + resp.Status)
    }

    if c.conn != nil {
        c.conn.Close()
    }

    c.conn = &sc
    sc.stream = c.Stream
    go sc.readStream(resp)
    c.connLock.Unlock()

    return nil
}
