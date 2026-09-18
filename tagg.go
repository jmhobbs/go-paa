package paa

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/jmhobbs/go-paa/quantize"
)

var TaggSignature uint32 = 0x54414747

const (
	Tagg_AVGC uint32 = 0x41564743
	Tagg_MAXC uint32 = 0x4d415843
	Tagg_OFFS uint32 = 0x4f464653
	Tagg_SWIZ uint32 = 0x5357495a
)

/////////////////// AVGC

type TaggAVGC struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

func (t TaggAVGC) String() string {
	return fmt.Sprintf("R=%d, G=%d, B=%d, A=%d <#%2X%2X%2X%2X>", t.Red, t.Green, t.Blue, t.Alpha, t.Red, t.Green, t.Blue, t.Alpha)
}

func (t TaggAVGC) Write(out io.Writer) error {
	err := binary.Write(out, binary.LittleEndian, Tagg_AVGC)
	if err != nil {
		return err
	}
	err = binary.Write(out, binary.LittleEndian, uint32(4))
	if err != nil {
		return err
	}
	return binary.Write(out, binary.LittleEndian, t)
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

func CalculateAVGC(img image.Image) TaggAVGC {
	var avgc [3]byte
	switch img := img.(type) {
	case *image.NRGBA:
		avgc = quantize.MedianCut(fourChannelTothreeChannel(img.Pix))
	case *image.RGBA:
		avgc = quantize.MedianCut(fourChannelTothreeChannel(img.Pix))
	default:
		avgc = medianCutSlow(img)
	}

	return TaggAVGC{
		Red:   avgc[0],
		Green: avgc[1],
		Blue:  avgc[2],
		Alpha: 0xff,
	}
}

func medianCutSlow(img image.Image) [3]byte {
	buf := make([][3]byte, img.Bounds().Dx()*img.Bounds().Dy())

	for y := range img.Bounds().Dy() {
		rowoffset := y * img.Bounds().Dx()
		for x := range img.Bounds().Dx() {
			r, g, b, _ := img.At(x, y).RGBA()
			buf[rowoffset+x] = [3]byte{uint8(r), uint8(g), uint8(b)}
		}
	}

	return quantize.MedianCut(buf)
}

func fourChannelTothreeChannel(buf []byte) [][3]byte {
	out := make([][3]byte, len(buf)/4)
	for i := 0; i < len(buf); i += 4 {
		out[i/4] = [3]byte{buf[i], buf[i+1], buf[i+2]}
	}
	return out
}

/////////////////// MAXC

type TaggMAXC struct {
	Data [4]uint8
}

func (t TaggMAXC) Write(out io.Writer) error {
	err := binary.Write(out, binary.LittleEndian, Tagg_MAXC)
	if err != nil {
		return err
	}
	err = binary.Write(out, binary.LittleEndian, uint32(4))
	if err != nil {
		return err
	}
	return binary.Write(out, binary.LittleEndian, t)
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

/////////////////// OFFS

type TaggOFFS struct {
	Offsets [16]uint32
}

func (t TaggOFFS) Write(out io.Writer) error {
	err := binary.Write(out, binary.LittleEndian, Tagg_OFFS)
	if err != nil {
		return err
	}
	err = binary.Write(out, binary.LittleEndian, uint32(16*4))
	if err != nil {
		return err
	}
	return binary.Write(out, binary.LittleEndian, t)
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

/////////////////// SWIZ

type Channel uint8

const (
	Alpha Channel = 0
	Red   Channel = 1
	Green Channel = 2
	Blue  Channel = 3
)

func (c Channel) String() string {
	switch c {
	case Alpha:
		return "Alpha"
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Blue:
		return "Blue"
	default:
		return fmt.Sprintf("Unknown(%d)", uint8(c))
	}
}

type SwizzleChannel struct {
	Destination Channel
	Negated     bool
	ForceFF     bool
}

func DecodeSwizzleChannel(b uint8) SwizzleChannel {
	return SwizzleChannel{
		Destination: Channel(b & 0x03),
		Negated:     b&0x04 != 0,
		ForceFF:     b&0x08 != 0,
	}
}

func (c SwizzleChannel) String() string {
	if c.ForceFF {
		return "forced to 0xff"
	}
	if c.Negated {
		return fmt.Sprintf("negated, stored in %s", c.Destination)
	}
	return fmt.Sprintf("stored in %s", c.Destination)
}

type TaggSWIZ struct {
	Alpha uint8
	Red   uint8
	Green uint8
	Blue  uint8
}

func (t TaggSWIZ) AlphaChannel() SwizzleChannel { return DecodeSwizzleChannel(t.Alpha) }
func (t TaggSWIZ) RedChannel() SwizzleChannel   { return DecodeSwizzleChannel(t.Red) }
func (t TaggSWIZ) GreenChannel() SwizzleChannel { return DecodeSwizzleChannel(t.Green) }
func (t TaggSWIZ) BlueChannel() SwizzleChannel  { return DecodeSwizzleChannel(t.Blue) }

func (t TaggSWIZ) String() string {
	return fmt.Sprintf(
		"A=%s, R=%s, G=%s, B=%s",
		t.AlphaChannel(), t.RedChannel(), t.GreenChannel(), t.BlueChannel(),
	)
}

func (t TaggSWIZ) Write(out io.Writer) error {
	err := binary.Write(out, binary.LittleEndian, Tagg_SWIZ)
	if err != nil {
		return err
	}
	err = binary.Write(out, binary.LittleEndian, uint32(4))
	if err != nil {
		return err
	}
	return binary.Write(out, binary.LittleEndian, t)
}

func DecodeTaggSWIZ(in io.Reader) (*TaggSWIZ, error) {
	var length uint32
	err := binary.Read(in, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	if length != 4 {
		return nil, fmt.Errorf("error: unexpected length for SWIZ tag: %d", length)
	}

	var swiz TaggSWIZ
	err = binary.Read(in, binary.LittleEndian, &swiz)
	return &swiz, err
}
