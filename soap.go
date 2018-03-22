// Package soap implements SOAP-logic for the TransIP API.
package transip

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"fmt"
	"strconv"
)

const API_URL = "api.transip.nl"
const API_NAMESPACE = "http://www.transip.nl/soap"

// PHP's uniqid 'clone'
// https://github.com/php/php-src/blob/master/ext/standard/uniqid.c
func uniqid() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(
		"%08x%05x%.8d",
		time.Now().Unix(),
		time.Now().UnixNano()/1000,
		rand.Intn(100000000),
	)
}

type request struct {
	Service     string         // Service to call on TransIP side
	ExtraParams []kV // Additional params for the signature-code
	Body        string         // XML body to send in envelope
	Method      string         // Method to call on service
}

func lookup(c Client, in request) ([]byte, error) {
	raw := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"
	xmlns:ns1="%s" xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/"
	SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
		<SOAP-ENV:Body>%s</SOAP-ENV:Body>
</SOAP-ENV:Envelope>`, API_NAMESPACE, in.Body)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/soap/?service=%s", API_URL, in.Service), strings.NewReader(raw))
	if err != nil {
		return []byte{}, err
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := uniqid()

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Set("User-Agent", "mpdroog/soapclient")

	req.AddCookie(&http.Cookie{
		Name:  "login",
		Value: c.Login,
	})

	mode := "readonly"
	if c.ReadWrite {
		mode = "readwrite"
	}
	req.AddCookie(&http.Cookie{
		Name:  "mode",
		Value: mode,
	})

	req.AddCookie(&http.Cookie{
		Name:  "timestamp",
		Value: now,
	})
	req.AddCookie(&http.Cookie{
		Name:  "nonce",
		Value: nonce,
	})
	req.AddCookie(&http.Cookie{
		Name:  "clientVersion",
		Value: "0.1",
	})

	kv := in.ExtraParams
	kv = append(kv, []kV{
		{"__method", in.Method},
		{"__service", in.Service},
		{"__hostname", API_URL},
		{"__timestamp", now},
		{"__nonce", nonce}}...,
	)
	sig, e := sign(c.PrivateKey, kv)
	if e != nil {
		return []byte{}, e
	}
	req.AddCookie(&http.Cookie{
		Name:  "signature",
		Value: sig,
	})

	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	if len(rawbody) == 0 {
		return []byte{}, errors.New("empty response")
	}
	return rawbody, nil
}
