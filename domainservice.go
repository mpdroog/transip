package transip

import (
	"fmt"
	"github.com/mpdroog/transip/creds"
	"github.com/mpdroog/transip/soap"
	"github.com/mpdroog/transip/soap/signature"
)

const domainService = "DomainService"

type DomainService struct {
	Creds creds.Client
}

func (c *Client) DomainNames() (*soap.Domains, error) {
	rawbody, e := soap.Lookup(c.Creds, soap.Request{Service: domainService, Method: "getDomainNames", Body: `<ns1:getDomainNames/>`})
	if e != nil {
		return nil, e
	}

	domains := &soap.Domains{}
	e := soap.Decode(rawbody, &domains)
	return domains, e
}

func (c *Client) Domain(name string) (*soap.Domains, error) {
	rawbody, e := soap.Lookup(c.Creds, soap.Request{
		Service: domainService,
		ExtraParams: []signature.KV{
			signature.KV{Key: "0", Value: name},
		},
		Method: "getInfo",
		Body:   fmt.Sprintf(`<ns1:getInfo><domainName xsi:type="xsd:string">%s</domainName></ns1:getInfo>`, name),
	})
	if e != nil {
		return nil, e
	}

	domain := &soap.Domain{}
	e := soap.Decode(rawbody, &domain)
	return domain, e
}
