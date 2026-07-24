package paa_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/jmhobbs/go-paa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DecodeTaggAVGC(t *testing.T) {
	src := bytes.NewReader([]byte{
		0x04, 0x00, 0x00, 0x00, // Length: 4
		0xb6, // Red
		0xaa, // Green
		0xa9, // Blue
		0xff, // Alpha
	})

	avgc, err := paa.DecodeTaggAVGC(src)
	require.NoError(t, err)
	require.NotNil(t, avgc)
	assert.Equal(t, uint8(0xb6), avgc.Red)
	assert.Equal(t, uint8(0xaa), avgc.Green)
	assert.Equal(t, uint8(0xa9), avgc.Blue)
	assert.Equal(t, uint8(0xff), avgc.Alpha)
}

func Test_DecodeTaggMAXC(t *testing.T) {
	src := bytes.NewReader([]byte{
		0x04, 0x00, 0x00, 0x00, // Length: 4
		0xff, 0xff, 0xff, 0xff, // Data
	})

	maxc, err := paa.DecodeTaggMAXC(src)
	require.NoError(t, err)
	require.NotNil(t, maxc)
	assert.Equal(t, [4]uint8{0xff, 0xff, 0xff, 0xff}, maxc.Data)
}

func Test_DecodeTagg_LengthReadError(t *testing.T) {
	decoders := map[string]func(io.Reader) (any, error){
		"AVGC": func(r io.Reader) (any, error) { return paa.DecodeTaggAVGC(r) },
		"MAXC": func(r io.Reader) (any, error) { return paa.DecodeTaggMAXC(r) },
		"OFFS": func(r io.Reader) (any, error) { return paa.DecodeTaggOFFS(r) },
	}

	for name, decode := range decoders {
		t.Run(name, func(t *testing.T) {
			// Empty reader: even the 4 byte length prefix can't be read.
			_, err := decode(bytes.NewReader(nil))
			require.Error(t, err)
		})
	}
}

func Test_DecodeTagg_LengthMismatch(t *testing.T) {
	cases := []struct {
		name   string
		decode func(io.Reader) (any, error)
	}{
		{"AVGC", func(r io.Reader) (any, error) { return paa.DecodeTaggAVGC(r) }},
		{"MAXC", func(r io.Reader) (any, error) { return paa.DecodeTaggMAXC(r) }},
		{"OFFS", func(r io.Reader) (any, error) { return paa.DecodeTaggOFFS(r) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := bytes.NewReader([]byte{0x01, 0x00, 0x00, 0x00}) // Length: 1, wrong for all three tags
			result, err := c.decode(src)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "unexpected length")
			assert.Nil(t, result)
		})
	}
}

func Test_DecodeTaggOFFS(t *testing.T) {
	src := bytes.NewReader([]byte{
		0x40, 0x00, 0x00, 0x00, // Length: 64 (16*4)
		0x01, 0x00, 0x00, 0x00, // Offsets[0]
		0x02, 0x00, 0x00, 0x00, // Offsets[1]
		0x03, 0x00, 0x00, 0x00, // Offsets[2]
		0x04, 0x00, 0x00, 0x00, // Offsets[3]
		0x05, 0x00, 0x00, 0x00, // Offsets[4]
		0x06, 0x00, 0x00, 0x00, // Offsets[5]
		0x07, 0x00, 0x00, 0x00, // Offsets[6]
		0x08, 0x00, 0x00, 0x00, // Offsets[7]
		0x09, 0x00, 0x00, 0x00, // Offsets[8]
		0x0a, 0x00, 0x00, 0x00, // Offsets[9]
		0x0b, 0x00, 0x00, 0x00, // Offsets[10]
		0x0c, 0x00, 0x00, 0x00, // Offsets[11]
		0x0d, 0x00, 0x00, 0x00, // Offsets[12]
		0x0e, 0x00, 0x00, 0x00, // Offsets[13]
		0x0f, 0x00, 0x00, 0x00, // Offsets[14]
		0x10, 0x00, 0x00, 0x00, // Offsets[15]
	})

	offs, err := paa.DecodeTaggOFFS(src)
	require.NoError(t, err)
	require.NotNil(t, offs)
	for i := range 16 {
		assert.Equal(t, uint32(i+1), offs.Offsets[i])
	}
}
