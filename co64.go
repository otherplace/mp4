package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box Type: ‘co64’
// Container: Sample Table Box (‘stbl’)
// Mandatory: Yes
// Quantity: Exactly one variant must be present

type Co64Box struct {
	FullBox
	EntryCount  uint32
	ChunkOffset []uint64
}

func DecodeCo64(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &Co64Box{
		FullBox:     DecodeFullBox(data[0:4]),
		EntryCount:  binary.BigEndian.Uint32(data[4:8]),
		ChunkOffset: []uint64{},
	}
	for i := 0; i < int(b.EntryCount); i++ {
		chunk := binary.BigEndian.Uint64(data[(8 + 8*i):(16 + 8*i)])
		b.ChunkOffset = append(b.ChunkOffset, chunk)
	}
	return b, nil
}

func (b *Co64Box) Box() Box {
	return b
}

func (b *Co64Box) Type() string {
	return "co64"
}

func (b *Co64Box) Size() int {
	return BoxHeaderSize + 8 + len(b.ChunkOffset)*4
}

func (b *Co64Box) Dump() {
	fmt.Println("Chunk byte offsets:")
	for i, o := range b.ChunkOffset {
		fmt.Printf(" #%d : starts at %d\n", i, o)
	}
}

func (b *Co64Box) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	EncodeFullBox(b.FullBox, buf)
	binary.BigEndian.PutUint32(buf[4:], b.EntryCount)
	for i := range b.ChunkOffset {
		binary.BigEndian.PutUint64(buf[12+8*i:], b.ChunkOffset[i])
	}
	_, err = w.Write(buf)
	return err
}
