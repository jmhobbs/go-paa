package paa

import (
	"fmt"
	"image"
	"io"

	"github.com/anchore/go-lzo"
	"github.com/mauserzjeh/dxt"
)

type mipmapHeader struct {
	Width  uint16
	Height uint16
}

type Mipmap struct {
	Width      uint16
	Height     uint16
	Compressed bool
	Type       TypeOfPaX
	Size       uint32
	Offset     int64
}

// Get access to the raw mipmap texture data.
func (m *Mipmap) Reader(src io.ReadSeeker) (io.Reader, error) {
	_, err := src.Seek(m.Offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	reader := io.LimitedReader{R: src, N: int64(m.Size)}

	if !m.Compressed {
		return &reader, nil
	}
	return lzo.NewReader(&reader), nil
}

func (m *Mipmap) Image(src io.ReadSeeker) (*image.RGBA, error) {
	var (
		rgbaBytes []byte
		err       error
	)

	reader, err := m.Reader(src)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	switch m.Type {
	case Type_DXT1:
		rgbaBytes, err = dxt.DecodeDXT1(data, uint(m.Width), uint(m.Height))
	case Type_DXT3:
		rgbaBytes, err = dxt.DecodeDXT3(data, uint(m.Width), uint(m.Height))
	case Type_DXT5:
		rgbaBytes, err = dxt.DecodeDXT5(data, uint(m.Width), uint(m.Height))
	default:
		return nil, fmt.Errorf("error: unsupported mipmap type: %#x", TypeOfPaXStrings[m.Type])
	}
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))
	rgba.Pix = rgbaBytes

	return rgba, nil
}
