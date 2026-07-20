package paa

import (
	"encoding/binary"
	"fmt"
	"io"
)

var TaggSignature uint32 = 0x54414747

const (
	Tagg_AVGC uint32 = 0x41564743
	Tagg_MAXC uint32 = 0x4d415843
	Tagg_OFFS uint32 = 0x4f464653
)

type TaggAVGC struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

type TaggMAXC struct {
	Data [4]uint8
}

type TaggOFFS struct {
	Offsets [16]uint32
}

func DecodeTaggAVGC(in io.Reader) (*TaggAVGC, error) {
	var length uint32
	err := binary.Read(in, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	if length != 4 {
		return nil, fmt.Errorf("error: unexpected length for AVGC tag: %d", length)
	}

	var avgc TaggAVGC
	err = binary.Read(in, binary.LittleEndian, &avgc)
	return &avgc, err
}

func DecodeTaggMAXC(in io.Reader) (*TaggMAXC, error) {
	var length uint32
	err := binary.Read(in, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	if length != 4 {
		return nil, fmt.Errorf("error: unexpected length for MAXC tag: %d", length)
	}

	var maxc TaggMAXC
	err = binary.Read(in, binary.LittleEndian, &maxc)
	return &maxc, err
}

func DecodeTaggOFFS(in io.Reader) (*TaggOFFS, error) {
	var length uint32
	err := binary.Read(in, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	if length != 64 {
		return nil, fmt.Errorf("error: unexpected length for OFFS tag: %d", length)
	}

	var offs TaggOFFS
	err = binary.Read(in, binary.LittleEndian, &offs)
	return &offs, err
}
