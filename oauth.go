package twitterstream

import (
    "bytes"
    "container/vector"
    "crypto/hmac"
    "encoding/base64"
    "fmt"
    "http"
    "httplib"
    "os"
    "rand"
    "sort"
    "strconv"
    "time"
)

var requestTokenUrl, _ = http.ParseURL("https://api.twitter.com/oauth/request_token")
var accessTokenUrl, _ = http.ParseURL("https://api.twitter.com/oauth/access_token")
var authorizeUrl, _ = http.ParseURL("https://api.twitter.com/oauth/authorize")

type OAuthClient struct {
    ConsumerKey      string
    ConsumerSecret   string
    OAuthToken       string
    OAuthTokenSecret string
}

type RequestToken struct {
    OAuthTokenSecret       string
    OAuthToken             string
    OAuthCallbackConfirmed bool
}

type AccessToken struct {
    OAuthToken       string
    OAuthTokenSecret string
    UserId           string
    ScreenName       string
}

func getNonce(n int) string {
    var alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    buf := make([]byte, n)
    for i := 0; i < n; i++ {
        buf[i] = alphabet[rand.Intn(len(alphabet))]
    }
    return string(buf)
}

func signatureBase(httpMethod string, base_uri string, params map[string]string) string {
    var buf bytes.Buffer

    buf.WriteString(httpMethod)
    buf.WriteString("&")
    buf.WriteString(URLEscape(base_uri))
    buf.WriteString("&")

    var keys vector.StringVector
    for k, _ := range params {
        keys.Push(k)
    }

    sort.SortStrings(keys)
    for i, k := range keys {
        v := params[k]
        buf.WriteString(URLEscape(k))
        buf.WriteString("%3D")
        buf.WriteString(URLEscape(v))
        //don't include the dangling %26
        if i < len(params)-1 {
            buf.WriteString("%26")
        }
        i++
    }
    return buf.String()
}

func signRequest(base string, consumerSecret string, tokenSecret string) string {
    signingKey := URLEscape(consumerSecret) + "&"
    if tokenSecret != "" {
        signingKey += URLEscape(tokenSecret)
    }

    hash := hmac.NewSHA1([]byte(signingKey))
    hash.Write([]byte(base))
    sum := hash.Sum()
    bb := new(bytes.Buffer)
    encoder := base64.NewEncoder(base64.StdEncoding, bb)
    encoder.Write(sum)
    encoder.Close()
    return bb.String()
}

func (o *OAuthClient) GetRequestToken(callback string) *RequestToken {
    nonce := getNonce(40)
    params := map[string]string{
        "oauth_nonce":            nonce,
        "oauth_callback":         URLEscape(callback),
        "oauth_signature_method": "HMAC-SHA1",
        "oauth_timestamp":        strconv.Itoa64(time.Seconds()),
        "oauth_consumer_key":     o.ConsumerKey,
        "oauth_version":          "1.0",
    }

    base := signatureBase("POST", requestTokenUrl.Raw, params)
    signature := signRequest(base, o.ConsumerSecret, "")
    params["oauth_signature"] = URLEscape(signature)

    authBuf := bytes.NewBufferString("OAuth ")
    i := 0
    for k, v := range params {
        authBuf.WriteString(fmt.Sprintf("%s=%q", k, v))
        if i < len(params)-1 {
            authBuf.WriteString(", ")
        }
        i++
    }
    request := httplib.Post(requestTokenUrl.Raw)
    request.Header("Authorization", authBuf.String())
    request.Body("")
    resp, err := request.AsString()
    tokens, err := http.ParseQuery(resp)
    if err != nil {
        println(err.String())
    }

    confirmed, _ := strconv.Atob(tokens["oauth_callback_confirmed"][0])
    rt := RequestToken{
        OAuthTokenSecret:       tokens["oauth_token_secret"][0],
        OAuthToken:             tokens["oauth_token"][0],
        OAuthCallbackConfirmed: confirmed,
    }
    return &rt
}

func (rt *RequestToken) AuthorizeUrl() string {
    return fmt.Sprintf("%s?oauth_token=%s", authorizeUrl.Raw, rt.OAuthToken)
}

func (o *OAuthClient) GetAccessToken(requestToken *RequestToken, OAuthVerifier string) *AccessToken {
    nonce := getNonce(40)
    params := map[string]string{
        "oauth_nonce":            nonce,
        "oauth_token":            requestToken.OAuthToken,
        "oauth_verifier":         OAuthVerifier,
        "oauth_signature_method": "HMAC-SHA1",
        "oauth_timestamp":        strconv.Itoa64(time.Seconds()),
        "oauth_consumer_key":     o.ConsumerKey,
        "oauth_version":          "1.0",
    }

    base := signatureBase("POST", requestTokenUrl.Raw, params)
    signature := signRequest(base, o.ConsumerSecret, requestToken.OAuthTokenSecret)
    params["oauth_signature"] = URLEscape(signature)

    authBuf := bytes.NewBufferString("OAuth ")
    i := 0
    for k, v := range params {
        authBuf.WriteString(fmt.Sprintf("%s=%q", k, v))
        if i < len(params)-1 {
            authBuf.WriteString(", ")
        }
        i++
    }
    request := httplib.Post(accessTokenUrl.Raw)
    request.Header("Authorization", authBuf.String())
    request.Body("")
    resp, err := request.AsString()
    tokens, err := http.ParseQuery(resp)
    if err != nil {
        println(err.String())
    }

    at := AccessToken{
        OAuthTokenSecret: tokens["oauth_token_secret"][0],
        OAuthToken:       tokens["oauth_token"][0],
        UserId:           tokens["user_id"][0],
        ScreenName:       tokens["screen_name"][0],
    }
    return &at

}

func (o *OAuthClient) OAuthConnect(url string) (*http.Response, os.Error) {
    nonce := getNonce(40)
    params := map[string]string{
        "oauth_nonce":            nonce,
        "oauth_token":            o.OAuthToken,
        "oauth_signature_method": "HMAC-SHA1",
        "oauth_timestamp":        strconv.Itoa64(time.Seconds()),
        "oauth_consumer_key":     o.ConsumerKey,
        "oauth_version":          "1.0",
    }

    base := signatureBase("GET", url, params)
    signature := signRequest(base, o.ConsumerSecret, o.OAuthTokenSecret)

    params["oauth_signature"] = URLEscape(signature)

    authBuf := bytes.NewBufferString("OAuth ")
    i := 0
    for k, v := range params {
        authBuf.WriteString(fmt.Sprintf("%s=%q", k, v))
        if i < len(params)-1 {
            authBuf.WriteString(", ")
        }
        i++
    }

    request := httplib.Get(url)
    request.Header("Authorization", authBuf.String())
    request.Body("")
    return request.AsResponse()
}
