package mp4

import (
	"fmt"
	"io"
)

// Box	Type:	 ‘traf’
// Container:	 Movie	Fragment	Box	('moof')
// Mandatory:	No
// Quantity:	 Zero	or	more
type TrafBox struct {
	Tfhd *TfhdBox   `json:"tfhd,omitempty"`
	Trun []*TrunBox `json:"trun,omitempty"`
	//Sbgp *SbgpBox `json:"sbgp,omitempty"`
	//Sgpd *SbgpBox `json:"sgpd,omitempty"`
	//Subs *SubsBox `json:"subs,omitempty"`
	//Saiz *SaizBox `json:"saiz,omitempty"`
	//Saiz *SaioBox `json:"saio,omitempty"`
	Tfdt  *TfdtBox `json:"tfdt,omitempty"`
	Meta  *MetaBox `json:"meta,omitempty"`
	Boxes []Box    `json:",omitempty"`
}

func (b *TrafBox) Box() Box {
	return b
}

func (b *TrafBox) Type() string {
	return "traf"
}

func (b *TrafBox) Size() int {
	l := BoxHeaderSize
	return l
}
func (b *TrafBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return err
}
func (b *TrafBox) Dump() {
	fmt.Printf("Track Fragment Box\n")
	if b.Tfhd != nil {
		b.Tfhd.Dump()
	}
	for _, t := range b.Trun {
		t.Dump()
	}
	if b.Tfdt != nil {
		b.Tfdt.Dump()
	}
	if b.Meta != nil {
		b.Meta.Dump()
	}
}

func DecodeTraf(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	t := &TrafBox{}
	for _, b := range l {
		switch b.Type() {
		case "tfhd":
			t.Tfhd = b.(*TfhdBox)
		case "tfdt":
			t.Tfdt = b.(*TfdtBox)
		case "trun":
			t.Trun = append(t.Trun, b.(*TrunBox))
		case "meta":
			t.Meta = b.(*MetaBox)
		default:
			t.Boxes = append(t.Boxes, b.Box())
		}
	}
	return t, err
}
