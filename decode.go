package paa

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type PAA struct {
	Type TypeOfPaX
	AVGC *TaggAVGC
	MAXC *TaggMAXC
	OFFS *TaggOFFS
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

	return img, nil
}

func readAndDecodeTagg(img *PAA, in io.Reader) error {
	var ulong uint32
	err := binary.Read(in, binary.LittleEndian, &ulong)
	if err != nil {
		return err
	}

	log.Printf("data length: %d", ulong)

	switch ulong {
	case Tagg_AVGC:
		img.AVGC, err = decodeTaggAVGC(in)
	case Tagg_MAXC:
		img.MAXC, err = decodeTaggMAXC(in)
	case Tagg_OFFS:
		img.OFFS, err = decodeTaggOFFS(in)
	default:
		return fmt.Errorf("unknown TAGG: %#x", ulong)
	}

	return err
}
