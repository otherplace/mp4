package mp4

import (
	"fmt"
	"io"
)

// File Type Box (ftyp - mandatory)
//
// Status: decoded
type FreeBox struct {
	notDecoded []byte
}

func DecodeFree(h BoxHeader, r io.Reader) (Box, error) {
	if h.Size <= BoxHeaderSize {
		return &FreeBox{nil}, nil
	}
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &FreeBox{data}, nil
}

func (b *FreeBox) Box() Box {
	return b
}
func (b *FreeBox) Type() string {
	return "free"
}

func (b *FreeBox) Size() int {
	return BoxHeaderSize + len(b.notDecoded)
}

func (b *FreeBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	_, err = w.Write(b.notDecoded)
	return err
}

func (b *FreeBox) Dump() {
	fmt.Printf("Free Space Box\n")
}
