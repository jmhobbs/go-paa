package paa

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/jmhobbs/go-paa/formats"

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

	compressed, err := io.ReadAll(&reader)
	if err != nil {
		return nil, err
	}

	size, err := m.decompressedSize()
	if err != nil {
		return nil, err
	}

	// lzo.NewReader as a bound buffer size of 64kb, which
	// PAA images can easily overrun, so we have to eat it all
	// in one go instead.
	decompressed := make([]byte, size)
	n, err := lzo.Decompress(compressed, decompressed)
	if err != nil {
		return nil, fmt.Errorf("error: failed to decompress lzo mipmap data: %w", err)
	}

	return bytes.NewReader(decompressed[:n]), nil
}

func (m *Mipmap) decompressedSize() (int, error) {
	switch m.Type {
	case Type_DXT1:
		return blockCount(m.Width, m.Height) * 8, nil
	case Type_DXT2, Type_DXT3, Type_DXT4, Type_DXT5:
		return blockCount(m.Width, m.Height) * 16, nil
	case Type_RGBA4, Type_RGBA5:
		return int(m.Width) * int(m.Height) * 2, nil
	case Type_RGBA8:
		return int(m.Width) * int(m.Height) * 4, nil
	case Type_Gray:
		return int(m.Width) * int(m.Height), nil
	default:
		return 0, fmt.Errorf("error: unsupported mipmap type: %s", m.Type)
	}
}

func blockCount(width, height uint16) int {
	blockCountX := (int(width) + 3) / 4
	blockCountY := (int(height) + 3) / 4
	return blockCountX * blockCountY
}

func (m *Mipmap) Image(src io.ReadSeeker) (*image.NRGBA, error) {
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
	case Type_RGBA4:
		fmt.Fprintln(os.Stderr, "warning: RGBA4 is not fully supported and may produce incorrect results")
		rgbaBytes, err = formats.DecodeRGBA4(data, uint(m.Width), uint(m.Height))
	case Type_RGBA5:
		fmt.Fprintln(os.Stderr, "warning: RGBA5 is not fully supported and may produce incorrect results")
		rgbaBytes, err = formats.DecodeRGBA5(data, uint(m.Width), uint(m.Height))
	default:
		return nil, fmt.Errorf("error: unsupported mipmap type: %s", m.Type)
	}
	if err != nil {
		return nil, err
	}

	rgba := image.NewNRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))
	rgba.Pix = rgbaBytes

	return rgba, nil
}
