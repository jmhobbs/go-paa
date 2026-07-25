package formats_test

import (
	"testing"

	"github.com/jmhobbs/go-paa/formats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// synthetic smoke test until we get a real sample file
func Test_DecodeRGBA4_Basic(t *testing.T) {
	rgba4 := []byte{
		0x01, 0x23,
		0x32, 0x10,
		0x11, 0x11,
		0x00, 0x00,
	}

	rgba8, err := formats.DecodeRGBA4(rgba4, 2, 2)
	require.NoError(t, err)
	require.Equal(t, 16, len(rgba8))

	assert.Equal(
		t,
		[]uint8{
			0x00, 0x01, 0x02, 0x03,
			0x03, 0x02, 0x01, 0x00,
			0x01, 0x01, 0x01, 0x01,
			0x00, 0x00, 0x00, 0x00,
		},
		rgba8,
	)
}
