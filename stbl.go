package mp4

import (
	"fmt"
	"io"
)

// Sample Table Box (stbl - mandatory)
//
// Contained in : Media Information Box (minf)
//
// Status: partially decoded (anything other than stsd, stts, stsc, stss, stsz, stco, ctts is ignored)
//
// The table contains all information relevant to data samples (times, chunks, sizes, ...)
type StblBox struct {
	Sbgp  []*SbgpBox `json:"sbgp,omitempty"`
	Stsd  *StsdBox   `json:"stsd,"`
	Stts  *SttsBox   `json:"sttd,"`
	Stss  *StssBox   `json:"stss,omitempty"`
	Stsc  *StscBox   `json:"stsc,"`
	Stsz  *StszBox   `json:"stsz,"`
	Stco  *StcoBox   `json:"stco,omitempty"`
	Co64  *Co64Box   `json:"co64,omitempty"`
	Ctts  *CttsBox   `json:"ctts,omitempty"`
	Sgpd  *SgpdBox   `json:"jspd,omitempty"`
	Boxes []Box      `json:",omitempty"`
}

func DecodeStbl(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	s := &StblBox{}
	for _, b := range l {
		switch b.Type() {
		case "sbgp":
			s.Sbgp = append(s.Sbgp, b.(*SbgpBox))
		case "sgpd":
			s.Sgpd = b.(*SgpdBox)
		case "stsd":
			s.Stsd = b.(*StsdBox)
		case "stts":
			s.Stts = b.(*SttsBox)
		case "stsc":
			s.Stsc = b.(*StscBox)
		case "stss":
			s.Stss = b.(*StssBox)
		case "stsz":
			s.Stsz = b.(*StszBox)
		case "stco":
			s.Stco = b.(*StcoBox)
		case "co64":
			s.Co64 = b.(*Co64Box)
		case "ctts":
			s.Ctts = b.(*CttsBox)
		default:
			s.Boxes = append(s.Boxes, b.Box())
		}
	}
	return s, nil
}

func (b *StblBox) Box() Box {
	return b
}

func (b *StblBox) Type() string {
	return "stbl"
}

func (b *StblBox) Size() int {
	sz := b.Stsd.Size()
	for _, s := range b.Sbgp {
		sz += s.Size()
	}
	if b.Stts != nil {
		sz += b.Stts.Size()
	}
	if b.Stss != nil {
		sz += b.Stss.Size()
	}
	if b.Stsc != nil {
		sz += b.Stsc.Size()
	}
	if b.Stsz != nil {
		sz += b.Stsz.Size()
	}
	if b.Stco != nil {
		sz += b.Stco.Size()
	}
	if b.Ctts != nil {
		sz += b.Ctts.Size()
	}
	return sz + BoxHeaderSize
}

func (b *StblBox) Dump() {
	fmt.Printf("Sample Table Box\n")
	for _, s := range b.Sbgp {
		s.Dump()
	}
	if b.Stsc != nil {
		b.Stsc.Dump()
	}
	if b.Stts != nil {
		b.Stts.Dump()
	}
	if b.Stsz != nil {
		b.Stsz.Dump()
	}
	if b.Stss != nil {
		b.Stss.Dump()
	}
	if b.Stco != nil {
		b.Stco.Dump()
	}
}

func (b *StblBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	for _, s := range b.Sbgp {
		err = s.Encode(w)
		if err != nil {
			return err
		}
	}
	err = b.Stsd.Encode(w)
	if err != nil {
		return err
	}
	err = b.Stts.Encode(w)
	if err != nil {
		return err
	}
	if b.Stss != nil {
		err = b.Stss.Encode(w)
		if err != nil {
			return err
		}
	}
	err = b.Stsc.Encode(w)
	if err != nil {
		return err
	}
	err = b.Stsz.Encode(w)
	if err != nil {
		return err
	}
	err = b.Stco.Encode(w)
	if err != nil {
		return err
	}
	if b.Ctts != nil {
		return b.Ctts.Encode(w)
	}
	return nil
}
