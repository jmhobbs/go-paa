package paa_test

import (
	"os"
	"testing"

	"github.com/jmhobbs/go-paa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Decode(t *testing.T) {
	f, err := os.Open("testdata/test-pattern.paa")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	img, err := paa.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if img.Type != paa.Type_DXT1 {
		t.Error("unexpected type:", img.Type)
	}

	require.NotNil(t, img.AVGC)
	assert.Equal(t, uint8(0xb6), img.AVGC.Red)
	assert.Equal(t, uint8(0xaa), img.AVGC.Green)
	assert.Equal(t, uint8(0xa9), img.AVGC.Blue)
	assert.Equal(t, uint8(0xff), img.AVGC.Alpha)
}
