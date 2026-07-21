package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/jmhobbs/go-paa"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: paa-unpack [options] <input-file>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "options:")
		flag.PrintDefaults()
	}

	var (
		outputFilename = flag.String("output", "", "output filename (default: <input-file>.png)")
	)
	flag.Parse()

	in, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to open input file:", err)
		os.Exit(1)
	}
	defer in.Close()

	paaImg, err := paa.Decode(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to decode input file:", err)
		os.Exit(1)
	}

	if outputFilename == nil || *outputFilename == "" {
		*outputFilename = strings.TrimSuffix(flag.Arg(0), ".paa") + ".png"
	}

	rgba, err := paaImg.Mipmaps[0].Image()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to decode mipmap image:", err)
		os.Exit(1)
	}

	sink, err := os.Create(*outputFilename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to open output file:", err)
		os.Exit(1)
	}
	defer sink.Close()

	err = png.Encode(sink, rgba)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to encode output png:", err)
		os.Exit(1)
	}

	fmt.Println("Converted", flag.Arg(0), "to", *outputFilename)
}
