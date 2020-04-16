package mp4

import (
	"fmt"
	"io"
)

type UkwnBox struct {
	Header      BoxHeader
	PayloadSize uint32
	Data        []byte
}

func DecodeUkwnBox(h BoxHeader, r io.Reader) (Box, error) {
	if lr, limited := r.(*io.LimitedReader); limited {
		r = lr.R
	}
	data := make([]byte, h.Size-BoxHeaderSize)
	n, _ := r.Read(data)

	b := &UkwnBox{
		Header: BoxHeader{
			Size: h.Size,
			Type: h.Type,
		},
		Data:        data,
		PayloadSize: uint32(n),
	}
	return b, nil
}

func (b *UkwnBox) Type() string {
	return b.Header.Type
}

func (b *UkwnBox) Size() int {
	return int(b.PayloadSize)
}

func (b *UkwnBox) Dump() {
	fmt.Printf("Box type: %s\n", b.Type())
	fmt.Printf(" Data length: %d\n", len(b.Data))
}
func (b *UkwnBox) Encode(w io.Writer) error {
	return nil
}

func (b *UkwnBox) Box() Box {
	return b
}
