package mp4

import (
	"fmt"
	"io"
)

// Media Data Box (mdat - optional)
//
// Status: not decoded
//
// The mdat box contains media chunks/samples.
//
// It is not read, only the io.Reader is stored, and will be used to Encode (io.Copy) the box to a io.Writer.
type MdatBox struct {
	ContentSize uint32
	data        []byte
	r           io.Reader
}

func DecodeMdat(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	n, _ := r.Read(data)
	// r is a LimitedReader
	if lr, limited := r.(*io.LimitedReader); limited {
		r = lr.R
	}
	return &MdatBox{
		data:        data,
		r:           r,
		ContentSize: uint32(n),
	}, nil
}

func (b *MdatBox) Box() Box {
	return b
}

func (b *MdatBox) Type() string {
	return "mdat"
}

func (b *MdatBox) Size() int {
	return BoxHeaderSize + int(b.ContentSize)
}

func (b *MdatBox) Reader() io.Reader {
	return b.r
}

func (b *MdatBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	_, err = w.Write(b.data)
	return err
}

func (b *MdatBox) Dump() {
	fmt.Printf("Media Data Box\n")
	fmt.Printf("+- ContentSize: %d\n", b.ContentSize)
}
