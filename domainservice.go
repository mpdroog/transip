package transip

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mpdroog/transip/creds"
	"github.com/mpdroog/transip/soap"
	"github.com/mpdroog/transip/soap/signature"
	"strconv"
)

const domainService = "DomainService"

type DomainService struct {
	Creds creds.Client
}

func (c *DomainService) DomainNames() (*soap.Domains, error) {
	rawbody, e := soap.Lookup(c.Creds, soap.Request{Service: domainService, Method: "getDomainNames", Body: `<ns1:getDomainNames/>`})
	if e != nil {
		return nil, e
	}

	domains := &soap.Domains{}
	e = soap.Decode(rawbody, &domains)
	return domains, e
}

func (c *DomainService) Domain(name string) (*soap.Domain, error) {
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
	e = soap.Decode(rawbody, &domain)
	return domain, e
}

func (c *DomainService) SetDNSEntries(domain string, entries []soap.DomainDNSentry) error {
	entryTemplate := `<item xsi:type="ns1:DnsEntry"><name xsi:type="xsd:string">%s</name><expire xsi:type="xsd:int">%d</expire><type xsi:type="xsd:string">%s</type><content xsi:type="xsd:string">%s</content></item>`

	params := []signature.KV{
		signature.KV{Key: "0", Value: domain},
	}
	xml := ``

	for idx, entry := range entries {
		offset := idx + 1
		xml = xml + fmt.Sprintf(entryTemplate, entry.Name, entry.Expire, entry.Type, entry.Content)
		params = append(params, []signature.KV{
			signature.KV{fmt.Sprintf("%d[0][name]", offset), entry.Name},
			signature.KV{fmt.Sprintf("%d[0][expire]", offset), strconv.Itoa(entry.Expire)},
			signature.KV{fmt.Sprintf("%d[0][type]", offset), entry.Type},
			signature.KV{fmt.Sprintf("%d[0][content]", offset), entry.Content},
		}...)
	}

	rawbody, e := soap.Lookup(c.Creds, soap.Request{
		Service:     domainService,
		ExtraParams: params,
		Method:      "setDnsEntries",
		Body:        fmt.Sprintf(`<ns1:setDnsEntries><domainName xsi:type="xsd:string">%s</domainName><dnsEntries SOAP-ENC:arrayType="ns1:DnsEntry[1]" xsi:type="ns1:ArrayOfDnsEntry">%s</dnsEntries></ns1:setDnsEntries>`, domain, xml),
	})
	if e != nil {
		return e
	}

	if !bytes.Contains(rawbody, []byte(`<ns1:setDnsEntriesResponse/>`)) {
		return errors.New("Unexpected XML-reply: " + string(rawbody))
	}
	return nil
}
