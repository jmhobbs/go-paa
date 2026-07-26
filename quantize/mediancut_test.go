package quantize_test

import (
	"testing"

	"github.com/jmhobbs/go-paa/quantize"

	"github.com/stretchr/testify/assert"
)

func Test_MedianCut_ZeroElements(t *testing.T) {
	img := [][3]byte{}

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{0, 0, 0}, result)
}

func Test_MedianCut_SingleElement(t *testing.T) {
	img := [][3]byte{{12, 34, 56}}

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{12, 34, 56}, result)
}

func Test_MedianCut_TwoElements_RedDominant(t *testing.T) {
	img := [][3]byte{
		{50, 10, 10},
		{200, 10, 10},
	}

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{200, 10, 10}, result)
}

func Test_MedianCut_TieBreaksToRed(t *testing.T) {
	img := [][3]byte{
		{0, 0, 0},
		{100, 100, 100},
	}

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{100, 100, 100}, result)
}

func Test_MedianCut_BlueThenGreenDominant(t *testing.T) {
	img := [][3]byte{
		{10, 10, 250}, // color0
		{10, 10, 5},   // color1
		{5, 200, 100}, // color2
		{8, 5, 90},    // color3
	}

	// First pass: rrange=5 (5-10), grange=195 (5-200), bbrange=245 (5-250)
	// Sorted by blue; color0(250), color2(100), color3(90), color1(5)
	// Half: [color0, color2]
	//
	// Second pass: rrange=5, grange=190, bbrange=150
	// Sorted by green: color2(200), color0(10)
	// Half: [color2]

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{5, 200, 100}, result)
}

func Test_MedianCut_GreenDominantAcrossTwoPasses(t *testing.T) {
	img := [][3]byte{
		{0, 255, 0}, // color0
		{0, 0, 0},   // color1
		{64, 1, 0},  // color2
		{0, 2, 128}, // color3
	}

	// First pass: rrange=64, grange=255, bbrange=128
	// Sorted by green: color0(255), color3(2), color2(1), color1(0)
	// Half: [color0, color3]
	//
	// Second pass: rrange=0, grange=253, bbrange=128
	// Sorted by green: color0(255), color3(2)
	// Half: [color0]

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{0, 255, 0}, result)
}

func Test_MedianCut_OddLengthSliceTruncatesDown(t *testing.T) {
	img := [][3]byte{
		{50, 0, 0},
		{10, 0, 0},
		{90, 0, 0},
		{30, 0, 0},
		{70, 0, 0},
	}

	// First pass: sorted by red: 90, 70, 50, 30, 10
	// Half: [90, 70] (5/2 truncates to 2, the rest are dropped)
	//
	// Second pass: sorted by red: 90, 70
	// Half: [90]

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{90, 0, 0}, result)
}

func Test_MedianCut_IdenticalColors(t *testing.T) {
	img := [][3]byte{
		{77, 77, 77},
		{77, 77, 77},
		{77, 77, 77},
	}

	result := quantize.MedianCut(img)
	assert.Equal(t, [3]byte{77, 77, 77}, result)
}
