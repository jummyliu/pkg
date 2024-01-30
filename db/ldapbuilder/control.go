package ldapbuilder

import (
	"fmt"

	ber "github.com/go-asn1-ber/asn1-ber"
)

// ControlInteger implements the Control interface for simple controls
type ControlInteger struct {
	ControlType  string
	Criticality  bool
	ControlValue int64
}

// GetControlType returns the OID
func (c *ControlInteger) GetControlType() string {
	return c.ControlType
}

// Encode returns the ber packet representation
func (c *ControlInteger) Encode() *ber.Packet {
	packet := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Control")
	packet.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, c.ControlType, "Control Type ("+c.ControlType+")"))
	if c.Criticality {
		packet.AppendChild(ber.NewBoolean(ber.ClassUniversal, ber.TypePrimitive, ber.TagBoolean, c.Criticality, "Criticality"))
	}

	// p2 ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, nil, "Control Value")
	p2 := ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, nil, "Control Value")
	value := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Control Value Sequence")
	value.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, c.ControlValue, "Integer"))
	p2.AppendChild(value)
	packet.AppendChild(p2)

	return packet
}

// String returns a human-readable description
func (c *ControlInteger) String() string {
	return fmt.Sprintf("Control Type: %v  Critiality: %t  Control Value: %v", c.ControlType, c.Criticality, c.ControlValue)
}

// NewControlInteger returns a generic control
func NewControlInteger(controlType string, criticality bool, controlValue int64) *ControlInteger {
	return &ControlInteger{
		ControlType:  controlType,
		Criticality:  criticality,
		ControlValue: controlValue,
	}
}

// NewControlSD LDAP_SERVER_SD_FLAGS_OID
func NewControlSDFlags() *ControlInteger {
	return &ControlInteger{
		ControlType:  "1.2.840.113556.1.4.801",
		Criticality:  true,
		ControlValue: int64(7),
	}
}
