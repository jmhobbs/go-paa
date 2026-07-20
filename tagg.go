package paa

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var TaggSignature uint32 = 0x47474154

type Tagg struct {
	Name   [4]byte
	Length uint32
	Data   []byte
}

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
	Length uint32
	Data   [4]byte
}

type TaggOFFS struct {
	Length  uint32
	Offsets [16]uint32
}

func decodeTaggAVGC(in io.Reader) (*TaggAVGC, error) {
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

func decodeTaggMAXC(in io.Reader) (*TaggMAXC, error) {
	return nil, errors.New("not implemented")

}
func decodeTaggOFFS(in io.Reader) (*TaggOFFS, error) {
	return nil, errors.New("not implemented")
}
