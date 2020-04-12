package mp4

import (
	"fmt"
	"io"
)

type MvexBox struct {
	Trex *TrexBox
	//Mehd *MehdBox
}

func DecodeMvex(r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	m := &MvexBox{}
	for _, b := range l {
		switch b.Type() {
		case "trex":
			m.Trex = b.(*TrexBox)
		case "mehd":
			//m.Mehd = b.(*MehdBox)
		}
	}
	return m, err
}

func (b *MvexBox) Type() string {
	return "mvex"
}

func (b *MvexBox) Size() int {
	l := BoxHeaderSize + b.Trex.Size()
	//if b.Mehd != nil {
	//	l = l + b.Mehd.Size()
	//}
	return l
}

func (b *MvexBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	err = b.Trex.Encode(w)
	if err != nil {
		return err
	}
	return err
}

func (b *MvexBox) Dump() {
	fmt.Printf("Movie Extends Box\n")
	b.Trex.Dump()
}
