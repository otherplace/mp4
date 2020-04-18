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
// Other boxes can also be present (pdin, moof, mfra, free, ...)
type MP4 struct {
	Ftyp *FtypBox   `json:"ftyp,omitempty"`
	Styp *StypBox   `json:"styp,omitempty"`
	Pdin *PdinBox   `json:"pdin,omitempty"`
	Moov *MoovBox   `json:"moov,omitempty"`
	Moof []*MoofBox `json:"moof,omitempty"`
	//Mfra  *MfraBox `json:"mfra,omitempty"`
	Mdat *MdatBox   `json:"mdat,omitempty"`
	Free []*FreeBox `json:"free,omitempty"`
	//Skip  []*SkipBox `json:"skip,omitempty"`
	Udta *UdtaBox   `json:"udta,omitempty"`
	Sidx []*SidxBox `json:"sidx,omitempty"`
	//Ssix  []*SSixBox `json:"ssix,omitempty"`
	//Prft  []*PrftBox `json:"prft,omitempty"`
	Mfra *MfraBox `json:"mfra,omitempty"`
	Meta *MetaBox `json:"meta,omitempty"`
	//Meco  *MetaBox `json:"meco,omitempty"`
	boxes []Box `json:",omitempty"`
}

type fMP4 struct {
	Styp *StypBox `json:"styp,omitempty"`
	//Pdin  *PdinBox `json:"pdin,omitempty"`
	Moof []*MoofBox `json:"moof,omitempty"`
	//Mfra  *MfraBox `json:"mfra,omitempty"`
	Mdat  *MdatBox   `json:"mdat,omitempty"`
	Free  []*FreeBox `json:"free,omitempty"`
	boxes []Box      `json:",omitempty"`
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
		case "mfra":
			v.Mfra = b.(*MfraBox)
		case "styp":
			v.Styp = b.(*StypBox)
		case "meta":
			v.Meta = b.(*MetaBox)
		case "moof":
			v.Moof = append(v.Moof, b.(*MoofBox))
		case "moov":
			v.Moov = b.(*MoovBox)
		case "mdat":
			v.Mdat = b.(*MdatBox)
		case "sidx":
			v.Sidx = append(v.Sidx, b.(*SidxBox))
		case "free":
			v.Free = append(v.Free, b.(*FreeBox))
		case "pdin":
			v.Pdin = b.(*PdinBox)
		//case "udta":
		//	v.Udta = b.(*UdtaBox)
		default:
			if decoders[b.Type()] != nil {
				v.boxes = append(v.boxes, b.Box())
			} else {
				v.boxes = append(v.boxes, b.(*UkwnBox))
			}
		}
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
	if m.Mfra != nil {
		m.Mfra.Dump()
	}
	if m.Moov != nil {
		m.Moov.Dump()
	}
	for _, f := range m.Moof {
		f.Dump()
	}
	if m.Meta != nil {
		m.Meta.Dump()
	}
	if m.Mdat != nil {
		m.Mdat.Dump()
	}
	for _, s := range m.Sidx {
		s.Dump()
	}
	for _, f := range m.Moof {
		f.Dump()
	}
	for _, b := range m.boxes {
		b.Dump()
	}
}

// Boxes lists the top-level boxes from a media
func (m *MP4) Boxes() []Box {
	return m.boxes
}

// Encode encodes a media to a Writer
func (m *MP4) Encode(w io.Writer) error {
	if m.Ftyp != nil {
		err := m.Ftyp.Encode(w)
		if err != nil {
			return err
		}
	}
	if m.Styp != nil {
		err := m.Styp.Encode(w)
		if err != nil {
			return err
		}
	}
	if m.Moov != nil {
		err := m.Moov.Encode(w)
		if err != nil {
			return err
		}
	}
	if m.Mdat != nil {
		err := m.Mdat.Encode(w)
		if err != nil {
			return err
		}
	}
	for _, f := range m.Free {
		err := f.Encode(w)
		if err != nil {
			return err
		}
	}
	for _, s := range m.Sidx {
		err := s.Encode(w)
		if err != nil {
			return err
		}
	}
	if m.Mfra != nil {
		err := m.Mfra.Encode(w)
		if err != nil {
			return err
		}
	}
	for _, b := range m.boxes {
		if b.Type() != "ftyp" && b.Type() != "moov" {
			err := b.Encode(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
