package paa

import (
	"bytes"
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
	mipmapHeader
	Data       []byte
	Compressed bool
	Type       TypeOfPaX
}

// Get access to the raw mipmap texture data.
func (m *Mipmap) Reader() (io.Reader, error) {
	if !m.Compressed {
		return bytes.NewReader(m.Data), nil
	}
	return lzo.NewReader(bytes.NewReader(m.Data)), nil
}

func (m *Mipmap) Image() (*image.RGBA, error) {
	if m.Type != Type_DXT1 {
		return nil, fmt.Errorf("error: unsupported mipmap type: %#x", TypeOfPaXStrings[m.Type])
	}

	var (
		rgbaBytes []byte
		err       error
	)

	if m.Compressed {
		// DXT1
		blockCountX := (int(m.Width) + 3) / 4
		blockCountY := (int(m.Height) + 3) / 4
		dxt1Bytes := make([]byte, blockCountX*blockCountY*8)
		_, err = lzo.Decompress(m.Data, dxt1Bytes)
		if err != nil {
			return nil, err
		}
		rgbaBytes, err = dxt.DecodeDXT1(dxt1Bytes, uint(m.Width), uint(m.Height))
	} else {
		rgbaBytes, err = dxt.DecodeDXT1(m.Data, uint(m.Width), uint(m.Height))
	}
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))
	rgba.Pix = rgbaBytes

	return rgba, nil
}
