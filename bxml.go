package mp4

import (
	"fmt"
	"io"
)

// Box	Type:	 ‘xml ‘	or	‘bxml’
// Container:	 Meta	box	(‘meta’)
// Mandatory:	No
// Quantity:	 Zero	or	one
type BxmlBox struct {
	Version byte
	Flags   [3]byte
	Data    []byte
}

func DecodeBxml(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &BxmlBox{
		Version: data[0],
		Flags:   [3]byte{data[1], data[2], data[3]},
		Data:    data[4:],
	}, nil
}

func (b *BxmlBox) Box() Box {
	return b
}

func (b *BxmlBox) Type() string {
	return "bxml"
}

func (b *BxmlBox) Size() int {
	return int(BoxHeaderSize + 4 + len(b.Data))
}

func (b *BxmlBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	copy(buf[4:], b.Data[:])
	_, err = w.Write(buf)
	return err
}

func (b *BxmlBox) Dump() {
	fmt.Printf("Binary XML Box\n")
	fmt.Printf("+- Version: %d\n", b.Version)
	fmt.Printf("+- Flag: %v\n", b.Flags)
	fmt.Printf("+- Data: %s\n", string(b.Data))
}
