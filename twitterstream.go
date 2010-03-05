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
    "strings"
)


type User struct {
    Lang                         string
    Verified                     bool
    Followers_count              int
    Location                     string
    Screen_name                  string
    Following                    bool
    Friends_count                int
    Profile_background_color     string
    Favourites_count             int
    Description                  string
    Notifications                string
    Profile_text_color           string
    Url                          string
    Time_zone                    string
    Statuses_count               int
    Profile_link_color           string
    Geo_enabled                  bool
    Profile_background_image_url string
    Protected                    bool
    Contributors_enabled         bool
    Profile_sidebar_fill_color   string
    Name                         string
    Profile_background_tile      string
    Created_at                   string
    Profile_image_url            string
    Id                           int64
    Utc_offset                   int
    Profile_sidebar_border_color string
}

type Tweet struct {
    Text                    string
    Truncated               bool
    Geo                     string
    In_reply_to_screen_name string
    Favorited               bool
    Source                  string
    Contributors            string
    In_reply_to_status_id   string
    In_reply_to_user_id     int64
    Id                      int64
    Created_at              string
    User                    User
}

var filterUrl, _ = http.ParseURL("http://stream.twitter.com/1/statuses/filter.json")

type Client struct {
    Username string
    Password string
    conn     *http.ClientConn
    tweets   chan Tweet
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

        if len(line) == 0 {
            continue
        }

        var tweet Tweet
        json.Unmarshal(line, &tweet)

        c.tweets <- tweet
    }
}

// Follow a list of user ids
func (c *Client) Follow(ids []int64) (chan Tweet, os.Error) {
    conn, err := net.Dial("tcp", "", filterUrl.Host+":80")
    if err != nil {
        return nil, err
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
            "Content-Type":  "application/x-www-form-urlencoded",
            "Authorization": "Basic " + encodedAuth(c.Username, c.Password),
        }
    }
    err = c.conn.Write(&req)
    if err != nil {
        return nil, err
    }

    resp, err := c.conn.Read()
    if err != nil {
        return nil, err
    }

    if c.tweets == nil {
        c.tweets = make(chan Tweet)
    }

    go c.readStream(resp)

    return c.tweets, nil
}
