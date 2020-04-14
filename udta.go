package mp4

import (
	"io"
	"log"
)

// User Data Box (udta - optional)
//
// Contained in: Movie Box (moov) or Track Box (trak)
type UdtaBox struct {
	Cprt []*CprtBox `json:"cprt,omitempty"`
	//Tsel *TselBox `json:"tsel,omitempty"`
	//Kind []*KindBox `json:"kind,omitempty"` // udta in a track
	//Strk []*StrkBox `json:"strk,omitempry"` // udts of trak
	Meta       *MetaBox `json:"meta,omitempty"` // ISO IEC 14496-12, does not contain meta in udta
	notDecoded []byte
}

//FIXME: Udta can contain several boxes, also actual user data
func DecodeUdta(h BoxHeader, r io.Reader) (Box, error) {
	if h.Size-BoxHeaderSize < BoxHeaderSize {
		data := make([]byte, h.Size-BoxHeaderSize)
		_, err := r.Read(data)
		if err != nil {
			return nil, err
		}
		log.Printf("udta Dump: %s", data)
		return &UdtaBox{
			notDecoded: data,
		}, nil

	}
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	u := &UdtaBox{}
	for _, b := range l {
		switch b.Type() {
		case "meta":
			u.Meta = b.(*MetaBox)
		case "cprt":
			u.Cprt = append(u.Cprt, b.(*CprtBox))
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
