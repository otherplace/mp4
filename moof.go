package mp4

import (
	"fmt"
	"io"
)

// Box	Type:	 ‘moof’
// Container:	 File
// Mandatory:	No
// Quantity:	 Zero	or	more
type MoofBox struct {
	Mfhd *MfhdBox   `json:"mfhd,"`
	Traf []*TrafBox `json:"traf,omitempty"`
}

func (b *MoofBox) Type() string {
	return "moof"
}

func (b *MoofBox) Size() int {
	l := BoxHeaderSize + b.Mfhd.Size()
	for _, e := range b.Traf {
		l = l + e.Size()
	}
	return l
}
func (b *MoofBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return err
}
func (b *MoofBox) Dump() {
	fmt.Printf("Movie Fragment Box\n")
}

func DecodeMoof(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MoofBox{}
	for _, b := range l {
		switch b.Type() {
		case "mfhd":
			m.Mfhd = b.(*MfhdBox)
		case "traf":
			m.Traf = append(m.Traf, b.(*TrafBox))
		}
	}
	return m, err
}
