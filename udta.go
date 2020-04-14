package mp4

import "io"

// User Data Box (udta - optional)
//
// Contained in: Movie Box (moov) or Track Box (trak)
type UdtaBox struct {
	//Cprt []*CprtBox `json:"cprt,omitempty"`
	//Tsel *TselBox `json:"tsel,omitempty"`
	//Kind []*KindBox `json:"kind,omitempty"` // udta in a track
	//Strk []*StrkBox `json:"strk,omitempry"` // udts of trak
	Meta *MetaBox `json:"meta,omitempty"` // ISO IEC 14496-12, does not contain meta in udta
}

func DecodeUdta(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	u := &UdtaBox{}
	for _, b := range l {
		switch b.Type() {
		case "meta":
			u.Meta = b.(*MetaBox)
		}
	}
	return u, nil
}

func (b *UdtaBox) Box() Box {
	return b
}

func (b *UdtaBox) Type() string {
	return "udta"
}

func (b *UdtaBox) Size() int {
	l := BoxHeaderSize
	if b.Meta != nil {
		l = l + b.Meta.Size()
	}
	return l
}

func (b *UdtaBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return b.Meta.Encode(w)
}
