package mp4

import (
	"io"
)

// Object Descriptor Container Box (iods - optional)
//
// Contained in : Movie Box (‘moov’)
//
// Status: not decoded
type IodsBox struct {
	notDecoded []byte
}

func DecodeIods(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &IodsBox{
		notDecoded: data,
	}, nil
}

func (b *IodsBox) Box() Box {
	return b
}
func (b *IodsBox) Type() string {
	return "iods"
}

func (b *IodsBox) Size() int {
	return BoxHeaderSize + len(b.notDecoded)
}

func (b *IodsBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	_, err = w.Write(b.notDecoded)
	return err
}
