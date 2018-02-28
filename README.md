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

func printDomainNames(username, privKeyPath string) error {
    creds := creds.Client{
        Login:     username,
        ReadWrite: false,
    }
    if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
        return fmt.Errorf("could not load private key from path %s: %s",
            privKeyPath, err)
    }
    	
    domainService := transip.DomainService{creds}
	domains, err := domainService.DomainNames()
	if err != nil {
		return err
	}

	fmt.Printf("Managed domains:\n")
	for _, domain := range domains.Item {
		fmt.Printf("★ %s\n", domain)

		info, err := domainService.Domain(domain)
		if err != nil {
			return err
		}
		fmt.Printf("\t%+v\n\n", info)
	}
	
	return nil
}
```

Overwrite a domain's DNS entries
```go
package main

import (
	"github.com/mpdroog/transip"
	"github.com/mpdroog/transip/creds"
    "github.com/mpdroog/transip/soap"
	"fmt"
)

func overWriteDnsEntries(username, privKeyPath, domain string) error {
    creds := creds.Client{
        Login:     username,
        ReadWrite: false,
    }
    if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
        return fmt.Errorf("could not load private key from path %s: %s",
            privKeyPath, err)
    }
        
    // 360 = 6min (TTL in seconds)
    recordSet := []soap.DomainDNSentry{
    	{Name: "@", Expire: 360, Type: "A", Content: "127.0.0.1"},
    }
    
    domainService := transip.DomainService{creds}
    err := domainService.SetDNSEntries(domain, recordSet)
    return err
}
```