// Package signature creates a cookie signature to validate
// a request for the TransIP API
package transip

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"net/url"
)

// Creates a digest of the given data, with an asn1 header.
// $digest = self::_sha512Asn1(self::_encodeParameters($parameters));
func sha512ASN1(data []byte) []byte {
	asn1 := []byte{
		0x30, 0x51,
		0x30, 0x0d,
		0x06, 0x09,
		0x60, 0x86, 0x48, 0x01, 0x65,
		0x03, 0x04,
		0x02, 0x03,
		0x05, 0x00,
		0x04, 0x40,
	}
	h := sha512.New()
	h.Write(data)
	return append(asn1, h.Sum(nil)...)
}

func Sign(privKey *rsa.PrivateKey, params []KV) (string, error) {
	asn1 := sha512ASN1(urlencode(params))

	sig, e := rsa.SignPKCS1v15(nil, privKey, crypto.Hash(0), asn1)
	if e != nil {
		return "", e
	}

	return url.QueryEscape(base64.StdEncoding.EncodeToString(sig)), nil
}
