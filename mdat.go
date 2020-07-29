package mp4

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Media Data Box (mdat - optional)
//
// Status: not decoded
//
// The mdat box contains media chunks/samples.
//
// It is not read, only the io.Reader is stored, and will be used to Encode (io.Copy) the box to a io.Writer.
type MdatBox struct {
	ContentSize uint64
	r           io.Reader
}

func DecodeMdat(h BoxHeader, r io.Reader) (Box, error) {
	dataSize := uint64(h.Size - BoxHeaderSize)
	if h.Size == 1 {
		dataSize = h.LargeSize - uint64(BoxHeaderSize*2)
	}
	// r is a LimitedReader
	if lr, limited := r.(*io.LimitedReader); limited {
		r = lr.R
	}
	// FIXME:
	io.CopyN(ioutil.Discard, r, int64(dataSize))
	return &MdatBox{
		r:           r,
		ContentSize: dataSize,
	}, nil
}

func (b *MdatBox) Box() Box {
	return b
}

func (b *MdatBox) Type() string {
	return "mdat"
}

func (b *MdatBox) Size() int {
	// FIXME:
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
	_, err = io.Copy(w, b.r)
	return err
}

func (b *MdatBox) Dump() {
	fmt.Printf("Media Data Box\n")
	fmt.Printf("+- ContentSize: %d\n", b.ContentSize)
}
