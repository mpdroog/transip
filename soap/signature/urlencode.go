package signature

import (
	"bytes"
	"net/url"
)

type KV struct {
	Key   string
	Value string
}

// Encode encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") NOT sorted.
// https://golang.org/src/net/url/url.go?s=24497:24528#L850
func urlencode(v []KV) []byte {
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
		buf.WriteString(url.QueryEscape(v.Value))
	}
	return buf.Bytes()
}
