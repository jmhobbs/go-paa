package paa_test

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/jmhobbs/go-paa"
	"github.com/jmhobbs/go-psnr"
	"github.com/stretchr/testify/require"
)

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

			require.Greater(t, db, 30.0, "PSNR should be greater than 30 dB for %s", format.Name)
		})
	}
}
