package mp4

import (
	"fmt"
	"io"
)

// (styp - mandatory for fMP4)
//
// Status: decoded
type StypBox struct {
	MajorBrand       string
	MinorVersion     []byte
	CompatibleBrands []string
}

func DecodeStyp(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	b := &StypBox{
		MajorBrand:       string(data[0:4]),
		MinorVersion:     data[4:8],
		CompatibleBrands: []string{},
	}
	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(data[i:i+4]))
		}
	}
	return b, nil
}

func (b *StypBox) Type() string {
	return "styp"
}

func (b *StypBox) Size() int {
	return BoxHeaderSize + 8 + 4*len(b.CompatibleBrands)
}

func (b *StypBox) Dump() {
	fmt.Printf("Box type: %s\n", b.Type())
	fmt.Printf("+- Major brand: %s\n", b.MajorBrand)
	fmt.Printf("+- Minor version: %s\n", b.MinorVersion)
	fmt.Printf("+- Compatible brands: sizes = %d\n", len(b.CompatibleBrands))
	for i, e := range b.CompatibleBrands {
		fmt.Printf(" +- [%d]\t: %s\n", i, e)
	}
}

func (b *StypBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	strtobuf(buf, b.MajorBrand, 4)
	copy(buf[4:], b.MinorVersion)
	for i, c := range b.CompatibleBrands {
		strtobuf(buf[8+i*4:], c, 4)
	}
	_, err = w.Write(buf)
	return err
}
