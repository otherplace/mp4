package mp4

import (
	"fmt"
	"io"
)

// Box	Type:	 ‘meta’
// Container:	 File,	Movie	Box	(‘moov’),	Track	Box	(‘trak’),
//	 Additional	Metadata	Container	Box	(‘meco’),
//	 Movie	Fragment	Box	(‘moof’)	or	Track	Fragment	Box	(‘traf’)
// Mandatory:	No
// Quantity:	 Zero	or	one	(in	File,	‘moov’,	and	‘trak’),	One	or	more	(in	‘meco’)
type MetaBox struct {
	Hdlr *HdlrBox `json:"hdlr,"`
	Bxml *BxmlBox `json:"bxml,omitempty"`
	Dinf *DinfBox `json:"dinf,omitempty"`
	Iloc *IlocBox `json:"iloc,omitempty"`
	//Pitm *PitmBox `json:"pitm,omitempty"`
	//Ipro *IproBox `json:"ipro,omitempty
	//Iinf *IinfBox `json:"iinf,omitempty"`
	Boxes []Box `json:",omitempty"`
}

func DecodeMeta(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MetaBox{}
	for _, b := range l {
		switch b.Type() {
		case "hdlr":
			m.Hdlr = b.(*HdlrBox)
		case "dinf":
			m.Dinf = b.(*DinfBox)
		case "bxml":
			m.Bxml = b.(*BxmlBox)
		case "iloc":
			m.Iloc = b.(*IlocBox)
		default:
			m.Boxes = append(m.Boxes, b.Box())
		}
	}
	return m, err
}

func (b *MetaBox) Box() Box {
	return b
}

func (b *MetaBox) Type() string {
	return "meta"
}

func (b *MetaBox) Size() int {
	l := BoxHeaderSize
	if b.Hdlr != nil {
		l += b.Hdlr.Size()
	}
	if b.Dinf != nil {
		l += b.Dinf.Size()
	}
	return l
}

func (b *MetaBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	if b.Hdlr != nil {
		err = b.Encode(w)
		if err != nil {
			return err
		}
	}
	if b.Dinf != nil {
		err = b.Encode(w)
		if err != nil {
			return err
		}
	}
	return err
}

func (b *MetaBox) Dump() {
	fmt.Printf("Meta Box\n")
	if b.Hdlr != nil {
		b.Hdlr.Dump()
	}
	if b.Dinf != nil {
		b.Dinf.Dump()
	}
}
