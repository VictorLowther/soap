package soap

import (
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
	_, err := IsSoap(msg.Doc)
	if err != nil {
		t.Fatalf("Skeleton SOAP generator not generating valid SOAP\nGot: %s\n\nError: %v", msg.Doc.String(), err)
	}
}

func TestSimpleSoap(t *testing.T) {
	_, err := Parse(strings.NewReader(simpleSoap))
	if err != nil {
		t.Fatalf("Cannot parse test document. Error: %v", err)
	}
}
