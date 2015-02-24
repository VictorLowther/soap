// Package soap wraps github.com/VictorLowther/simplexml/dom to provide
// convienent methods for dealing with SOAP messages as a SOAP client.
package soap

import (
	"encoding/xml"
	"fmt"
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
	Doc          *dom.Document
	Header, Body *dom.Element
}

// NewMessage creates the skeleton of a new SOAP message.
func NewMessage() *Message {
	res := &Message{
		Doc:    dom.CreateDocument(),
		Body:   dom.CreateElement(bodyName),
		Header: dom.CreateElement(headerName),
	}
	res.Doc.SetRoot(dom.CreateElement(envName))
	res.Doc.Root().AddChild(res.Header)
	res.Doc.Root().AddChild(res.Body)
	return res
}

// IsSoap takes a simplexml dom.Document and validates that
// it contains a valid SOAP message.  If it does, it returns a Message.
// If it does not, it returns an error explaining why not.
func IsSoap(doc *dom.Document) (res *Message, err error) {
	envelope := doc.Root()
	if envelope == nil {
		return nil, fmt.Errorf("Invalid SOAP: Document does not have a root element")
	}
	if envelope.Name != envName {
		return nil, fmt.Errorf("Invalid SOAP: Root should be %v, not '%v'", envName, envelope.Name)
	}
	children := envelope.Children()
	if len(children) > 2 {
		return nil, fmt.Errorf("Invalid SOAP: Envelope must have at most 2 children, not %d",
			len(children))
	}
	var header, body *dom.Element
	for _, c := range children {
		if c.Name == headerName {
			if header != nil {
				return nil, fmt.Errorf("Invalid SOAP: More than one Header element!")
			}
			header = c
		} else if c.Name == bodyName {
			if body != nil {
				return nil, fmt.Errorf("Invalid SOAP: More than one Body element!")
			}
			body = c
		} else {
			return nil, fmt.Errorf("Invalid SOAP: Unexpected tag %v", c.Name)
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
	return &Message{Doc: doc, Header: header, Body: body}, nil
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
