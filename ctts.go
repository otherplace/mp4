package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Composition Time to Sample Box (ctts - optional)
//
// Contained in: Sample Table Box (stbl)
//
// Status: version 0 decoded. version 1 uses int32 for offsets
type CttsBox struct {
	Version      byte
	Flags        [3]byte
	EntryCount   uint32
	SampleCount  []uint32
	SampleOffset []uint32 // int32 for version 1
}

func DecodeCtts(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &CttsBox{
		Version:      data[0],
		Flags:        [3]byte{data[1], data[2], data[3]},
		SampleCount:  []uint32{},
		SampleOffset: []uint32{},
	}
	ec := binary.BigEndian.Uint32(data[4:8])
	b.EntryCount = ec
	for i := 0; i < int(ec); i++ {
		s_count := binary.BigEndian.Uint32(data[(8 + 8*i):(12 + 8*i)])
		s_offset := binary.BigEndian.Uint32(data[(12 + 8*i):(16 + 8*i)])
		b.SampleCount = append(b.SampleCount, s_count)
		b.SampleOffset = append(b.SampleOffset, s_offset)
	}
	return b, nil
}

func (b *CttsBox) Box() Box {
	return b
}
func (b *CttsBox) Type() string {
	return "ctts"
}

func (b *CttsBox) Size() int {
	return BoxHeaderSize + 8 + len(b.SampleCount)*8
}

func (b *CttsBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.EntryCount)
	for i := range b.SampleCount {
		binary.BigEndian.PutUint32(buf[8+8*i:], b.SampleCount[i])
		binary.BigEndian.PutUint32(buf[12+8*i:], b.SampleOffset[i])
	}
	_, err = w.Write(buf)
	return err
}

func (b *CttsBox) Dump() {
	fmt.Printf("Composition Time to Sample Box\n")
	fmt.Printf("+- Entry Count: %d\n", b.EntryCount)
	for i := 0; i < int(b.EntryCount); i++ {
		fmt.Printf(" +- #%d Sample Count: %d\n", i, b.SampleCount[i])
		fmt.Printf(" +- #%d Sample Offset: %d\n", i, b.SampleOffset[i])
	}
}
