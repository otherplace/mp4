package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Types:	 ‘pdin’
// Container:	 File
// Mandatory:	 No
// Quantity:	 Zero	or	One
type InfoBox struct {
	Rate         uint32
	InitialDelay uint32
}
type PdinBox struct {
	Version   byte
	Flags     []byte
	InfoBoxes []*InfoBox
}

func (b *PdinBox) Box() Box {
	return b
}

func (b *PdinBox) Type() string {
	return "pdin"
}

func (b *PdinBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]

	for _, p := range b.InfoBoxes {
		binary.BigEndian.PutUint32(buf[4:], p.Rate)
		binary.BigEndian.PutUint32(buf[8:], p.InitialDelay)
	}
	return err
}

func (b *PdinBox) Size() int {
	l := BoxHeaderSize
	for _, _ = range b.InfoBoxes {
		l = l + 8 + 8
	}
	return l
}

func DecodePdin(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &PdinBox{
		Version: data[0],
		Flags:   []byte{data[1], data[2], data[3]},
	}
	for i := 4; 4 < h.Size-BoxHeaderSize; i += 4 {
		ib := &InfoBox{
			Rate:         binary.BigEndian.Uint32(data[i : i+4]),
			InitialDelay: binary.BigEndian.Uint32(data[i+4 : i+8]),
		}
		b.InfoBoxes = append(b.InfoBoxes, ib)
	}
	return b, nil
}

func (b *PdinBox) Dump() {
	fmt.Printf("Progressive Download Information Box\n")
	fmt.Printf("+- Version: %d\n", b.Version)
	fmt.Printf("+- Flags: %v\n", b.Flags)
	for _, i := range b.InfoBoxes {
		fmt.Printf(" +- Rate: %d\n", i.Rate)
		fmt.Printf(" +- InitialDelay: %d\n", i.InitialDelay)
	}
}
