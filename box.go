package mp4

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
)

const (
	BoxHeaderSize = 8
)

var (
	ErrUnknownBoxType  = errors.New("unknown box type")
	ErrTruncatedHeader = errors.New("truncated header")
	ErrTruncatedBody   = errors.New("truncated body")
	ErrBadFormat       = errors.New("bad format")
)

var decoders map[string]BoxDecoder

func init() {
	decoders = map[string]BoxDecoder{
		"cprt": DecodeCprt,
		"ftyp": DecodeFtyp,
		"styp": DecodeStyp,
		"moof": DecodeMoof,
		"moov": DecodeMoov,
		"mvhd": DecodeMvhd,
		"iods": DecodeIods,
		"trak": DecodeTrak,
		//"udta": DecodeUdta,
		"tkhd": DecodeTkhd,
		"edts": DecodeEdts,
		"elst": DecodeElst,
		"mdia": DecodeMdia,
		"minf": DecodeMinf,
		"mdhd": DecodeMdhd,
		"mfhd": DecodeMfhd,
		"mvex": DecodeMvex,
		"hdlr": DecodeHdlr,
		"vmhd": DecodeVmhd,
		"smhd": DecodeSmhd,
		"dinf": DecodeDinf,
		"dref": DecodeDref,
		"pdin": DecodePdin,
		"sidx": DecodeSidx,
		"stbl": DecodeStbl,
		"stco": DecodeStco,
		"stsc": DecodeStsc,
		"stsz": DecodeStsz,
		"ctts": DecodeCtts,
		"stsd": DecodeStsd,
		"stts": DecodeStts,
		"stss": DecodeStss,
		"traf": DecodeTraf,
		"tfdt": DecodeTfdt,
		"tfhd": DecodeTfhd,
		"trex": DecodeTrex,
		"trun": DecodeTrun,
		"meta": DecodeMeta,
		"mdat": DecodeMdat,
		"free": DecodeFree,
	}
}

// The header of a box
type BoxHeader struct {
	Type string
	Size uint32
}

// DecodeHeader decodes a box header (size + box type)
func DecodeHeader(r io.Reader) (BoxHeader, error) {
	buf := make([]byte, BoxHeaderSize)
	n, err := r.Read(buf)
	if err != nil {
		return BoxHeader{}, err
	}
	if n != BoxHeaderSize {
		return BoxHeader{}, ErrTruncatedHeader
	}
	return BoxHeader{string(buf[4:8]), binary.BigEndian.Uint32(buf[0:4])}, nil
}

// EncodeHeader encodes a box header to a writer
func EncodeHeader(b Box, w io.Writer) error {
	buf := make([]byte, BoxHeaderSize)
	binary.BigEndian.PutUint32(buf, uint32(b.Size()))
	strtobuf(buf[4:], b.Type(), 4)
	_, err := w.Write(buf)
	return err
}

// A box
type Box interface {
	Box() Box
	Type() string
	Size() int
	Encode(w io.Writer) error
}

type UkwnBox struct {
	h           BoxHeader
	PayloadSize uint32
	Data        []byte
}

func DecodeUkwnBox(h BoxHeader, r io.Reader) (Box, error) {
	if lr, limited := r.(*io.LimitedReader); limited {
		r = lr.R
	}
	data := make([]byte, h.Size-BoxHeaderSize)
	n, _ := r.Read(data)

	b := &UkwnBox{
		Data:        data,
		PayloadSize: uint32(n),
	}
	return b, nil
}

func (b *UkwnBox) Type() string {
	return b.h.Type
}

func (b *UkwnBox) Size() int {
	return int(b.PayloadSize)
}

func (b *UkwnBox) Dump() {
	fmt.Printf("Box type: %s\n", b.Type())
	fmt.Printf(" Data length: %d\n", len(b.Data))
}
func (b *UkwnBox) Encode(w io.Writer) error {
	return nil
}

func (b *UkwnBox) Box() Box {
	return b
}

type BoxDecoder func(h BoxHeader, r io.Reader) (Box, error)

// DecodeBox decodes a box
func DecodeBox(h BoxHeader, r io.Reader) (Box, error) {
	d := decoders[h.Type]
	if d == nil {
		log.Printf("Error while decoding %s: unknown box type, len: %d", h.Type, h.Size)
		d = DecodeUkwnBox
	}
	b, err := d(h, io.LimitReader(r, int64(h.Size-BoxHeaderSize)))
	if err != nil {
		log.Printf("Error while decoding %v:%s", h, err)
		return nil, err
	}
	return b, nil
}

// DecodeContainer decodes a container box
func DecodeContainer(r io.Reader) ([]Box, error) {
	l := []Box{}
	for {
		h, err := DecodeHeader(r)
		if err == io.EOF {
			return l, nil
		}
		log.Printf("Decode header, %v\n", h)
		if err != nil {
			log.Printf("Decode header fail, %s:%v\n", err, h)
			return l, err
		}
		b, err := DecodeBox(h, r)
		if err != nil {
			log.Printf("Decode Box fail, %s:%v\n", err, h)
			return l, err
		}
		log.Printf("Decode box, %v\n", h)
		l = append(l, b)
	}
}

// An 8.8 fixed point number
type Fixed16 uint16

func (f Fixed16) String() string {
	return fmt.Sprintf("%d.%d", uint16(f)>>8, uint16(f)&7)
}

func fixed16(bytes []byte) Fixed16 {
	return Fixed16(binary.BigEndian.Uint16(bytes))
}

func putFixed16(bytes []byte, i Fixed16) {
	binary.BigEndian.PutUint16(bytes, uint16(i))
}

// A 16.16 fixed point number
type Fixed32 uint32

func (f Fixed32) String() string {
	return fmt.Sprintf("%d.%d", uint32(f)>>16, uint32(f)&15)
}

func fixed32(bytes []byte) Fixed32 {
	return Fixed32(binary.BigEndian.Uint32(bytes))
}

func putFixed32(bytes []byte, i Fixed32) {
	binary.BigEndian.PutUint32(bytes, uint32(i))
}

func strtobuf(out []byte, str string, l int) {
	in := []byte(str)
	if l < len(in) {
		copy(out, in)
	} else {
		copy(out, in[0:l])
	}
}

func makebuf(b Box) []byte {
	return make([]byte, b.Size()-BoxHeaderSize)
}

// utilities
func BEUint28(b []byte) uint32 {
	_ = b[2]
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}
func BEPutUint28(b []byte, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}

func compareFlag(flag uint32, bitmask uint32) bool {
	return flag&bitmask == bitmask
}
