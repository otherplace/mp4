package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 ‘iloc’
// Container:	 Meta	box	(‘meta’)
// Mandatory:	No
// Quantity:	 Zero	or	one
type IlocBox struct {
	Version        byte
	Flags          [3]byte
	OffsetSize     uint8
	LenghSize      uint8
	BaseOffsetSize uint8
	IndexSize      uint8
	ItemCount      uint16
	Items          []*Item
}
type Item struct {
	ItemId             uint16
	Reserved           uint16
	ConstructionMethod uint8
	DataReferenceIndex uint16
	BaseOffset         uint32
	ExtentCount        uint16
	Extents            []*Extent
}

type Extent struct {
	ExtentIndex  uint32
	ExtentOffset uint32
	ExtentLength uint32
}

func DecodeIloc(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &IlocBox{
		Version: data[0],
		Flags:   [3]byte{data[1], data[2], data[3]},
	}
	offsetLength := data[4]
	b.OffsetSize = (0xf0 ^ offsetLength) >> 4
	b.LenghSize = (0x0f ^ offsetLength)
	baseIndex := data[5]
	b.BaseOffsetSize = (0xf0 ^ baseIndex) >> 4
	b.IndexSize = (0x0f ^ baseIndex)
	b.ItemCount = binary.BigEndian.Uint16(data[6:8])
	offset := 8
	for i := 0; i < int(b.ItemCount); i++ {
		bi := &Item{}
		bi.ItemId = binary.BigEndian.Uint16(data[offset+i*2 : offset+2+i*2])
		offset += 2
		revCon := binary.BigEndian.Uint16(data[offset+i*2 : offset+2+i*2])
		offset += 2
		bi.Reserved = (0xf0 ^ revCon) >> 4
		bi.ConstructionMethod = uint8(0x0f ^ revCon)
		bi.DataReferenceIndex = binary.BigEndian.Uint16(data[offset+i*2 : offset+2+i*2])
		offset += 2
		bi.BaseOffset = binary.BigEndian.Uint32(data[offset+i*4 : offset+4+i*4])
		offset += 4
		bi.ExtentCount = binary.BigEndian.Uint16(data[offset+i*2 : offset+2+i*2])
		offset += 2
		offset2 := offset
		for j := 0; i < int(bi.ExtentCount); j++ {
			e := &Extent{}
			e.ExtentIndex = binary.BigEndian.Uint32(data[offset2+i*4 : offset2+4+i*4])
			e.ExtentOffset = binary.BigEndian.Uint32(data[offset2+i*4 : offset2+4+i*4])
			e.ExtentLength = binary.BigEndian.Uint32(data[offset2+i*4 : offset2+4+i*4])
			bi.Extents = append(bi.Extents, e)
		}
		b.Items = append(b.Items, bi)
	}
	return b, nil
}

func (b *IlocBox) Box() Box {
	return b
}

func (b *IlocBox) Type() string {
	return "iloc"
}

func (b *IlocBox) Size() int {
	l := BoxHeaderSize + 8
	for _, i := range b.Items {
		l = l + 12 + int(i.ExtentCount)*24
	}

	return l
}

func (b *IlocBox) Encode(w io.Writer) error {
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

func (b *IlocBox) Dump() {
	fmt.Printf("Item Location Box\n")
}
