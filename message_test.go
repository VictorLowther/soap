package soap

import (
	"github.com/VictorLowther/simplexml/dom"
	"strings"
	"testing"
)

var simpleSoap string = `<?xml version="1.0"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope">
  <s:Header>
    <Action s:mustUnderstand="true">Frob</Action>
    <To s:mustUnderstand="true">this</To>
  </s:Header>
  <s:Body>
    <Frob>
      <OptimizeFrobnication/>
      <MaxFrobs>999</MaxFrobs>
    </Frob>
  </s:Body>
</s:Envelope>`

func TestSkeletonSoap(t *testing.T) {
	msg := NewMessage()
	_, err := IsSoap(msg.Document)
	if err != nil {
		t.Fatalf("Skeleton SOAP generator not generating valid SOAP\nGot: %s\n\nError: %v", msg.String(), err)
	}
}

func TestSimpleSoap(t *testing.T) {
	_, err := Parse(strings.NewReader(simpleSoap))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
}

func TestSoapNoEnvelope(t *testing.T) {
	_, err := Parse(strings.NewReader(`<?xml version="1.0"?>`))
	if err == nil || err.Error() != NoEnvelope {
		t.Errorf("IsSoap should have failed with NoEnvelope, got %v", err)
	}
}

func TestSoapBadEnvelope(t *testing.T) {
	msg := NewMessage()
	msg.SetRoot(dom.Elem("BadEnvelope", ""))
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != BadEnvelope {
		t.Errorf("IsSoap should have failed with BadEnvelope, got %v", err)
	}
}

func TestSoapEnvelopeOverstuffed(t *testing.T) {
	msg := NewMessage()
	msg.Root().AddChild(dom.Elem("ExtraThing", ""))
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != EnvelopeOverstuffed {
		t.Errorf("IsSoap should have failed with EnvelopeOverstuffed, got %v", err)
	}
}

func TestSoapTooManyHeader(t *testing.T) {
	msg := NewMessage()
	msg.Body.Name = headerName
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != TooManyHeader {
		t.Errorf("IsSoap should have failed with TooManyHeader, got %v", err)
	}
}

func TestSoapTooManyBody(t *testing.T) {
	msg := NewMessage()
	msg.Header.Name = bodyName
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != TooManyBody {
		t.Errorf("IsSoap should have failed with TooManyBody, got %v", err)
	}
}

func TestSoapBadTag(t *testing.T) {
	msg := NewMessage()
	msg.Header.Name = faultName
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != BadTag {
		t.Errorf("IsSoap should have failed with BadTag, got %v", err)
	}
}

func TestSoapAddHeaderAndBody(t *testing.T) {
	doc := dom.CreateDocument()
	doc.SetRoot(dom.CreateElement(envName))
	msg, err := IsSoap(doc)
	if err != nil {
		t.Error("IsSoap should have passed")
	}
	if msg.Header == nil || msg.Body == nil {
		t.Error("IsSoap failed to add a header and a body")
	}
}

var soapFault string = `<?xml version="1.0"?>
<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope"
            xmlns:xml="http://www.w3.org/XML/1998/namespace">
  <s:Body>
    <s:Fault>
      <s:Code>
        <s:Value>s:Sender</s:Value>
        <s:Subcode>
          <s:Value>It died.</s:Value>
        </s:Subcode>
      </s:Code>
      <s:Reason>
        <s:Text xml:lang="en">Death By Chocolate</s:Text>
      </s:Reason>
      <s:Detail>
        <maxChoc>5 bars</maxChoc>
      </s:Detail>
    </s:Fault>
  </s:Body>
</s:Envelope>`

func TestSoapFault(t *testing.T) {
	msg, err := Parse(strings.NewReader(soapFault))
	if err != nil {
		t.Errorf("Error %v parsing %s", err, soapFault)
	}
	if f := msg.Fault(); f == nil {
		t.Error("Expected message to be a SOAP Fault!")
	}
}
