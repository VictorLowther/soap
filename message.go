package soap

import (
	"encoding/xml"
	"fmt"
	"github.com/VictorLowther/simplexml/dom"
	"io"
)

var (
	envName    xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Envelope"}
	headerName xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Header"}
	bodyName   xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Body"}
	faultName  xml.Name = xml.Name{Space: NS_ENVELOPE, Local: "Fault"}
)

type Message struct {
	Doc          *dom.Document
	header, body *dom.Element
}

func NewMessage() *Message {
	doc := dom.CreateDocument()
	root := dom.CreateElement(envName)
	doc.SetRoot(root)
	return &Message{Doc: doc}
}

func (m *Message) Headers() []*dom.Element {
	if m.header == nil {
		return []*dom.Element{}
	}
	return m.header.Children()
}

func (m *Message) AddHeader(e *dom.Element) {
	if m.header == nil {
		envelope := m.Doc.Root()
		m.header = dom.CreateElement(headerName)
		envelope.AddChild(m.header)
	}
	m.header.AddChild(e)
}

func (m *Message) RemoveHeader(e *dom.Element) *dom.Element {
	if m.header == nil {
		return nil
	}
	return m.header.RemoveChild(e)
}

func (m *Message) Body() []*dom.Element {
	if m.body == nil {
		return []*dom.Element{}
	}
	return m.body.Children()
}

func (m *Message) AddBody(e *dom.Element) {
	if m.body == nil {
		envelope := m.Doc.Root()
		m.body = dom.CreateElement(bodyName)
		envelope.AddChild(m.body)
	}
	m.body.AddChild(e)
}

func (m *Message) RemoveBody(e *dom.Element) *dom.Element {
	if m.body == nil {
		return nil
	}
	return m.body.RemoveChild(e)
}

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
	return &Message{Doc: doc, header: header, body: body}, nil
}
