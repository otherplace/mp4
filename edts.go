package mp4

import "io"

// Edit Box (edts - optional)
//
// Contained in: Track Box ("trak")
//
// Status: decoded
//
// The edit box maps the presentation timeline to the media-time line
type EdtsBox struct {
	Elst *ElstBox `json:"elst,omitempty"`
}

func DecodeEdts(h BoxHeader, r io.Reader) (Box, error) {
	l, err := DecodeContainer(r)
	if err != nil {
		return nil, err
	}
	e := &EdtsBox{}
	for _, b := range l {
		switch b.Type() {
		case "elst":
			e.Elst = b.(*ElstBox)
		default:
			return nil, ErrBadFormat
		}
	}
	return e, nil
}

func (b *EdtsBox) Box() Box {
	return b
}
func (b *EdtsBox) Type() string {
	return "edts"
}

func (b *EdtsBox) Size() int {
	return BoxHeaderSize + b.Elst.Size()
}

func (b *EdtsBox) Dump() {
	b.Elst.Dump()
}

func (b *EdtsBox) Encode(w io.Writer) error {
	err := EncodeHeader(b, w)
	if err != nil {
		return err
	}
	return b.Elst.Encode(w)
}
