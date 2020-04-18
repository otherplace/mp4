package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘mfro’
// Container:	 Movie	Fragment	Random	Access	Box	(‘mfra’)
// Mandatory:	Yes
// Quantity:	 Exactly	one
type MfroBox struct {
	Version byte
	Flags   []byte
	MSize   uint32 `json:"Size,"`
}

func DecodeMfro(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &MfroBox{
		Version: data[0],
		Flags:   []byte{data[1], data[2], data[3]},
		MSize:   binary.BigEndian.Uint32(data[4:8]),
	}, nil
}

func (b *MfroBox) Box() Box {
	return b
}

func (b *MfroBox) Type() string {
	return "mfro"
}

func (b *MfroBox) Size() int {
	return int(BoxHeaderSize + 4 + 4)
}

func (b *MfroBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.MSize)
	_, err = w.Write(buf)
	return err
}

func (b *MfroBox) Dump() {
	fmt.Printf("Binary XML Box\n")
	fmt.Printf("+- Version: %d\n", b.Version)
	fmt.Printf("+- Flag: %v\n", b.Flags)
	fmt.Printf("+- Size: %d\n", b.MSize)
}
