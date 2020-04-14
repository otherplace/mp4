package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

type TrexBox struct {
	Version                byte
	Flags                  [3]byte
	TrackId                uint32
	SampleDescriptionIndex uint32
	SampleDuration         uint32
	SampleSize             uint32
	SampleFlags            uint32
}

func DecodeTrex(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &TrexBox{
		Version:                data[0],
		Flags:                  [3]byte{data[1], data[2], data[3]},
		TrackId:                binary.BigEndian.Uint32(data[4:8]),
		SampleDescriptionIndex: binary.BigEndian.Uint32(data[8:12]),
		SampleDuration:         binary.BigEndian.Uint32(data[12:16]),
		SampleSize:             binary.BigEndian.Uint32(data[16:20]),
		SampleFlags:            binary.BigEndian.Uint32(data[20:24]),
	}, nil
}

func (b *TrexBox) Box() Box {
	return b
}

func (b *TrexBox) Type() string {
	return "trex"
}

func (b *TrexBox) Size() int {
	return BoxHeaderSize + 24
}

func (b *TrexBox) Dump() {
	fmt.Printf("TrackExtendsBox:\n")
	fmt.Printf(" Version: %d\n", b.Version)
	fmt.Printf(" Flags: 0x%x\n", b.Flags)
	fmt.Printf(" TrackId: %d\n", b.TrackId)
	fmt.Printf(" Sample Description index: %d\n", b.SampleDescriptionIndex)
	fmt.Printf(" Sample Duration: %d\n", b.SampleDuration)
	fmt.Printf(" Sample Size: %d\n", b.SampleSize)
	fmt.Printf(" Sample Flags: 0x%x\n", b.SampleFlags)
}

func (b *TrexBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.TrackId)
	binary.BigEndian.PutUint32(buf[8:], b.SampleDescriptionIndex)
	binary.BigEndian.PutUint32(buf[12:], b.SampleDuration)
	binary.BigEndian.PutUint32(buf[16:], b.SampleSize)
	binary.BigEndian.PutUint32(buf[20:], b.SampleFlags)
	_, err = w.Write(buf)
	return err
}
