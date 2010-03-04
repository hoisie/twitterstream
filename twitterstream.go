package twitterstream

import (
    "bufio"
    "bytes"
    "encoding/base64"
    "http"
    "io"
    "os"
    "net"
    "strconv"
    "strings"
)

var filterUrl, _ = http.ParseURL("http://stream.twitter.com/1/statuses/filter.json")

type Client struct {
    Username string
    Password string
    conn     *http.ClientConn
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

func (c *Client) readStream(resp *http.Response) {
    reader := bufio.NewReader(resp.Body)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            println(err.String())
            return
        }
        line = strings.TrimSpace(line)
        if len(line) > 0 {
            println(line)
        } else {
            println("nop")
        }
    }
}

// Follow a list of user ids
func (c *Client) Follow(ids []int64) os.Error {
    conn, err := net.Dial("tcp", "", filterUrl.Host+":80")
    if err != nil {
        return err
    }
    c.conn = http.NewClientConn(conn, nil)

    var req http.Request
    req.URL = filterUrl
    req.Method = "POST"

    var body bytes.Buffer
    body.WriteString("follow=")
    for i, id := range ids {
        body.WriteString(strconv.Itoa64(id))
        if i != len(ids)-1 {
            body.WriteString(",")
        }
    }

    req.Body = nopCloser{&body}
    req.ContentLength = int64(len(body.String()))
    if c.Username != "" {
        req.Header = map[string]string{
            "Content-Type": "application/x-www-form-urlencoded",
            "Authorization": "Basic " + encodedAuth(c.Username, c.Password),
        }
    }
    err = c.conn.Write(&req)
    if err != nil {
        return err
    }

    resp, err := c.conn.Read()
    if err != nil {
        return err
    }

    go c.readStream(resp)

    return nil
}
