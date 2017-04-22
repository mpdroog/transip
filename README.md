TransIP API
==================
Small library in Golang that implements:
* DomainService/getDomainNames
* DomainService/getInfo
* DomainService/setDnsEntries

If you need other methods on the TransIP API be free to fork my
code and I'll merge it with love. :)

Why use this library instead of writing own:
* 'Correctly' implementing __nonce because there's no uniqid in Golang;
* Correctly implementing signature, RSA Private key signing of cookie signature (unittests included);
* Tokenized XML-parser to keep data structures free of clutter (stripping SOAP envelopes);

Note
=======
As far as I know Golang/crypto doesn't support loading PKCS8-certificates as RSA.PrivateKey. So if you
want to start using this code please convert your privatekey from PKCS8 to PKCS1 with OpenSSL:
```
openssl rsa -in privkeyfromtransip.pem -out privkey.pem
```

Example code
=======
Get a list of domains and return all of their details
```go
package main

import (
	"github.com/mpdroog/transip"
	"github.com/mpdroog/transip/creds"
	"fmt"
)

func main() {
	c := transip.DomainService{
		Creds: creds.Client{
			Login: "mdroog",
			PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
REMOVED :)
-----END RSA PRIVATE KEY-----`,
			ReadWrite: false,
		},
	}

	domains, e := c.DomainNames()
	if e != nil {
		panic(e)
	}

	fmt.Printf("Managed domains:\n")
	for _, domain := range domains.Item {
		fmt.Printf("★ %s\n", domain)

		info, e := c.Domain(domain)
		if e != nil {
			panic(e)
		}
		fmt.Printf("\t%+v\n\n", info)
	}
}
```

Change a domain's DNS entries
```go
package main

import (
	"github.com/mpdroog/transip"
	"github.com/mpdroog/transip/creds"
	"fmt"
)

func main() {
	c := transip.DomainService{
		Creds: creds.Client{
			Login: "mdroog",
			PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
REMOVED :)
-----END RSA PRIVATE KEY-----`,
			ReadWrite: false,
		},
	}

	// 360 = 6min (TTL in seconds)
	if e := c.SetDNSEntries(
		"yourdomain.com",
		[]soap.DomainDNSentry{soap.DomainDNSentry{Name: "@", Expire: 360, Type: "A", Content: "127.0.0.1"}},
	); e != nil {
		panic(e)
	}
}
```