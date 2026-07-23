//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"image/png"
	"syscall/js"

	"github.com/jmhobbs/go-paa"
)

func main() {
	js.Global().Set("paaToPng", js.FuncOf(paaToPng))
	select {}
}

func paaToPng(_ js.Value, args []js.Value) (result any) {
	defer func() {
		if r := recover(); r != nil {
			result = errorResult(fmt.Sprintf("panic: %v", r))
		}
	}()

	data := make([]byte, args[0].Get("length").Int())
	js.CopyBytesToGo(data, args[0])

	src := bytes.NewReader(data)

	paaFile, err := paa.Decode(src)
	if err != nil {
		return errorResult(fmt.Sprintf("failed to decode PAA file: %v", err))
	}

	img, err := paaFile.Mipmaps[0].Image(src)
	if err != nil {
		return errorResult(fmt.Sprintf("failed to decode mipmap image: %v", err))
	}

	var buf bytes.Buffer

	err = png.Encode(&buf, img)
	if err != nil {
		return errorResult(fmt.Sprintf("failed to encode PNG: %v", err))
	}

	jsData := js.Global().Get("Uint8Array").New(buf.Len())
	js.CopyBytesToJS(jsData, buf.Bytes())
	return js.ValueOf(map[string]any{
		"ok":  true,
		"png": jsData,
	})
}

func errorResult(message string) js.Value {
	return js.ValueOf(map[string]any{
		"ok":    false,
		"error": message,
	})
}
