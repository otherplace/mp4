package mp4

import "io"

// Track Box (tkhd - mandatory)
//
// Contained in : Movie Box (moov)
//
// A media file can contain one or more tracks.
type TrakBox struct {
	Tkhd *TkhdBox `json:"tkhd,"`
	//Tref *TrefBox `json:"tref,omitempty"`
	//Trgr *TrgrBox `json"trgr,omitempty"`
	Mdia  *MdiaBox `json:"mdia,"`
	Edts  *EdtsBox `json:"edts,omitempty"`
	Meta  *MetaBox `json:"meta,omitempty"`
	Boxes []Box
}

func DecodeTrak(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	t := &TrakBox{}
	for _, b := range l {
		switch b.Type() {
		case "tkhd":
			t.Tkhd = b.(*TkhdBox)
		case "mdia":
			t.Mdia = b.(*MdiaBox)
		case "edts":
			t.Edts = b.(*EdtsBox)
		case "meta":
			t.Meta = b.(*MetaBox)
		default:
			t.Boxes = append(t.Boxes, b.Box())
		}
	}
	return t, nil
}

func (b *TrakBox) Box() Box {
	return b
}

func (b *TrakBox) Type() string {
	return "trak"
}

func (b *TrakBox) Size() int {
	sz := b.Tkhd.Size()
	sz += b.Mdia.Size()
	if b.Edts != nil {
		sz += b.Edts.Size()
	}
	if b.Meta != nil {
		sz += b.Meta.Size()
	}
	return sz + BoxHeaderSize
}

func (b *TrakBox) Dump() {
	b.Tkhd.Dump()
	if b.Edts != nil {
		b.Edts.Dump()
	}
	b.Mdia.Dump()
	if b.Meta != nil {
		b.Meta.Dump()
	}
}

func (b *TrakBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	err = b.Tkhd.Encode(w)
	if err != nil {
		return err
	}
	if b.Edts != nil {
		err = b.Edts.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Meta != nil {
		b.Meta.Encode(w)
	}
	return b.Mdia.Encode(w)
}
