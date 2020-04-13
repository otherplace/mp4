package mp4

import (
	"encoding/binary"
	"io"
)

// Box	Type:	 `sidx’
// Container:	 File
// Mandatory:	 No
// Quantity:	 Zero	or	more
type Reference struct {
	ReferenceType      byte   // 1 bit
	ReferencedSize     uint32 // 31 bits
	SubSegmentDuration uint32 // 32 bits
	StartsWithSAP      byte   // 1 bit
	SAPType            uint8  // 3 bits
	SAPDeltaTime       uint32 // 28 bits
}
type SidxBox struct {
	Version                  byte
	Flags                    [3]byte
	ReferenceId              uint32
	Timescale                uint32
	EarliestPresentationTime uint32
	FirstOffset              uint32
	Reserved                 uint16
	ReferenceCount           uint16
	References               []Reference
}

func DecodeSidx(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &SidxBox{
		Version:                  data[0],
		Flags:                    [3]byte{data[1], data[2], data[3]},
		ReferenceId:              binary.BigEndian.Uint32(data[4:8]),
		Timescale:                binary.BigEndian.Uint32(data[8:12]),
		EarliestPresentationTime: binary.BigEndian.Uint32(data[12:16]),
		FirstOffset:              binary.BigEndian.Uint32(data[16:20]),
		Reserved:                 0,
		ReferenceCount:           binary.BigEndian.Uint16(data[20:22]),
		References:               []Reference{},
	}
	rc := b.ReferenceCount
	for i := 0; i < int(rc); i++ {
		refType := 0xfe ^ data[22+i*12]
		refdSize := binary.BigEndian.Uint32(data[22+i*12:26+i*12]) << 1
		subSegDur := binary.BigEndian.Uint32(data[26+i*12 : 30+i*12])
		startWithSAP := 0xfe ^ data[30+i*12]
		sapType := data[31+i*12] << 1
		sapDeltaTime := binary.BigEndian.Uint32(data[34+i*12:30+i*12]) << 4
		b.References = append(b.References, Reference{
			ReferenceType:      refType,
			ReferencedSize:     refdSize,
			SubSegmentDuration: subSegDur,
			StartsWithSAP:      startWithSAP,
			SAPType:            sapType,
			SAPDeltaTime:       sapDeltaTime,
		})
	}
	return b, nil
}

func (b *SidxBox) Type() string {
	return "sidx"
}

func (b *SidxBox) Size() int {
	return BoxHeaderSize + 44
}

func (b *SidxBox) Dump() {

}

func (b *SidxBox) Encode(w io.Writer) error {

	return nil
}
