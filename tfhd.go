package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘tfhd’
// Container:	 Track	Fragment	Box	('traf')
// Mandatory:	Yes
// Quantity:	 Exactly	one
type TfhdBox struct {
	Version                byte
	Flags                  [3]byte
	TrackId                uint32 `json:"TrackId,"` // ISO IEC 14496-12 said this is 64 bits, but actual media samples got 32 bits
	BaseDataOffset         uint32 `json:"BaseDataOffset,omitempty"`
	SampleDescriptionIndex uint32 `json:"SampleDescriptionIndex,omitempty"`
	DefaultSampleDuration  uint32 `json:"DefaultSampleDuration,omitempty"`
	DefaultSampleSize      uint32 `json:"DefaultSampleSize,omitempty"`
	DefaultSampleFlags     uint32 `json:"DefaultSampleFlags,omitempty"`
	DurationIsEmpty        bool   `json:"DurationIsEmpty,omitempty"`
	DefaultBaseIsMoof      bool   `json:"DefaultBaseIsMoof,omitempty"`
}

const (
	BaseDataOffsetPresent         = 0x000001
	SampleDescriptionIndexPresent = 0x000002
	DefaultSampleDurationPresent  = 0x000008
	DefaultSampleSizePresent      = 0x000010
	DefaultSampleFlagsPresent     = 0x000020
	DurationIsEmpty               = 0x010000
	DefaultBaseIsMoof             = 0x020000
)

func (b *TfhdBox) Type() string {
	return "tfhd"
}

func (b *TfhdBox) Size() int {
	l := BoxHeaderSize + 8
	flag := BEUint28(b.Flags[:])
	if compareFlag(flag, BaseDataOffsetPresent) {
		l = l + 8
	}
	if compareFlag(flag, SampleDescriptionIndexPresent) {
		l = l + 8
	}
	if compareFlag(flag, DefaultSampleDurationPresent) {
		l = l + 8
	}
	if compareFlag(flag, DefaultSampleSizePresent) {
		l = l + 8
	}
	if compareFlag(flag, DefaultSampleFlagsPresent) {
		l = l + 8
	}

	return l
}
func (b *TfhdBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:8], b.TrackId)

	startOffset := 8
	flag := BEUint28(b.Flags[:])
	if compareFlag(flag, BaseDataOffsetPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.BaseDataOffset)
		startOffset = startOffset + 4
	}
	if compareFlag(flag, SampleDescriptionIndexPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.SampleDescriptionIndex)
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleDurationPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.DefaultSampleDuration)
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleSizePresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.DefaultSampleSize)
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleFlagsPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.DefaultSampleFlags)
		startOffset = startOffset + 4
	}

	return err
}
func (b *TfhdBox) Dump() {
	fmt.Printf("Track fragment Header Box\n")
}

func DecodeTfhd(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &TfhdBox{
		Version: data[0],
		Flags:   [3]byte{data[1], data[2], data[3]},
		TrackId: binary.BigEndian.Uint32(data[4:8]),
	}
	startOffset := 8
	flag := BEUint28(b.Flags[:])
	if compareFlag(flag, BaseDataOffsetPresent) {
		b.BaseDataOffset = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	if compareFlag(flag, SampleDescriptionIndexPresent) {
		b.SampleDescriptionIndex = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleDurationPresent) {
		b.DefaultSampleDuration = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleSizePresent) {
		b.DefaultSampleSize = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DefaultSampleFlagsPresent) {
		b.DefaultSampleFlags = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	if compareFlag(flag, DurationIsEmpty) {
		b.DurationIsEmpty = true
	}
	if compareFlag(flag, DefaultBaseIsMoof) {
		b.DefaultBaseIsMoof = true
	}
	return b, nil
}
