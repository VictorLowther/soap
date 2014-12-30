package soap

import (
	"github.com/masterzen/simplexml/dom"
	"strconv"
)

type HeaderOption struct {
	key   string
	value string
}

func NewHeaderOption(name string, value string) *HeaderOption {
	return &HeaderOption{key: name, value: value}
}

type Header struct {
	to              string
	replyTo         string
	maxEnvelopeSize string
	timeout         string
	locale          string
	id              string
	action          string
	shellId         string
	resourceURI     string
	options         []HeaderOption
	message         *Message
}

type HeaderBuilder interface {
	To(string) *Header
	ReplyTo(string) *Header
	MaxEnvelopeSize(int) *Header
	Timeout(string) *Header
	Locale(string) *Header
	Id(string) *Header
	Action(string) *Header
	ShellId(string) *Header
	resourceURI(string) *Header
	AddOption(*HeaderOption) *Header
	Options([]HeaderOption) *Header
	Build(*Message) *Message
}

func (self *Header) To(uri string) *Header {
	self.to = uri
	return self
}

func (self *Header) ReplyTo(uri string) *Header {
	self.replyTo = uri
	return self
}

func (self *Header) MaxEnvelopeSize(size int) *Header {
	self.maxEnvelopeSize = strconv.Itoa(size)
	return self
}

func (self *Header) Timeout(timeout string) *Header {
	self.timeout = timeout
	return self
}

func (self *Header) Id(id string) *Header {
	self.id = id
	return self
}

func (self *Header) Action(action string) *Header {
	self.action = action
	return self
}

func (self *Header) Locale(locale string) *Header {
	self.locale = locale
	return self
}

func (self *Header) ShellId(shellId string) *Header {
	self.shellId = shellId
	return self
}

func (self *Header) ResourceURI(resourceURI string) *Header {
	self.resourceURI = resourceURI
	return self
}

func (self *Header) AddOption(option *HeaderOption) *Header {
	self.options = append(self.options, *option)
	return self
}

func (self *Header) Options(options []HeaderOption) *Header {
	self.options = options
	return self
}

func (self *Header) Build() *Message {
	header := self.createElement(self.message.envelope, "Header", NS_ENVELOPE)

	if self.to != "" {
		to := self.createElement(header, "To", NS_WSA)
		to.SetContent(self.to)
	}

	if self.replyTo != "" {
		replyTo := self.createElement(header, "ReplyTo", NS_WSA)
		a := self.createMUElement(replyTo, "Address", NS_WSA, true)
		a.SetContent(self.replyTo)
	}

	if self.maxEnvelopeSize != "" {
		envelope := self.createMUElement(header, "MaxEnvelopeSize", NS_WSMAN, true)
		envelope.SetContent(self.maxEnvelopeSize)
	}

	if self.timeout != "" {
		timeout := self.createElement(header, "OperationTimeout", NS_WSMAN)
		timeout.SetContent(self.timeout)
	}

	if self.id != "" {
		id := self.createElement(header, "MessageID", NS_WSA)
		id.SetContent(self.id)
	}

	if self.locale != "" {
		locale := self.createMUElement(header, "Locale", NS_WSMAN, false)
		locale.SetAttr("xml:lang", self.locale)
	}

	if self.action != "" {
		action := self.createMUElement(header, "Action", NS_WSA, true)
		action.SetContent(self.action)
	}

	if self.shellId != "" {
		selectorSet := self.createElement(header, "SelectorSet", NS_WSMAN)
		selector := self.createElement(selectorSet, "Selector", NS_WSMAN)
		selector.SetAttr("Name", "ShellId")
		selector.SetContent(self.shellId)
	}

	if self.resourceURI != "" {
		resource := self.createMUElement(header, "ResourceURI", NS_WSMAN, true)
		resource.SetContent(self.resourceURI)
	}

	if len(self.options) > 0 {
		set := self.createElement(header, "OptionSet", NS_WSMAN)
		for _, option := range self.options {
			e := self.createElement(set, "Option", NS_WSMAN)
			e.SetAttr("Name", option.key)
			e.SetContent(option.value)
		}
	}

	return self.message
}

func (self *Header) createElement(parent *dom.Element, name string, ns dom.Namespace) (element *dom.Element) {
	element = dom.CreateElement(name)
	parent.AddChild(element)
	ns.SetTo(element)
	return
}

func (self *Header) createMUElement(parent *dom.Element, name string, ns dom.Namespace, mustUnderstand bool) (element *dom.Element) {
	element = self.createElement(parent, name, ns)
	value := "false"
	if mustUnderstand {
		value = "true"
	}
	element.SetAttr("mustUnderstand", value)
	return
}
