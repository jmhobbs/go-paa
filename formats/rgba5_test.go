package formats_test

import (
	"testing"

	"github.com/jmhobbs/go-paa/formats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// synthetic smoke test until we get a real sample file
func Test_DecodeRGBA5_Basic(t *testing.T) {
	rgba5 := []byte{
		0b00001000, 0b10000110,
		0b00100001, 0b01001101,
	}

	rgba8, err := formats.DecodeRGBA5(rgba5, 2, 1)
	require.NoError(t, err)
	require.Equal(t, 8, len(rgba8))

	assert.Equal(
		t,
		[]uint8{
			0x01, 0x02, 0x03, 0x00,
			0x04, 0x05, 0x06, 0xFF,
		},
		rgba8,
	)
}
