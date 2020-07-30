package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box Type: ‘sgpd’
// Container: Sample Table Box (‘stbl’) or Track Fragment Box (‘traf’)
// Mandatory: No
// Quantity: Zero or more, with one for each Sample to Group Box.

type SgpdBox struct {
	FullBox
	GroupType                     uint32
	DefaultLength                 uint32 // Version == 1
	DefaultSampleDescriptionIndex uint32 // Version >= 2
	EntryCount                    uint32
	DescriptionLength             []uint32 // Version == 1 && DefaultLength == 0
	SampleGroupEntry              []uint32
}

func DecodeSgpd(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &SgpdBox{
		FullBox:           DecodeFullBox(data[0:4]),
		GroupType:         binary.BigEndian.Uint32(data[4:8]),
		EntryCount:        binary.BigEndian.Uint32(data[12:16]),
		DescriptionLength: []uint32{},
		SampleGroupEntry:  []uint32{},
	}
	if b.Version == 1 {
		b.DefaultLength = binary.BigEndian.Uint32(data[8:12])
	} else if b.Version >= 2 {
		b.DefaultSampleDescriptionIndex = binary.BigEndian.Uint32(data[8:12])
	}
	for i := 0; i < int(b.EntryCount); i++ {
		var sampleGroupEntry uint32
		if b.Version == 1 && b.DefaultLength == 0 {
			descriptionLength := binary.BigEndian.Uint32(data[(12 + i*4):(16 + i*4)])
			b.DescriptionLength = append(b.DescriptionLength, descriptionLength)
			sampleGroupEntry = binary.BigEndian.Uint32(data[(16 + i*4):(20 + i*4)])
		} else {
			sampleGroupEntry = binary.BigEndian.Uint32(data[(12 + i*4):(16 + i*4)])
		}
		b.SampleGroupEntry = append(b.SampleGroupEntry, sampleGroupEntry)
	}
	return b, nil
}

func (b *SgpdBox) Box() Box {
	return b
}

func (b *SgpdBox) Type() string {
	return "sgpd"
}

func (b *SgpdBox) Size() int {
	sz := BoxHeaderSize + b.FullBox.Size()
	if b.Version == 1 && b.DefaultLength == 0 {
		sz += len(b.DescriptionLength)
	}
	sz += len(b.SampleGroupEntry)
	return sz
}

func (b *SgpdBox) Dump() {
	fmt.Println("Sample Group Description Box:")
	// TODO:
}

func (b *SgpdBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	EncodeFullBox(b.FullBox, buf)
	binary.BigEndian.PutUint32(buf[4:], b.GroupType)
	if b.Version == 1 {
		binary.BigEndian.PutUint32(buf[8:], b.DefaultLength)
	} else if b.Version >= 2 {
		binary.BigEndian.PutUint32(buf[8:], b.DefaultSampleDescriptionIndex)
	}
	binary.BigEndian.PutUint32(buf[12:], b.EntryCount)
	for i := uint32(0); i < b.EntryCount; i++ {
		if b.Version == 1 && b.DefaultLength == 0 {
			binary.BigEndian.PutUint32(buf[(12+i*4):], b.DescriptionLength[i])
			binary.BigEndian.PutUint32(buf[(16+i*4):], b.SampleGroupEntry[i])
		} else {
			binary.BigEndian.PutUint32(buf[(14+i*4):], b.SampleGroupEntry[i])
		}
	}
	_, err = w.Write(buf)
	return err
}
