package mp4

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Box	Type:	 `tfdt’
// Container:	 Track	Fragment	box	(‘traf’)
// Mandatory:	 No
// Quantity:	 Zero	or	one
type TfdtBox struct {
	Version             byte
	Flags               [3]byte
	BaseMediaDecodeTime uint32
}

func (b *TfdtBox) Type() string {
	return "tfdt"
}

func (b *TfdtBox) Size() int {
	return BoxHeaderSize + 8
}
func (b *TfdtBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	buf := makebuf(b)
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	binary.BigEndian.PutUint32(buf[4:], b.BaseMediaDecodeTime)

	return err
}
func (b *TfdtBox) Dump() {
	fmt.Printf("Track fragment decode time\n")
	fmt.Printf("  Base Media Decode Time: %d\n", b.BaseMediaDecodeTime)
}

func DecodeTfdt(h BoxHeader, r io.Reader) (Box, error) {
	data := make([]byte, h.Size-BoxHeaderSize)
	_, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	return &TfdtBox{
		Version:             data[0],
		Flags:               [3]byte{data[1], data[2], data[3]},
		BaseMediaDecodeTime: binary.BigEndian.Uint32(data[4:8]),
	}, nil
}
