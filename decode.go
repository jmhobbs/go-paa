package paa

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/anchore/go-lzo"
)

type PAA struct {
	Type    TypeOfPaX
	AVGC    *TaggAVGC
	MAXC    *TaggMAXC
	OFFS    *TaggOFFS
	Mipmaps []Mipmap
}

type mipmapHeader struct {
	Width  uint16
	Height uint16
}

type Mipmap struct {
	mipmapHeader
	Data       []byte
	Compressed bool
}

func (m *Mipmap) Reader() (io.Reader, error) {
	if !m.Compressed {
		return bytes.NewReader(m.Data), nil
	}
	return lzo.NewReader(bytes.NewReader(m.Data)), nil
}

func Decode(in io.Reader) (*PAA, error) {
	img := &PAA{}

	// temporary storage
	var ulong uint32
	var ushort uint16

	// read first 4 bytes
	err := binary.Read(in, binary.LittleEndian, &ulong)
	if err != nil {
		return nil, err
	}

	// if tagg, assume OFP index palette
	if ulong == TaggSignature {
		img.Type = Type_OFP

		err = readAndDecodeTagg(img, in)
		if err != nil {
			return nil, err
		}
	} else {
		img.Type = TypeOfPaX(ulong)

		if ulong&0xFFFF0000 != 0x47470000 {
			// not a TAGG next, we don't support this yet
			return nil, fmt.Errorf("error: A unsupported PAX type: %#x", img.Type)
		}

		// read 2 more bytes for the tagg header
		err := binary.Read(in, binary.LittleEndian, &ushort)
		if err != nil {
			return nil, err
		}

		if ushort == 0x5441 {
			// It is TAGG, decode it
			err = readAndDecodeTagg(img, in)
			if err != nil {
				return nil, err
			}
		} else {
			// also not a TAGG
			return nil, fmt.Errorf("error: unsupported PAX type: %#x", img.Type)
		}
	}

	for {
		// read first 2 bytes
		err = binary.Read(in, binary.LittleEndian, &ushort)
		if err != nil {
			return nil, err
		}

		if ushort != 0x4747 {
			// not a tag, move on to palette
			break
		}

		err = binary.Read(in, binary.LittleEndian, &ushort)
		if err != nil {
			return nil, err
		}
		if ushort != 0x5441 {
			// not a tagg? something is very wrong
			return nil, fmt.Errorf("error: expected remainder of TAGG signature, got: %#x", ushort)
		}

		err = readAndDecodeTagg(img, in)
		if err != nil {
			return nil, err
		}
	}

	// palette
	if ushort != 0 {
		return nil, fmt.Errorf("error: paletted images not supported (yet)")
	}

	// mipmaps
	var mmHeader mipmapHeader

	err = binary.Read(in, binary.LittleEndian, &mmHeader)
	if err != nil {
		return nil, err
	}

	var mmCompressed = mmHeader.Width&0x8000 == 0x8000
	if mmCompressed {
		mmHeader.Width = mmHeader.Width & 0x7FFF
	}

	if mmHeader.Width == 1234 && mmHeader.Height == 8765 {
		return nil, fmt.Errorf("error: paletted images not supported (yet)")
	}

	// size is a 24 bit unsigned
	mmSizeBytes := make([]uint8, 3)
	_, err = in.Read(mmSizeBytes)
	if err != nil {
		return nil, err
	}

	mmSize := uint32(mmSizeBytes[0]) | uint32(mmSizeBytes[1])<<8 | uint32(mmSizeBytes[2])<<16

	var mmData = make([]byte, mmSize)

	_, err = in.Read(mmData)
	if err != nil {
		return nil, err
	}

	img.Mipmaps = append(img.Mipmaps, Mipmap{mmHeader, mmData, mmCompressed})

	return img, nil
}

func readAndDecodeTagg(img *PAA, in io.Reader) error {
	var ulong uint32
	err := binary.Read(in, binary.LittleEndian, &ulong)
	if err != nil {
		return err
	}

	switch ulong {
	case Tagg_AVGC:
		img.AVGC, err = DecodeTaggAVGC(in)
	case Tagg_MAXC:
		img.MAXC, err = DecodeTaggMAXC(in)
	case Tagg_OFFS:
		img.OFFS, err = DecodeTaggOFFS(in)
	default:
		return fmt.Errorf("unknown TAGG: %#x", ulong)
	}

	return err
}
