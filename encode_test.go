package paa_test

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/jmhobbs/go-paa"

	"github.com/stretchr/testify/assert"
)

func Test_CalculateAVGC(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, image.Black)
	img.Set(1, 1, image.Black)
	img.Set(0, 1, image.White)
	img.Set(1, 0, image.White)

	avgc := paa.CalculateAVGC(img)
	assert.Equal(
		t,
		paa.TaggAVGC{
			Red:   0x00,
			Green: 0x00,
			Blue:  0x00,
			Alpha: 0xFF,
		},
		avgc,
	)
}

// opaqueImage wraps an image.Image without exposing its concrete type,
// forcing CalculateAVGC's type switch to fall through to the slow path
// even when the wrapped image is an *image.RGBA or *image.NRGBA.
type opaqueImage struct {
	image.Image
}

func Test_CalculateAVGC_FastPath_MultiPixel(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.SetRGBA(0, 0, color.RGBA{R: 10, G: 10, B: 250, A: 255})
	img.SetRGBA(1, 0, color.RGBA{R: 10, G: 10, B: 5, A: 255})
	img.SetRGBA(0, 1, color.RGBA{R: 5, G: 200, B: 100, A: 255})
	img.SetRGBA(1, 1, color.RGBA{R: 8, G: 5, B: 90, A: 255})

	avgc := paa.CalculateAVGC(img)
	assert.Equal(
		t,
		paa.TaggAVGC{
			Red:   0x05,
			Green: 0xC8,
			Blue:  0x64,
			Alpha: 0xFF,
		},
		avgc,
	)
}

func Test_CalculateAVGC_SlowPath_MultiPixel(t *testing.T) {
	rgba := image.NewRGBA(image.Rect(0, 0, 2, 2))
	rgba.SetRGBA(0, 0, color.RGBA{R: 10, G: 10, B: 250, A: 255})
	rgba.SetRGBA(1, 0, color.RGBA{R: 10, G: 10, B: 5, A: 255})
	rgba.SetRGBA(0, 1, color.RGBA{R: 5, G: 200, B: 100, A: 255})
	rgba.SetRGBA(1, 1, color.RGBA{R: 8, G: 5, B: 90, A: 255})

	avgc := paa.CalculateAVGC(opaqueImage{rgba})
	assert.Equal(
		t,
		paa.TaggAVGC{
			Red:   0x05,
			Green: 0xC8,
			Blue:  0x64,
			Alpha: 0xFF,
		},
		avgc,
	)
}

func newBenchmarkRGBA(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(x),
				G: uint8(y),
				B: uint8(x + y),
				A: 255,
			})
		}
	}
	return img
}

var benchmarkImageSizes = []int{16, 64, 256}

func Benchmark_CalculateAVGC_FastPath(b *testing.B) {
	for _, size := range benchmarkImageSizes {
		b.Run(fmt.Sprintf("%dx%d", size, size), func(b *testing.B) {
			img := newBenchmarkRGBA(size)

			b.ResetTimer()
			for b.Loop() {
				paa.CalculateAVGC(img)
			}
		})
	}
}

func Benchmark_CalculateAVGC_SlowPath(b *testing.B) {
	for _, size := range benchmarkImageSizes {
		b.Run(fmt.Sprintf("%dx%d", size, size), func(b *testing.B) {
			img := opaqueImage{newBenchmarkRGBA(size)}

			b.ResetTimer()
			for b.Loop() {
				paa.CalculateAVGC(img)
			}
		})
	}
}

/*
func Test_New_DefaultsToDXT1(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	var buf bytes.Buffer
	builder, err := paa.New(img)
	require.NoError(t, err)
	require.NotNil(t, builder)

	require.NoError(t, builder.Write(&buf))

	var typeOfPaX uint32
	require.NoError(t, binary.Read(&buf, binary.LittleEndian, &typeOfPaX))
	assert.Equal(t, uint32(paa.Type_DXT1), typeOfPaX)
}

func Test_New_AppliesOptions(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	var buf bytes.Buffer
	builder, err := paa.New(img, paa.WithType(paa.Type_DXT3))
	require.NoError(t, err)
	require.NotNil(t, builder)

	require.NoError(t, builder.Write(&buf))

	var typeOfPaX uint32
	require.NoError(t, binary.Read(&buf, binary.LittleEndian, &typeOfPaX))
	assert.Equal(t, uint32(paa.Type_DXT3), typeOfPaX)
}

func Test_New_UnsupportedTypeOption(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))

	builder, err := paa.New(img, paa.WithType(paa.Type_RGBA8))
	require.Error(t, err)
	assert.Nil(t, builder)
}

func Test_WithType_Supported(t *testing.T) {
	for _, typeOfPaX := range []paa.TypeOfPaX{paa.Type_DXT1, paa.Type_DXT3, paa.Type_DXT5} {
		builder, err := paa.New(image.NewRGBA(image.Rect(0, 0, 1, 1)), paa.WithType(typeOfPaX))
		require.NoError(t, err)
		require.NotNil(t, builder)
	}
}

func Test_WithType_Unsupported(t *testing.T) {
	_, err := paa.New(image.NewRGBA(image.Rect(0, 0, 1, 1)), paa.WithType(paa.Type_RGBA8))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "currently unsupported type")
}

// failAfterWriter errors starting on its (failOn+1)th Write call, letting
// earlier calls through so both binary.Write sites in Write can be tested.
type failAfterWriter struct {
	failOn int
	calls  int
}

func (w *failAfterWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls > w.failOn {
		return 0, errors.New("boom")
	}
	return len(p), nil
}

func Test_Write_PropagatesTypeWriteError(t *testing.T) {
	builder, err := paa.New(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	require.NoError(t, err)

	err = builder.Write(&failAfterWriter{failOn: 0})
	require.Error(t, err)
}

func Test_Write_PropagatesSignatureWriteError(t *testing.T) {
	builder, err := paa.New(image.NewRGBA(image.Rect(0, 0, 1, 1)))
	require.NoError(t, err)

	err = builder.Write(&failAfterWriter{failOn: 1})
	require.Error(t, err)
}

func Test_Write_EmitsTypeAndTaggSignature(t *testing.T) {
	builder, err := paa.New(image.NewRGBA(image.Rect(0, 0, 1, 1)), paa.WithType(paa.Type_DXT5))
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, builder.Write(&buf))

	var typeOfPaX uint32
	require.NoError(t, binary.Read(&buf, binary.LittleEndian, &typeOfPaX))
	assert.Equal(t, uint32(paa.Type_DXT5), typeOfPaX)

	var signature uint32
	require.NoError(t, binary.Read(&buf, binary.LittleEndian, &signature))
	assert.Equal(t, paa.TaggSignature, signature)

	assert.Zero(t, buf.Len(), "current encoder does not yet write taggs or mipmap data")
}
*/
