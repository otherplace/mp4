package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘sbgp’
// Container:	 Sample	Table	Box	(‘stbl’)	or	Track	Fragment	Box	(‘traf’)
// Mandatory:	No
// Quantity:	 Zero	or	more.

type Entry struct {
	SampleCount           uint32
	GroupDescriptionIndex uint32
}

type SbgpBox struct {
	Version      byte
	Flags        [3]byte
	GroupingType uint32
	EntryCount   uint32
	Entries      []*Entry
}

func DecodeSbgp(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &SbgpBox{
		Version:      data[0],
		Flags:        [3]byte{data[1], data[2], data[3]},
		GroupingType: binary.BigEndian.Uint32(data[4:8]),
		EntryCount:   binary.BigEndian.Uint32(data[8:12]),
	}
	offset := 12
	for i := 0; i < int(b.EntryCount); i++ {
		e := &Entry{
			SampleCount:           binary.BigEndian.Uint32(data[offset+i*4 : offset+4+i*4]),
			GroupDescriptionIndex: binary.BigEndian.Uint32(data[offset+4+i*4 : offset+8+i*4]),
		}
		offset += 8
		b.Entries = append(b.Entries, e)
	}
	return b, nil
}

func (b *SbgpBox) Box() Box {
	return b
}

func (b *SbgpBox) Type() string {
	return "sbgp"
}

func (b *SbgpBox) Size() int {
	return int(BoxHeaderSize + 4 + 8 + (b.EntryCount * 8))
}

func (b *SbgpBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:8], b.GroupingType)
	binary.BigEndian.PutUint32(buf[8:12], b.EntryCount)
	offset := 12
	for _, e := range b.Entries {
		binary.BigEndian.PutUint32(buf[offset:], e.SampleCount)
		offset += 4
		binary.BigEndian.PutUint32(buf[offset:], e.GroupDescriptionIndex)
		offset += 4
	}
	_, err = w.Write(buf)
	return err
}

func (b *SbgpBox) Dump() {
	fmt.Printf("Sample To Group Box\n")
	fmt.Printf("+- Version: %d\n", b.Version)
	fmt.Printf("+- Flag: %v\n", b.Flags)
}
