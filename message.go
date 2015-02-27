// Package soap wraps github.com/VictorLowther/simplexml/dom to provide
// convienent methods for dealing with SOAP messages as a SOAP client.
package soap

import (
	"encoding/xml"
	"errors"
	"github.com/VictorLowther/simplexml/dom"
	"io"
)

const ContentType = "application/soap+xml; charset=utf-8"

var (
	envName    xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Envelope"}
	headerName xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Header"}
	bodyName   xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Body"}
	faultName  xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Fault"}
)

// Message represents a SOAP message.
type Message struct {
	*dom.Document
	Header, Body *dom.Element
}

// NewMessage creates the skeleton of a new SOAP message.
func NewMessage() *Message {
	res := &Message{
		Document: dom.CreateDocument(),
		Body:     dom.CreateElement(bodyName),
		Header:   dom.CreateElement(headerName),
	}
	res.SetRoot(dom.CreateElement(envName))
	res.Root().AddChild(res.Header)
	res.Root().AddChild(res.Body)
	return res
}

// IsSoap takes a simplexml dom.Document and validates that
// it contains a valid SOAP message.  If it does, it returns a Message.
// If it does not, it returns an error explaining why not.
func IsSoap(doc *dom.Document) (res *Message, err error) {
	envelope := doc.Root()
	if envelope == nil {
		return nil, errors.New(NoEnvelope)
	}
	if envelope.Name != envName {
		return nil, errors.New(BadEnvelope)
	}
	children := envelope.Children()
	if len(children) > 2 {
		return nil, errors.New(EnvelopeOverstuffed)
	}
	var header, body *dom.Element
	for _, c := range children {
		if c.Name == headerName {
			if header != nil {
				return nil, errors.New(TooManyHeader)
			}
			header = c
		} else if c.Name == bodyName {
			if body != nil {
				return nil, errors.New(TooManyBody)
			}
			body = c
		} else {
			return nil, errors.New(BadTag)
		}
	}
	if header == nil {
		header = dom.CreateElement(headerName)
		doc.Root().AddChild(header)
	}
	if body == nil {
		body = dom.CreateElement(bodyName)
		doc.Root().AddChild(body)
	}
	return &Message{Document: doc, Header: header, Body: body}, nil
}

// Parse parses what is hopefully a well-formed SOAP message
// from the passed io.Reader.  If it is not, err will say why not.
func Parse(r io.Reader) (msg *Message, err error) {
	doc, err := dom.Parse(r)
	if err != nil {
		return nil, err
	}
	msg, err = IsSoap(doc)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
