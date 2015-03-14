package soap

/*
Copyright 2015 Victor Lowther <victor.lowther@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/VictorLowther/simplexml/dom"
	"github.com/VictorLowther/simplexml/search"
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

func parseSoap() *Message {
	msg, err := Parse(strings.NewReader(simpleSoap))
	if err != nil {
		panic("Cannot parse test document.")
	}
	return msg
}

func TestSoapGetHeadersAndBody(t *testing.T) {
	msg := parseSoap()
	headers := msg.Headers()
	body := msg.Body()
	if len(headers) != 2 {
		t.Errorf("Expected 2 headers, got %d", len(headers))
	}
	if headers[0].Name.Local != "Action" {
		t.Errorf("Expected first header to be Action, not %v", headers[0].Name)
	}
	if s := string(headers[0].Content); s != "Frob" {
		t.Errorf("Expected Action header to be Frob, not %s", s)
	}
	if headers[1].Name.Local != "To" {
		t.Errorf("Expected second header to be To, not %v", headers[0].Name)
	}
	if len(body) != 1 {
		t.Errorf("Expected 1 body element, got %d", len(body))
	}
}

func TestSoapGetHeaders(t *testing.T) {
	msg := parseSoap()
	actionHdr := MuElem("Action","")
	toHdr := MuElem("To","")
	action := msg.GetHeader(actionHdr)
	if action == nil {
		t.Errorf("Expected to get %v, got nil!",actionHdr.Name)
	}
	if action.Name != actionHdr.Name {
		t.Errorf("Expected to get %v, got %v",actionHdr.Name, action.Name)
	}
	to := msg.GetHeader(toHdr)
	if to == nil {
		t.Errorf("Expected to get %v, got nil!",toHdr.Name)
	}
	if to.Name != toHdr.Name {
		t.Errorf("Expected to get %v, got %v",toHdr.Name, to.Name)
	}
}

func TestSoapGetBody(t *testing.T) {
	msg := parseSoap()
	frobBody := dom.Elem("Frob","")
	frob := msg.GetBody(frobBody)
	if frob == nil {
		t.Errorf("Expected to get %v, got nil!",frobBody.Name)
	}
	if frob.Name != frobBody.Name {
		t.Errorf("Expected to get %v, got %v",frobBody.Name, frob.Name)
	}
}


func TestSoapAddAndRemoveHeaders(t *testing.T) {
	msg := parseSoap()
	msg.SetHeader(
		MuElemC("ReplyTo", "", "Me"),
		MuElemC("AbandonAll", "", "Hope"),
		MuElemC("Action", "", "Abandon"))
	headers := msg.Headers()
	if len(headers) != 4 {
		t.Errorf("Expected 4 headers, got %d", len(headers))
	}
	if headers[0].Name.Local != "Action" {
		t.Errorf("Expected first header to be Action, not %v", headers[0].Name)
	}
	if s := string(headers[0].Content); s != "Abandon" {
		t.Errorf("Expected Action header to be Ababdon, not %s", s)
	}
	if headers[2].Name.Local != "ReplyTo" {
		t.Errorf("Expected third header to be ReplyTo, not %v", headers[2].Name)
	}
	if headers[3].Name.Local != "AbandonAll" {
		t.Errorf("Expected fourth header to be AbandonAll, not %v", headers[3].Name)
	}
	removed := msg.RemoveHeader(headers[2])
	if removed == nil {
		t.Errorf("Failed to remove %v", headers[2].Name)
	}
	if removed.Name.Local != "ReplyTo" {
		t.Errorf("Expected removed header to be ReplyTo, not %v", removed.Name)
	}
	headers = msg.Headers()
	if len(headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(headers))
	}
	removed = search.First(search.Tag("ReplyTo", ""), headers)
	if removed != nil {
		t.Errorf("Did not expect to find %v in the SOAP headers", removed.Name)
	}
}

func TestSoapAddAndRemoveBody(t *testing.T) {
	msg := parseSoap()
	msg.SetBody(
		MuElemC("ReplyTo", "", "Me"),
		MuElemC("AbandonAll", "", "Hope"),
		MuElemC("Frob", "", "Abandon"))
	body := msg.Body()
	if len(body) != 3 {
		t.Errorf("Expected 3 body elements, got %d", len(body))
	}
	if body[0].Name.Local != "Frob" {
		t.Errorf("Expected first body to be Frob, not %v", body[0].Name)
	}
	if s := string(body[0].Content); s != "Abandon" {
		t.Errorf("Expected Frob body to be Abandon, not %s", s)
	}
	if body[1].Name.Local != "ReplyTo" {
		t.Errorf("Expected third body to be ReplyTo, not %v", body[1].Name)
	}
	if body[2].Name.Local != "AbandonAll" {
		t.Errorf("Expected fourth body to be AbandonAll, not %v", body[2].Name)
	}
	removed := msg.RemoveBody(body[1])
	if removed == nil {
		t.Errorf("Failed to remove %v", body[1].Name)
	}
	if removed.Name.Local != "ReplyTo" {
		t.Errorf("Expected removed body to be ReplyTo, not %v", removed.Name)
	}
	body = msg.Body()
	if len(body) != 2 {
		t.Errorf("Expected 2 body elements, got %d", len(body))
	}
	removed = search.First(search.Tag("ReplyTo", ""), body)
	if removed != nil {
		t.Errorf("Did not expect to find %v in the SOAP body", removed.Name)
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
	msg.body.Name = headerName
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != TooManyHeader {
		t.Errorf("IsSoap should have failed with TooManyHeader, got %v", err)
	}
}

func TestSoapTooManyBody(t *testing.T) {
	msg := NewMessage()
	msg.header.Name = bodyName
	_, err := IsSoap(msg.Document)
	if err == nil || err.Error() != TooManyBody {
		t.Errorf("IsSoap should have failed with TooManyBody, got %v", err)
	}
}

func TestSoapBadTag(t *testing.T) {
	msg := NewMessage()
	msg.header.Name = faultName
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
	if msg.header == nil || msg.body == nil {
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
