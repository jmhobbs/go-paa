package paa

import (
	"encoding/binary"
	"fmt"
	"io"
)

type PAA struct {
	Type    TypeOfPaX
	AVGC    *TaggAVGC
	MAXC    *TaggMAXC
	OFFS    *TaggOFFS
	Mipmaps []Mipmap
}

func Decode(in io.ReadSeeker) (*PAA, error) {
	img := &PAA{}

	// temporary storage
	var ulong uint32
	var ushort uint16

	// read first 4 bytes
	err := binary.Read(in, binary.LittleEndian, &ulong)
	if err != nil {
		return nil, err
	}

	// if we go straight to a TAGG, assume OFP index palette
	if ulong == TaggSignature {
		img.Type = Type_OFP

		err = readAndDecodeTagg(img, in)
		if err != nil {
			return nil, err
		}
	} else {
		img.Type = TypeOfPaX(ulong)

		// rewind for TAGG loop
		_, err = in.Seek(-2, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
	}

	for {
		err = binary.Read(in, binary.LittleEndian, &ulong)
		if err != nil {
			return nil, err
		}

		if ulong != TaggSignature {
			// not a TAGG, move on to palette
			_, err = in.Seek(-4, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
			break
		}

		err = readAndDecodeTagg(img, in)
		if err != nil {
			return nil, err
		}
	}

	// palette
	err = binary.Read(in, binary.LittleEndian, &ushort)
	if err != nil {
		return nil, err
	}
	if ushort != 0 {
		return nil, fmt.Errorf("error: paletted images not supported")
	}

	// mipmaps
	for {
		var mmHeader mipmapHeader

		err = binary.Read(in, binary.LittleEndian, &mmHeader)
		if err != nil {
			if err == io.EOF {
				break
			}
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

		mmSize := int64(uint32(mmSizeBytes[0]) | uint32(mmSizeBytes[1])<<8 | uint32(mmSizeBytes[2])<<16)

		offset, err := in.Seek(mmSize, io.SeekCurrent)
		if err != nil {
			return nil, err
		}

		img.Mipmaps = append(
			img.Mipmaps,
			Mipmap{
				Width:      mmHeader.Width,
				Height:     mmHeader.Height,
				Compressed: mmCompressed,
				Type:       img.Type,
				Size:       uint32(mmSize),
				Offset:     offset - mmSize,
			},
		)
	}

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
