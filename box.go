package mp4

import (
	"bytes"
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
		"bxml": DecodeBxml,
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
		"iloc": DecodeIloc,
		"mdia": DecodeMdia,
		"minf": DecodeMinf,
		"mdhd": DecodeMdhd,
		"mfhd": DecodeMfhd,
		"mfra": DecodeMfra,
		"mfro": DecodeMfro,
		"mvex": DecodeMvex,
		"hdlr": DecodeHdlr,
		"vmhd": DecodeVmhd,
		"smhd": DecodeSmhd,
		"dinf": DecodeDinf,
		"dref": DecodeDref,
		"pdin": DecodePdin,
		"sbgp": DecodeSbgp,
		"sidx": DecodeSidx,
		"stbl": DecodeStbl,
		"stco": DecodeStco,
		"co64": DecodeCo64,
		"stsc": DecodeStsc,
		"stsz": DecodeStsz,
		"ctts": DecodeCtts,
		"stsd": DecodeStsd,
		"stts": DecodeStts,
		"stss": DecodeStss,
		"traf": DecodeTraf,
		"tfdt": DecodeTfdt,
		"tfhd": DecodeTfhd,
		"tfra": DecodeTfra,
		"trex": DecodeTrex,
		"trun": DecodeTrun,
		"meta": DecodeMeta,
		"mdat": DecodeMdat,
		"free": DecodeFree,
		"sgpd": DecodeSgpd,
	}
}

// The header of a box
type BoxHeader struct {
	Type      string
	Size      uint32
	LargeSize uint64
}

type FullBox struct {
	Version byte
	Flags   [3]byte
}

func (b *FullBox) Size() int {
	return 4
}

func DecodeFullBox(data []byte) FullBox {
	return FullBox{
		Version: data[0],
		Flags:   [3]byte{data[1], data[2], data[3]},
	}
}

func EncodeFullBox(b FullBox, buf []byte) error {
	buf[0] = b.Version
	buf[1], buf[2], buf[3] = b.Flags[0], b.Flags[1], b.Flags[2]
	return nil
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
	var lsz uint64
	sz := binary.BigEndian.Uint32(buf[0:4])
	if sz == 0 {
		secBuf := &bytes.Buffer{}
		nRead, err := io.Copy(secBuf, r)
		if err != nil {
			return BoxHeader{}, ErrTruncatedHeader
		}
		sz = uint32(nRead)
	} else if sz == 1 {
		secBuf := make([]byte, BoxHeaderSize)
		n, err := r.Read(secBuf)
		if err != nil {
			return BoxHeader{}, err
		}
		if n != BoxHeaderSize {
			return BoxHeader{}, ErrTruncatedHeader
		}
		lsz = binary.BigEndian.Uint64(secBuf)
	}
	return BoxHeader{string(buf[4:8]), sz, lsz}, nil
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
	Dump()
}

type BoxDecoder func(h BoxHeader, r io.Reader) (Box, error)

// DecodeBox decodes a box
func DecodeBox(h BoxHeader, r io.Reader) (Box, error) {
	d := decoders[h.Type]
	if d == nil {
		log.Printf("Error while decoding %s: unknown box type, len: %d", h.Type, h.Size)
		d = DecodeUkwnBox
	}
	var readSize int64
	readSize = int64(h.Size - BoxHeaderSize)
	if h.Size == 1 {
		readSize = int64(h.LargeSize - (BoxHeaderSize * 2))
	}
	b, err := d(h, io.LimitReader(r, readSize))
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
