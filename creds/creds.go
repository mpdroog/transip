package creds

type Client struct {
	Login      string // Login to TransIP website
	PrivateKey string // Private key for API
	ReadWrite  bool   // Read+Write mode?
}
