package twitterstream

// URLEscape converts a string into ``URL encoded'' form.
// Despite the name, this encoding applies only to individual
// components of the query portion of the URL.
func URLEscape(s string) string {
    return urlEscape(s)
}

func urlEscape(s string) string {
    spaceCount, hexCount := 0, 0
    for i := 0; i < len(s); i++ {
        c := s[i]
        if shouldEscape(c) {
            if c == ' ' {
                spaceCount++
            } else {
                hexCount++
            }
        }
    }

    if spaceCount == 0 && hexCount == 0 {
        return s
    }

    t := make([]byte, len(s)+2*hexCount)
    j := 0
    for i := 0; i < len(s); i++ {
        switch c := s[i]; {
        case c == ' ':
            t[j] = '+'
            j++
        case shouldEscape(c):
            t[j] = '%'
            t[j+1] = "0123456789ABCDEF"[c>>4]
            t[j+2] = "0123456789ABCDEF"[c&15]
            j += 3
        default:
            t[j] = s[i]
            j++
        }
    }
    return string(t)
}

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 2396.
// When 'all' is true the full range of reserved characters are matched.
func shouldEscape(c byte) bool {
    // RFC 2396 §2.3 Unreserved characters (alphanum)
    if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
        return false
    }
    switch c {
    case '-', '_', '.', '!', '~', '*', '\'', '(', ')': // §2.3 Unreserved characters (mark)
        return false

    case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
        // Different sections of the URL allow a few of
        // the reserved characters to appear unescaped.
        return true
    }

    // Everything else must be escaped.
    return true
}
