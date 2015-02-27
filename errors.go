package soap

const (
	BadEnvelope         = "Invalid SOAP: Root should be an Envelope"
	BadTag              = "Invalid SOAP: Unexpected tag in Envelope"
	EnvelopeOverstuffed = "Invalid SOAP: Envelope must have at most 2 children"
	NoEnvelope          = "Invalid SOAP: Document does not have a root element"
	TooManyBody         = "Invalid SOAP: More than one Body"
	TooManyHeader       = "Invalid SOAP: More than one Header"
)
