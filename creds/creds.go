package creds

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

type Client struct {
	Login      string          // Login to TransIP website
	PrivateKey *rsa.PrivateKey // Private key for API
	ReadWrite  bool            // Read+Write mode?
}

func (c *Client) SetPrivateKeyFromPath(path string) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return c.SetPrivateKeyFromBytes(contents)
}

func (c *Client) SetPrivateKeyFromBytes(keyContents []byte) error {
	block, _ := pem.Decode(keyContents)
	if block == nil {
		return errors.New("could not decode pem file")
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	c.PrivateKey = privKey.(*rsa.PrivateKey)

	return nil
}
