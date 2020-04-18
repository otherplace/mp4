package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

type TfraBox struct {
	Version               byte
	Flags                 []byte
	TrackId               uint32
	Reserved              uint32 // 26 bits
	LengthSizeOfTrafNum   uint8
	LengthSizeOfTrunNum   uint8
	LengthSizeOfSampleNum uint8
	NumberOfTfraEntry     uint32
	Entries               []*TfraEntry
}

type TfraEntry struct {
	Time         uint32
	MoofOffset   uint32
	TrafNumber   uint32
	TrunNumber   uint32
	SampleNumber uint32
}

func DecodeTfra(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &TfraBox{
		Version:  data[0],
		Flags:    []byte{data[1], data[2], data[3]},
		TrackId:  binary.BigEndian.Uint32(data[4:8]),
		Reserved: 0,
	}
	lengths := binary.BigEndian.Uint32(data[8:12])
	b.LengthSizeOfTrafNum = (0x30 ^ uint8(lengths)) >> 4
	b.LengthSizeOfTrunNum = (0x0C ^ uint8(lengths)) >> 2
	b.LengthSizeOfSampleNum = (0x03 ^ uint8(lengths))
	b.NumberOfTfraEntry = binary.BigEndian.Uint32(data[12:16])
	offset := 16
	for i := 0; i < int(b.NumberOfTfraEntry); i++ {
		e := &TfraEntry{
			Time:         binary.BigEndian.Uint32(data[offset : offset+4]),
			MoofOffset:   binary.BigEndian.Uint32(data[offset+4 : offset+8]),
			TrafNumber:   binary.BigEndian.Uint32(data[offset+8 : offset+12]),
			TrunNumber:   binary.BigEndian.Uint32(data[offset+12 : offset+16]),
			SampleNumber: binary.BigEndian.Uint32(data[offset+20 : offset+24]),
		}
		offset += 24
		b.Entries = append(b.Entries, e)
	}
	return b, nil
}

func (b *TfraBox) Box() Box {
	return b
}

func (b *TfraBox) Type() string {
	return "tfra"
}

func (b *TfraBox) Size() int {
	l := BoxHeaderSize + 32
	for _, _ = range b.Entries {
		l = l + 20
	}

	return l
}

func (b *TfraBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	// TODO:
	_, err = w.Write(buf)
	return err
}

func (b *TfraBox) Dump() {
	fmt.Printf("Track Fragment Random Access Box\n")
	fmt.Printf("+- TrackId: %d\n", b.TrackId)
	fmt.Printf("+- NumberOfTfraEntry: %d\n", b.NumberOfTfraEntry)
}
