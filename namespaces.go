package soap

import (
	"github.com/masterzen/simplexml/dom"
	"gopkg.in/xmlpath.v2"
)

// Default WSMAN namespaces.  Taked from appendix 1 of
// http://www.dmtf.org/sites/default/files/standards/documents/DSP0226_1.2.0.pdf
var (
	NS_WSMAN       = dom.Namespace{"wsman", "http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd"}
	NS_WSMID       = dom.Namespace{"wsmid", "http://schemas.dmtf.org/wbem/wsman/identity/1/wsmanidentity.xsd"}
	NS_ENVELOPE    = dom.Namespace{"s", "http://www.w3.org/2003/05/soap-envelope"}
	NS_SCHEMA      = dom.Namespace{"xs", "http://www.w3.org/2001/XMLSchema"}
	NS_SCHEMA_INST = dom.Namespace{"xsi", "http://www.w3.org/2001/XMLSchema-instance"}
	NS_WSDL        = dom.Namespace{"wsdl", "http://schemas.xmlsoap.org/wsdl"}
	NS_WSA04       = dom.Namespace{"wsa04", "http://schemas.xmlsoap.org/ws/2004/08/addressing"}
	NS_WSA10       = dom.Namespace{"wsa10", "http://www.w3.org/2005/08/addressing"}
	NS_WSA         = NS_WSA04
	NS_WSAM        = dom.Namespace{"wsam", "http://www.w3.org/2007/5/addressing/metadata"}
	NS_WSME        = dom.Namespace{"wsme", "http://schemas.xmlsoap.org/ws/2004/08/eventing"}
	NS_WSMEN       = dom.Namespace{"wsmen", "http://schemas.xmlsoap.org/ws/2004/09/enumeration"}
	NS_WSMT        = dom.Namespace{"wsmt", "http://schemas.xmlsoap.org/ws/2004/09/transfer"}
	NS_WSP         = dom.Namespace{"wsp", "http://schemas.xmlsoap.org/ws/2004/09/policy"}
)

var MostUsed = [...]dom.Namespace{NS_ENVELOPE, NS_WSA, NS_WSMAN}

func AddUsualNamespaces(node *dom.Element) {
	for _, ns := range MostUsed {
		node.DeclareNamespace(ns)
	}
}

func GetAllNamespaces() []xmlpath.Namespace {
	var ns = []dom.Namespace{NS_ENVELOPE, NS_WSA, NS_WSMAN}

	var xmlpathNs = make([]xmlpath.Namespace, 0, 4)
	for _, namespace := range ns {
		xmlpathNs = append(xmlpathNs, xmlpath.Namespace{Prefix: namespace.Prefix, Uri: namespace.Uri})
	}
	return xmlpathNs
}
