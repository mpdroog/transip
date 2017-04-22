// Package soap implements SOAP-logic for the TransIP API.
// decode.go contains the logic to convert the XML to the corresponding
// datastructures
package soap

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type Domains struct {
	Item []string `xml:"item"`
}

type DomainNameserver struct {
	Hostname string `xml:"hostname"`
	IPv4     string `xml:"ipv4"`
	IPv6     string `xml:"ipv6"`
}
type DomainContact struct {
	Type        string `xml:"type"`
	FirstName   string `xml:"firstName"`
	MiddleName  string `xml:"middleName"`
	LastName    string `xml:"lastName"`
	CompanyName string `xml:"companyName"`
	CompanyKvK  string `xml:"companyKvk"`
	CompanyType string `xml:"companyType"`
	Street      string `xml:"street"`
	Number      string `xml:"number"`
	PostalCode  string `xml:"postalCode"`
	City        string `xml:"city"`
	PhoneNumber string `xml:"phoneNumber"`
	FaxNumber   string `xml:"faxNumber"`
	Email       string `xml:"email"`
	Country     string `xml:"country"`
}
type DomainDNSentry struct {
	Name    string `xml:"name"`
	Expire  int    `xml:"expire"`
	Type    string `xml:"type"`
	Content string `xml:"content"`
}
type Domain struct {
	Name        string             `xml:"name"`
	Nameservers []DomainNameserver `xml:"nameservers>item"`
	Contacts    []DomainContact    `xml:"contacts>item"`
	DNSEntry    []DomainDNSentry   `xml:"dnsEntries>item"`

	AuthCode         string `xml:"authCode"`
	IsLocked         bool   `xml:"isLocked"`
	RegistrationDate string `xml:"registrationDate"`
	RenewalDate      string `xml:"renewalDate"`
}

// Convert rawbody to XML and subtract the 'body' from the
// SOAP-envelope into the struct given with out
func Decode(rawbody []byte, out interface{}) error {
	dec := xml.NewDecoder(bytes.NewReader(rawbody))

	for {
		tok, e := dec.Token()
		if e == io.EOF {
			return nil
		}
		if e != nil {
			return e
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "return" {
				// Start readin'!
				if e := dec.DecodeElement(out, &se); e != nil {
					return e
				}
			}
			if se.Name.Local == "Fault" {
				// Error!
				err := &SOAPFault{}
				if e := dec.DecodeElement(err, &se); e != nil {
					return e
				}
				return errors.New("soapFault: " + err.String)
			}

		}
	}

	return nil
}
