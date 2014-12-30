package soap

import (
	"github.com/masterzen/simplexml/dom"
)

type Message struct {
	document *dom.Document
	envelope *dom.Element
	header   *Header
	body     *dom.Element
}

type MessageBuilder interface {
	SetBody(*dom.Element)
	NewBody() *dom.Element
	CreateElement(*dom.Element, string, dom.Namespace) *dom.Element
	CreateBodyElement(string, dom.Namespace) *dom.Element
	Header() *Header
	Doc() *dom.Document
	Free()

	String() string
}

func NewMessage() (message *Message) {
	doc := dom.CreateDocument()
	e := dom.CreateElement("Envelope")
	doc.SetRoot(e)
	AddUsualNamespaces(e)
	NS_SOAP_ENV.SetTo(e)

	message = &Message{document: doc, envelope: e}
	return
}

func (message *Message) NewBody() (body *dom.Element) {
	body = dom.CreateElement("Body")
	message.envelope.AddChild(body)
	NS_SOAP_ENV.SetTo(body)
	return
}

func (message *Message) String() string {
	return message.document.String()
}

func (message *Message) Doc() *dom.Document {
	return message.document
}

func (message *Message) Free() {
}

func (message *Message) CreateElement(parent *dom.Element, name string, ns dom.Namespace) (element *dom.Element) {
	element = dom.CreateElement(name)
	parent.AddChild(element)
	ns.SetTo(element)
	return
}

func (message *Message) CreateBodyElement(name string, ns dom.Namespace) (element *dom.Element) {
	if message.body == nil {
		message.body = message.NewBody()
	}
	return message.CreateElement(message.body, name, ns)
}

func (message *Message) Header() *Header {
	if message.header == nil {
		message.header = &Header{message: message}
	}
	return message.header
}
