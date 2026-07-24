package paa_test

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/jmhobbs/go-paa"
	"github.com/jmhobbs/go-psnr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingSeeker always fails Seek, to exercise error propagation in
// Mipmap.Reader and Mipmap.Image without needing real broken I/O.
type failingSeeker struct{}

func (failingSeeker) Read(_ []byte) (int, error) {
	return 0, errors.New("read should not be called")
}

func (failingSeeker) Seek(_ int64, _ int) (int64, error) {
	return 0, errors.New("seek failed")
}

func Test_Mipmap_Reader_SeekError(t *testing.T) {
	m := paa.Mipmap{Width: 1, Height: 1, Size: 1}

	_, err := m.Reader(failingSeeker{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seek failed")
}

func Test_Mipmap_Image_ReaderError(t *testing.T) {
	m := paa.Mipmap{Width: 1, Height: 1, Size: 1, Type: paa.Type_DXT1}

	_, err := m.Image(failingSeeker{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seek failed")
}

func Test_Mipmap_Image_UnsupportedType(t *testing.T) {
	m := paa.Mipmap{
		Width:  1,
		Height: 1,
		Size:   1,
		Type:   paa.Type_RGBA8,
	}

	_, err := m.Image(bytes.NewReader([]byte{0x00}))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported mipmap type")
}

func Test_Mipmap_Image_DecodeError(t *testing.T) {
	// Width/Height imply far more data than is actually present, which
	// should surface as a decode error rather than a panic.
	m := paa.Mipmap{
		Width:  4,
		Height: 4,
		Size:   0,
		Type:   paa.Type_DXT1,
	}

	_, err := m.Image(bytes.NewReader(nil))
	require.Error(t, err)
}

func Test_Mipmap_Image(t *testing.T) {
	expected, err := func() (image.Image, error) {
		f, err := os.Open("testdata/test-pattern.png")
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := f.Close(); err != nil {
				t.Logf("failed to close test input file: %v", err)
			}
		}()

		return png.Decode(f)
	}()
	require.NoError(t, err)

	formats := []struct {
		Name string
		Type paa.TypeOfPaX
		File string
	}{
		{
			"DXT1", paa.Type_DXT1, "testdata/test-pattern.dxt1",
		},
		{
			"DXT3", paa.Type_DXT3, "testdata/test-pattern.dxt3",
		},
		{
			"DXT5", paa.Type_DXT5, "testdata/test-pattern.dxt5",
		},
	}
	for _, format := range formats {

		t.Run(format.Name, func(t *testing.T) {
			finfo, err := os.Stat(format.File)
			require.NoError(t, err)

			f, err := os.Open(format.File)
			require.NoError(t, err)
			defer f.Close()

			m := paa.Mipmap{
				Width:      512,
				Height:     256,
				Compressed: false,
				Type:       format.Type,
				Size:       uint32(finfo.Size()),
				Offset:     0,
			}

			rgba, err := m.Image(f)
			require.NoError(t, err)

			db, err := psnr.Image(expected, rgba)
			require.NoError(t, err)

			require.Greater(t, db, 40.0, "PSNR should be greater than 40dB for %s", format.Name)
		})
	}
}
