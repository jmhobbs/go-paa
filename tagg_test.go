package paa_test

import (
	"bytes"
	"image"
	"io"
	"testing"

	"github.com/jmhobbs/go-paa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EncodeTaggAVGC(t *testing.T) {
	var buf bytes.Buffer
	err := paa.TaggAVGC{Red: 0x12, Green: 0x34, Blue: 0x56, Alpha: 0x78}.Write(&buf)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]byte{
			0x43, 0x47, 0x56, 0x41, // header
			0x04, 0x00, 0x00, 0x00, // length
			0x12, 0x34, 0x56, 0x78, // data
		},
		buf.Bytes(),
	)
}

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

func Test_EncodeTaggMAXC(t *testing.T) {
	var buf bytes.Buffer
	err := paa.TaggMAXC{Data: [4]uint8{0x12, 0x23, 0x34, 0x45}}.Write(&buf)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]byte{
			0x43, 0x58, 0x41, 0x4d, // header
			0x04, 0x00, 0x00, 0x00, // length
			0x12, 0x23, 0x34, 0x45, // data
		},
		buf.Bytes(),
	)
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
		"SWIZ": func(r io.Reader) (any, error) { return paa.DecodeTaggSWIZ(r) },
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
		{"SWIZ", func(r io.Reader) (any, error) { return paa.DecodeTaggSWIZ(r) }},
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

func Test_EncodeTaggOFFS(t *testing.T) {
	var buf bytes.Buffer
	var offsets [16]uint32 = [16]uint32{1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	err := paa.TaggOFFS{Offsets: offsets}.Write(&buf)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]byte{
			0x53, 0x46, 0x46, 0x4f, // header
			0x40, 0x00, 0x00, 0x00, // length: 64 (16*4)
			0x01, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x00, 0x00,
			0x03, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00,
			0x05, 0x00, 0x00, 0x00,
			0x06, 0x00, 0x00, 0x00,
			0x07, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00,
		},
		buf.Bytes(),
	)
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

func Test_EncodeTaggSWIZ(t *testing.T) {
	var buf bytes.Buffer
	err := paa.TaggSWIZ{Alpha: 0x05, Red: 0x04, Green: 0x02, Blue: 0x03}.Write(&buf)
	require.NoError(t, err)
	assert.Equal(
		t,
		[]byte{
			0x5a, 0x49, 0x57, 0x53, // header
			0x04, 0x00, 0x00, 0x00, // length
			0x05, 0x04, 0x02, 0x03, // data
		},
		buf.Bytes(),
	)
}

func Test_DecodeTaggSWIZ(t *testing.T) {
	// alpha and red are swapped and negated, green and blue are stored unchanged.
	src := bytes.NewReader([]byte{
		0x04, 0x00, 0x00, 0x00, // Length: 4
		0x05, // Alpha
		0x04, // Red
		0x02, // Green
		0x03, // Blue
	})

	swiz, err := paa.DecodeTaggSWIZ(src)
	require.NoError(t, err)
	require.NotNil(t, swiz)
	assert.Equal(t, uint8(0x05), swiz.Alpha)
	assert.Equal(t, uint8(0x04), swiz.Red)
	assert.Equal(t, uint8(0x02), swiz.Green)
	assert.Equal(t, uint8(0x03), swiz.Blue)

	assert.Equal(t, paa.SwizzleChannel{Destination: paa.Red, Negated: true}, swiz.AlphaChannel())
	assert.Equal(t, paa.SwizzleChannel{Destination: paa.Alpha, Negated: true}, swiz.RedChannel())
	assert.Equal(t, paa.SwizzleChannel{Destination: paa.Green, Negated: false}, swiz.GreenChannel())
	assert.Equal(t, paa.SwizzleChannel{Destination: paa.Blue, Negated: false}, swiz.BlueChannel())
}

func Test_SwizzleChannel_ForceFF(t *testing.T) {
	// bit 3 set forces the channel's data to 0xff, regardless of the other bits.
	channel := paa.DecodeSwizzleChannel(0x0d)
	assert.True(t, channel.ForceFF)
	assert.Equal(t, "forced to 0xff", channel.String())
}

func Test_TaggSWIZ_Unswizzle(t *testing.T) {
	// Alpha and Red are swapped and negated,
	// Green and Blue untouched
	swiz := paa.TaggSWIZ{Alpha: 0x05, Red: 0x04, Green: 0x02, Blue: 0x03}

	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.Pix = []byte{10, 20, 30, 200}

	out := swiz.Unswizzle(src)

	require.Equal(t, src.Rect, out.Rect)
	assert.Equal(
		t,
		[]byte{
			0xff - 200, // Red = negated stored Alpha
			20,         // Green = stored Green, unchanged
			30,         // Blue = stored Blue, unchanged
			0xff - 10,  // Alpha = negated stored Red
		},
		out.Pix,
	)
}

func Test_TaggSWIZ_Unswizzle_ForceFF(t *testing.T) {
	// A ForceFF channel is always 0xff, regardless of what is stored.
	swiz := paa.TaggSWIZ{Alpha: 0x08, Red: 0x01, Green: 0x02, Blue: 0x03}

	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.Pix = []byte{10, 20, 30, 200}

	out := swiz.Unswizzle(src)

	assert.Equal(t, []byte{10, 20, 30, 0xff}, out.Pix)
}

func Test_TaggSWIZ_Unswizzle_Unchanged(t *testing.T) {
	swiz := paa.TaggSWIZ{Alpha: 0x00, Red: 0x01, Green: 0x02, Blue: 0x03}

	src := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	src.Pix = []byte{10, 20, 30, 200}

	out := swiz.Unswizzle(src)

	assert.Equal(t, src.Pix, out.Pix)
}
