package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘mfhd’
// Container:	 Movie	Fragment	Box	('moof')
// Mandatory:	Yes
// Quantity:	 Exactly	one
type MfhdBox struct {
	Version        byte
	Flags          [3]byte
	SequenceNumber uint32
}

func (b *MfhdBox) Box() Box {
	return b
}

func (b *MfhdBox) Type() string {
	return "mfhd"
}

func (b *MfhdBox) Size() int {
	return BoxHeaderSize + 8
}
func (b *MfhdBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.SequenceNumber)

	return err
}
func (b *MfhdBox) Dump() {
	fmt.Printf("Movie Fragment Header Box\n")
	fmt.Printf(" Sequence Number: %d\n", b.SequenceNumber)
}

func DecodeMfhd(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &MfhdBox{
		Version:        data[0],
		Flags:          [3]byte{data[1], data[2], data[3]},
		SequenceNumber: binary.BigEndian.Uint32(data[4:8]),
	}, nil
}
