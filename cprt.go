package mp4

import (
	"io"
)

// Box	Type:	 ‘cprt’
// Container:	 User	data	box	(‘udta’)
// Mandatory:	 No
// Quantity:	 Zero	or	more
type CprtBox struct {
	Version    byte
	Flags      [3]byte
	notDecoded []byte
}

// TODO
func DecodeCprt(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &CprtBox{
		Version:    data[0],
		Flags:      [3]byte{data[1], data[2], data[3]},
		notDecoded: data[4:],
	}, nil
}

func (b *CprtBox) Box() Box {
	return b
}

func (b *CprtBox) Type() string {
	return "meta"
}

func (b *CprtBox) Size() int {
	return BoxHeaderSize + 4 + len(b.notDecoded)
}

func (b *CprtBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	copy(buf[4:], b.notDecoded)
	_, err = w.Write(buf)
	return err
}
