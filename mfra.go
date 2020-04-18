package mp4

import (
	"fmt"
	"io"
)

// Box	Type:	 ‘mfra’
// Container:	 File
// Mandatory:	No
// Quantity:	 Zero	or	one
type MfraBox struct {
	Tfra  *TfraBox `json:"tfra,omitempty"`
	Mfro  *MfroBox `json:"mfro,"`
	Boxes []Box    `json:",omitempty"`
}

func DecodeMfra(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MfraBox{}
	for _, b := range l {
		switch b.Type() {
		case "tfra":
			m.Tfra = b.(*TfraBox)
		case "mfro":
			m.Mfro = b.(*MfroBox)
		default:
			m.Boxes = append(m.Boxes, b.Box())
		}
	}
	return m, err
}

func (b *MfraBox) Box() Box {
	return b
}

func (b *MfraBox) Type() string {
	return "mfra"
}

func (b *MfraBox) Size() int {
	l := BoxHeaderSize
	if b.Tfra != nil {
		l += b.Tfra.Size()
	}
	if b.Mfro != nil {
		l += b.Mfro.Size()
	}
	return l
}

func (b *MfraBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	if b.Tfra != nil {
		err = b.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Mfro != nil {
		err = b.Encode(w)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *MfraBox) Dump() {
	fmt.Printf("Movie Fragment Random Access Box\n")
	if b.Tfra != nil {
		b.Tfra.Dump()
	}
	if b.Mfro != nil {
		b.Mfro.Dump()
	}
}
