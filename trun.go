package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘trun’
// Container:	 Track	Fragment	Box	('traf')
// Mandatory:	No
// Quantity:	 Zero	or	more
const (
	DataOffsetPresent                   = 0x000001
	FirstSampleFlagsPresent             = 0x000004
	SampleDurationPresent               = 0x000100
	SampleSizePresent                   = 0x000200
	SampleFlagsPresent                  = 0x000400
	SampleCompositionTimeOffsetsPresent = 0x000800
)

type Sample struct {
	SampleDuration              uint32 `json:"SampleDuration,omitempty"`
	SampleSize                  uint32 `json:"SampleSize,omitempty"`
	SampleFlags                 uint32 `json:"SampleFlags,omitempty"`
	SampleCompositionTimeOffset uint32 `json:"SampleCompositionTimeOffset,omitempty"`
}

type TrunBox struct {
	Version          byte
	Flags            [3]byte
	SampleCount      uint32 `json:"SampleCount,"`
	DataOffset       int32  `json:"DataOffset,omitempty"`
	FirstSampleFlags uint32 `json:"FirstSampleFlags,omitempty"`
	Samples          []*Sample
}

func (b *TrunBox) Box() Box {
	return b
}

func (b *TrunBox) Type() string {
	return "trun"
}

func (b *TrunBox) Size() int {
	l := BoxHeaderSize + 12
	flag := BEUint28(b.Flags[:])
	for _, _ = range b.Samples {
		if compareFlag(flag, SampleDurationPresent) {
			l = l + 8
		}
		if compareFlag(flag, SampleSizePresent) {
			l = l + 8
		}
		if compareFlag(flag, SampleFlagsPresent) {
			l = l + 8
		}
		if compareFlag(flag, SampleCompositionTimeOffsetsPresent) {
			l = l + 8
		}
	}
	return l
}
func (b *TrunBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:8], b.SampleCount)

	startOffset := 8
	flag := BEUint28(b.Flags[:])
	if compareFlag(flag, DataOffsetPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], uint32(b.DataOffset))
		startOffset = startOffset + 8
	}
	if compareFlag(flag, FirstSampleFlagsPresent) {
		binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], b.FirstSampleFlags)
		startOffset = startOffset + 8
	}
	for _, s := range b.Samples {
		if compareFlag(flag, SampleDurationPresent) {
			binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], s.SampleDuration)
			startOffset = startOffset + 8
		}
		if compareFlag(flag, SampleSizePresent) {
			binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], s.SampleSize)
			startOffset = startOffset + 8
		}
		if compareFlag(flag, SampleFlagsPresent) {
			binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], s.SampleFlags)
			startOffset = startOffset + 8
		}
		if compareFlag(flag, SampleCompositionTimeOffsetsPresent) {
			binary.BigEndian.PutUint32(buf[startOffset:startOffset+4], s.SampleCompositionTimeOffset)
			startOffset = startOffset + 8
		}
	}
	return err
}
func (b *TrunBox) Dump() {
	fmt.Printf("Track fragment Run Box\n")
}

func DecodeTrun(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &TrunBox{
		Version:     data[0],
		Flags:       [3]byte{data[1], data[2], data[3]},
		SampleCount: binary.BigEndian.Uint32(data[4:8]),
	}
	startOffset := 8
	flag := BEUint28(b.Flags[:])
	if compareFlag(flag, DataOffsetPresent) {
		b.DataOffset = int32(binary.BigEndian.Uint32(data[startOffset : startOffset+4]))
		startOffset = startOffset + 4
	}
	if compareFlag(flag, FirstSampleFlagsPresent) {
		b.FirstSampleFlags = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
		startOffset = startOffset + 4
	}
	for i := 0; i < int(b.SampleCount); i++ {
		s := &Sample{}
		if compareFlag(flag, SampleDurationPresent) {
			s.SampleDuration = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
			startOffset = startOffset + 4
		}
		if compareFlag(flag, SampleSizePresent) {
			s.SampleSize = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
			startOffset = startOffset + 4
		}
		if compareFlag(flag, SampleFlagsPresent) {
			s.SampleFlags = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
			startOffset = startOffset + 4
		}
		if compareFlag(flag, SampleCompositionTimeOffsetsPresent) {
			s.SampleCompositionTimeOffset = binary.BigEndian.Uint32(data[startOffset : startOffset+4])
			startOffset = startOffset + 4
		}
		b.Samples = append(b.Samples, s)
	}
	return b, nil
}
