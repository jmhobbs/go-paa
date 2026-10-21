package paa

import (
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/mauserzjeh/dxt"
	"github.com/nfnt/resize"
)

type paaBuilder struct {
	source    image.Image
	typeOfPaX TypeOfPaX
}

type builderOption func(*paaBuilder) error

func WithType(t TypeOfPaX) builderOption {
	return func(b *paaBuilder) error {
		b.typeOfPaX = t
		switch t {
		case Type_DXT1, Type_DXT3, Type_DXT5:
			return nil
		default:
			return fmt.Errorf("currently unsupported type: %v", t)
		}
	}
}

func New(img image.Image, options ...builderOption) (*paaBuilder, error) {
	builder := &paaBuilder{img, Type_DXT1}

	var err error
	for _, option := range options {
		if err = option(builder); err != nil {
			return nil, err
		}
	}

	return builder, nil
}

func (b *paaBuilder) Write(sink io.Writer) error {
	// write type
	err := binary.Write(sink, binary.LittleEndian, uint32(b.typeOfPaX))
	if err != nil {
		return err
	}

	// write TAGG's
	err = binary.Write(sink, binary.LittleEndian, TaggSignature)
	if err != nil {
		return err
	}

	// AVGC
	err = CalculateAVGC(b.source).Write(sink)
	if err != nil {
		return err
	}

	// MAXC
	err = NewTaggMAXC().Write(sink)
	if err != nil {
		return err
	}

	mipmaps := calculateMipmaps(b.source, b.typeOfPaX)
	if err != nil {
		return err
	}

	mipmapdata := make([][]byte, len(mipmaps))

	for i := range len(mipmaps) {
		resized := resize.Resize(uint(mipmaps[i].Width), uint(mipmaps[i].Height), b.source, resize.Lanczos3)
		var out []byte
		switch b.typeOfPaX {
		case Type_DXT1:
			out, err = dxt.EncodeDXT1(rgbaData(resized), uint(mipmaps[i].Width), uint(mipmaps[i].Height))
		case Type_DXT3:
			out, err = dxt.EncodeDXT3(rgbaData(resized), uint(mipmaps[i].Width), uint(mipmaps[i].Height))
		case Type_DXT5:
			out, err = dxt.EncodeDXT5(rgbaData(resized), uint(mipmaps[i].Width), uint(mipmaps[i].Height))
		}
		if err != nil {
			return err
		}

		if len(out) > 1024 {
			// compress it!
		}

		mipmapdata[i] = make([]byte, len(out))
		copy(mipmapdata[i], out)
	}

	// OFFS
	err = TaggOFFS{}.Write(sink)
	if err != nil {
		return err
	}

	return nil
}

func calculateMipmaps(img image.Image, typeOfPaX TypeOfPaX) []Mipmap {
	mipmaps := []Mipmap{}
	mipmaps = append(mipmaps, Mipmap{
		Width:  uint16(img.Bounds().Dx()),
		Height: uint16(img.Bounds().Dy()),
		Type:   typeOfPaX,
	})
	x := uint16(img.Bounds().Dx())
	y := uint16(img.Bounds().Dy())

	for x > 4 && y > 4 {
		x = x / 2
		y = y / 2
		mipmaps = append(mipmaps, Mipmap{
			Width:  x,
			Height: y,
			Type:   typeOfPaX,
		})
	}

	return mipmaps
}

func rgbaData(img image.Image) []byte {
	switch img := img.(type) {
	case *image.RGBA:
		return img.Pix // TODO: undo pre-multiplied alpha
	case *image.NRGBA:
		return img.Pix
	}

	out := make([]byte, img.Bounds().Dx()*img.Bounds().Dy()*4)
	for y := range img.Bounds().Dy() {
		rowoffset := y * img.Bounds().Dx() * 4
		for x := range img.Bounds().Dx() {
			r, g, b, a := img.At(x, y).RGBA()
			out[rowoffset+x*4] = uint8(r)
			out[rowoffset+x*4+1] = uint8(g)
			out[rowoffset+x*4+2] = uint8(b)
			out[rowoffset+x*4+3] = uint8(a)
		}
	}

	return out
}
