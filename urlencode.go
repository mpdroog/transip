package transip

import (
	"bytes"
	"net/url"
	"strings"
)

type kV struct {
	Key   string
	Value string
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") NOT sorted.
// https://golang.org/src/net/url/url.go?s=24497:24528#L850
func urlencode(v []kV) []byte {
	if v == nil {
		return []byte{}
	}
	var buf bytes.Buffer
	for _, v := range v {
		prefix := v.Key + "="

		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(prefix)
		buf.WriteString(strings.Replace(url.QueryEscape(v.Value), "+", "%20", -1))
	}
	return buf.Bytes()
}
