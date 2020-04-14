package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// File Type Box (ftyp - mandatory)
//
// Status: decoded
type FtypBox struct {
	MajorBrand       string
	MinorVersion     uint32
	CompatibleBrands []string
}

func DecodeFtyp(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &FtypBox{
		MajorBrand:       string(data[0:4]),
		MinorVersion:     binary.BigEndian.Uint32(data[4:8]),
		CompatibleBrands: []string{},
	}
	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(data[i:i+4]))
		}
	}
	return b, nil
}

func (b *FtypBox) Box() Box {
	return b
}

func (b *FtypBox) Type() string {
	return "ftyp"
}

func (b *FtypBox) Size() int {
	return BoxHeaderSize + 8 + 4*len(b.CompatibleBrands)
}

func (b *FtypBox) Dump() {
	fmt.Printf("Box type: %s\n", b.Type())
	fmt.Printf("+- Major brand: %s\n", b.MajorBrand)
	fmt.Printf("+- Minor version: 0x%x\n", b.MinorVersion)
	fmt.Printf("+- Compatible brands: sizes = %d\n", len(b.CompatibleBrands))
	for i, e := range b.CompatibleBrands {
		fmt.Printf(" +- [%d]\t: %s\n", i, e)
	}
}

func (b *FtypBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	strtobuf(buf, b.MajorBrand, 4)
	binary.BigEndian.PutUint32(buf[4:], b.MinorVersion)
	for i, c := range b.CompatibleBrands {
		strtobuf(buf[8+i*4:], c, 4)
	}
	_, err = w.Write(buf)
	return err
}
