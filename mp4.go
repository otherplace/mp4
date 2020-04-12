package mp4

import (
	"io"
)

// A MPEG-4 media
//
// A MPEG-4 media contains three main boxes :
//
//   ftyp : the file type box
//   moov : the movie box (meta-data)
//   mdat : the media data (chunks and samples)
//
// Other boxes can also be present (pdin, moof, mfra, free, ...), but are not decoded.
type MP4 struct {
	Ftyp *FtypBox
	//Pdin  *PdinBox
	Moov *MoovBox
	//Moof  *MoofBox
	//Mfra  *MfraBox
	Mdat *MdatBox
	Free *FreeBox
	//Skip  *SkipBox
	//Meta  *MetaBox
	//Meco  *MecoBox
	Styp *StypBox
	//Sidx  *SidxBox
	//Ssix  *SSixBox
	//Prft  *PrftBox
	boxes []Box
}

// Decode decodes a media from a Reader
func Decode(r io.Reader) (*MP4, error) {
	v := &MP4{
		boxes: []Box{},
	}
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	for _, b := range l {
		switch b.Type() {
		case "ftyp":
			v.Ftyp = b.(*FtypBox)
		case "styp":
			v.Styp = b.(*StypBox)
		case "moov":
			v.Moov = b.(*MoovBox)
		case "mdat":
			v.Mdat = b.(*MdatBox)
			v.Mdat.ContentSize = uint32(b.Size() - BoxHeaderSize)
			break
		}
		if decoders[b.Type()] != nil {
			continue
		}
		u := b.(*UkwnBox)
		u.h.Type = b.Type()
		u.h.Size = uint32(b.Size())
		v.boxes = append(v.boxes, u)
	}
	return v, nil
}

// Dump displays some information about a media
func (m *MP4) Dump() {
	if m.Ftyp != nil {
		m.Ftyp.Dump()
	}
	if m.Styp != nil {
		m.Styp.Dump()
	}
	if m.Moov != nil {
		m.Moov.Dump()
	}
	for _, b := range m.boxes {
		if b != nil {
			ukwn := b.(*UkwnBox)
			ukwn.Dump()

		}
	}
}

// Boxes lists the top-level boxes from a media
func (m *MP4) Boxes() []Box {
	return m.boxes
}

// Encode encodes a media to a Writer
func (m *MP4) Encode(w io.Writer) error {
	err := m.Ftyp.Encode(w)
	if err != nil {
		return err
	}
	err = m.Moov.Encode(w)
	if err != nil {
		return err
	}
	for _, b := range m.boxes {
		if b.Type() != "ftyp" && b.Type() != "moov" {
			err = b.Encode(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
