TransIP API
==================
Small library in Golang that implements:
* DomainService/getDomainNames
* DomainService/getInfo
* DomainService/batchGetInfo
* DomainService/setDnsEntries
* DomainService/batchCheckAvailability
* DomainService/checkAvailability

If you need other methods on the TransIP API be free to fork my
code and I'll merge it with love. :)

Why use this library instead of writing own:
* 'Correctly' implementing \_\_nonce because there's no uniqid in Golang;
* Correctly implementing signature, RSA Private key signing of cookie signature (unittests included);
* Tokenized XML-parser to keep data structures free of clutter (stripping SOAP envelopes);

Example code
=======
Get info for a domain
```go
package main

import (
	"github.com/mpdroog/transip"
	"fmt"
)

func printDomainInfo(username, privKeyPath string) error {
	creds := transip.Client{
		Login:     username,
		ReadWrite: false,
	}
	if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
		return fmt.Errorf("could not load private key from path %s: %s",
			privKeyPath, err)
	}
    	
	domainService := transip.DomainService{creds}
	domain, err := domainService.Domain("example.com")
	if err != nil {
		return err
	}
	fmt.Printf("\t%+v\n\n", domain)
	return nil
}
```

Get a list of domain names and return all of their details
```go
package main

import (
	"github.com/mpdroog/transip"
	"fmt"
)

func printDomainNames(username, privKeyPath string) error {
	creds := transip.Client{
		Login:     username,
		ReadWrite: false,
	}
	if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
		return fmt.Errorf("could not load private key from path %s: %s",
			privKeyPath, err)
	}
    	
	domainService := transip.DomainService{creds}
	domainNames, err := domainService.DomainNames()
	if err != nil {
		return err
	}
	domains, err := domainService.Domains(domainNames)
	if err != nil {
		return err
	}
	fmt.Printf("Managed domains:\n")
	for _, domain := range domains {
		fmt.Printf("★ %s\n", domain.Name)
		fmt.Printf("\t%+v\n\n", domain)
	}

	return nil
}
```

Overwrite a domain's DNS entries
```go
package main

import (
	"github.com/mpdroog/transip"
	"fmt"
)

func overWriteDnsEntries(username, privKeyPath, domain string) error {
	creds := transip.Client{
		Login:     username,
		ReadWrite: true,
	}
	if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
		return fmt.Errorf("could not load private key from path %s: %s",
	    		privKeyPath, err)
	}

	// 360 = 6min (TTL in seconds)
	recordSet := []transip.DomainDNSentry{
		{Name: "@", Expire: 360, Type: "A", Content: "127.0.0.1"},
	}

	domainService := transip.DomainService{creds}
	err := domainService.SetDNSEntries(domain, recordSet)
	return err
}
```

Check availability of a domain
```go
package main

import (
	"github.com/mpdroog/transip"
	"fmt"
)

func checkAvailability(username, privKeyPath, domain string) (string, error) {
	creds := transip.Client{
        Login:     username,
        ReadWrite: false,
    }
    if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
        return "", fmt.Errorf("could not load private key from path %s: %s",
            privKeyPath, err)
    }

    domainService := transip.DomainService{creds}
	res, err := domainService.CheckAvailability(domain)
	if err != nil {
		return "", err
	}

    return res, nil
}
```

Check availability of a batch of domains
```go
package main

import (
	"github.com/mpdroog/transip"
	"fmt"
)

func checkAvailability(username, privKeyPath string, domains []string) ([]transip.DomainCheckResult, error) {
	creds := transip.Client{
        Login:     username,
        ReadWrite: false,
    }
    if err := creds.SetPrivateKeyFromPath(privKeyPath); err != nil {
        return nil, fmt.Errorf("could not load private key from path %s: %s",
            privKeyPath, err)
    }

    domainService := transip.DomainService{creds}
	res, err := domainService.BatchCheckAvailability(domains)
	if err != nil {
		return nil, err
	}

    return res, nil
}
```