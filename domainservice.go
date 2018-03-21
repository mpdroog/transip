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

func (c *DomainService) DomainNames() (*soap.DomainNames, error) {
	rawbody, e := soap.Lookup(c.Creds, soap.Request{Service: domainService, Method: "getDomainNames", Body: `<ns1:getDomainNames/>`})
	if e != nil {
		return nil, e
	}

	domains := &soap.DomainNames{}
	e = soap.Decode(rawbody, &domains)
	return domains, e
}

func (c *DomainService) Domain(name string) (*soap.Domain, error) {
	rawbody, e := soap.Lookup(c.Creds, soap.Request{
		Service: domainService,
		ExtraParams: []signature.KV{
			{Key: "0", Value: name},
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

func (c *DomainService) Domains(names []string) (*soap.Domains, error) {
    entryTemplate := `<item xsi:type="xsd:string">%s</item>`
    params := []signature.KV{}
    xml := ``

    for idx, v := range names {
        xml = xml + fmt.Sprintf(entryTemplate, v)
        params = append(params, []signature.KV{
            {Key: fmt.Sprintf("0[%d]", idx), Value: v},
        }...)
    }

	rawbody, e := soap.Lookup(c.Creds, soap.Request{
		Service: domainService,
		ExtraParams: params,
		Method: "batchGetInfo",
        Body:   fmt.Sprintf(`<ns1:batchGetInfo><domainNames SOAP-ENC:arrayType="xsd:string[%d]" xsi:type="ns1:ArrayOfString">%s</domainNames></ns1:batchGetInfo>`, len(names), xml),
	})
	if e != nil {
		return nil, e
	}

	domains := &soap.Domains{}
	e = soap.Decode(rawbody, &domains)
	return domains, e
}

func (c *DomainService) SetDNSEntries(domain string, entries []soap.DomainDNSentry) error {
	entryTemplate := `<item xsi:type="ns1:DnsEntry"><name xsi:type="xsd:string">%s</name><expire xsi:type="xsd:int">%d</expire><type xsi:type="xsd:string">%s</type><content xsi:type="xsd:string">%s</content></item>`

	params := []signature.KV{
		{Key: "0", Value: domain},
	}
	xml := ``

	for idx, entry := range entries {
		xml = xml + fmt.Sprintf(entryTemplate, entry.Name, entry.Expire, entry.Type, entry.Content)
		params = append(params, []signature.KV{
			{fmt.Sprintf("1[%d][name]", idx), entry.Name},
			{fmt.Sprintf("1[%d][expire]", idx), strconv.Itoa(entry.Expire)},
			{fmt.Sprintf("1[%d][type]", idx), entry.Type},
			{fmt.Sprintf("1[%d][content]", idx), entry.Content},
		}...)
	}

	rawbody, e := soap.Lookup(c.Creds, soap.Request{
		Service:     domainService,
		ExtraParams: params,
		Method:      "setDnsEntries",
		Body: fmt.Sprintf(
			`<ns1:setDnsEntries><domainName xsi:type="xsd:string">%s</domainName><dnsEntries SOAP-ENC:arrayType="ns1:DnsEntry[%d]" xsi:type="ns1:ArrayOfDnsEntry">%s</dnsEntries></ns1:setDnsEntries>`,
			domain, len(entries), xml,
		),
	})
	if e != nil {
		return e
	}

	if !bytes.Contains(rawbody, []byte(`<ns1:setDnsEntriesResponse/>`)) {
		return errors.New("Unexpected XML-reply: " + string(rawbody))
	}
	return nil
}
