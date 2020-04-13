package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// Box Type: ‘mehd’
// Container: Movie Extends Box(‘mvex’)
// Mandatory: No
// Quantity: Zero or one
type MehdBox struct {
	Version          byte
	Flags            []byte
	FragmentDuration uint32
}

func DecodeMehdBox(r io.Reader) (Box, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &MehdBox{
		Version:          data[0],
		Flags:            []byte{data[1], data[2], data[3]},
		FragmentDuration: binary.BigEndian.Uint32(data[4:8]),
	}, nil
}

func (b *MehdBox) Type() string {
	return "mehd"
}

func (b *MehdBox) Size() int {
	return BoxHeaderSize + 4 + 4
}

func (b *MehdBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.FragmentDuration)
	_, err = w.Write(buf)
	return err
}

func (b *MehdBox) Dump() {
	fmt.Printf("Movie Extends Header Box\n")
	fmt.Printf(" Fragment Duration: %d\n", b.FragmentDuration)

}
