package mp4

import (
	"fmt"
	"io"
)

// Movie Box (moov - mandatory)
//
// Status: partially decoded
//
// Contains all meta-data. To be able to stream a file, the moov box should be placed before the mdat box.
type MoovBox struct {
	Mvhd  *MvhdBox   `json:"mvhd"`
	Iods  *IodsBox   `json:"iods,omitempty"` // ISO IEC 13396-12, does not contain this box in moov box
	Trak  []*TrakBox `json:"trak"`
	Mvex  *MvexBox   `json:"mvex,omitempty"`
	Udta  *UdtaBox   `json:"udta,omitempty"`
	Meta  *MetaBox   `json:"meta,omitempty"`
	Boxes []Box      `json:",omitempty"`
}

func DecodeMoov(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MoovBox{}
	for _, b := range l {
		switch b.Type() {
		case "mvhd":
			m.Mvhd = b.(*MvhdBox)
		case "iods":
			m.Iods = b.(*IodsBox)
		case "trak":
			m.Trak = append(m.Trak, b.(*TrakBox))
		//case "udta":
		//	m.Udta = b.(*UdtaBox)
		case "mvex":
			m.Mvex = b.(*MvexBox)
		case "meta":
			m.Meta = b.(*MetaBox)
		default:
			m.Boxes = append(m.Boxes, b.Box())
		}
	}
	return m, err
}

func (b *MoovBox) Box() Box {
	return b
}

func (b *MoovBox) Type() string {
	return "moov"
}

func (b *MoovBox) Size() int {
	sz := b.Mvhd.Size()
	if b.Iods != nil {
		sz += b.Iods.Size()
	}
	for _, t := range b.Trak {
		sz += t.Size()
	}
	if b.Udta != nil {
		sz += b.Udta.Size()
	}
	if b.Mvex != nil {
		sz += b.Mvex.Size()
	}
	if b.Meta != nil {
		sz += b.Meta.Size()
	}
	return sz + BoxHeaderSize
}

func (b *MoovBox) Dump() {
	b.Mvhd.Dump()
	for i, t := range b.Trak {
		fmt.Println("Track", i)
		t.Dump()
	}
	if b.Mvex != nil {
		b.Mvex.Dump()
	}
}

func (b *MoovBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	err = b.Mvhd.Encode(w)
	if err != nil {
		return err
	}
	if b.Iods != nil {
		err = b.Iods.Encode(w)
		if err != nil {
			return err
		}
	}
	for _, t := range b.Trak {
		err = t.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Udta != nil {
		err = b.Udta.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Mvex != nil {
		err = b.Mvex.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Meta != nil {
		err = b.Meta.Encode(w)
		if err != nil {
			return err
		}
	}

	return nil
}
