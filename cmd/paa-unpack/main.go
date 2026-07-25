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
		all            = flag.Bool("all", false, "unpack all mipmaps (default: false)")
	)
	flag.Parse()

	in, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to open input file:", err)
		os.Exit(1)
	}
	defer func() {
		if err := in.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to close input file:", err)
		}
	}()

	paaImg, err := paa.Decode(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to decode input file:", err)
		os.Exit(1)
	}

	if outputFilename == nil || *outputFilename == "" {
		*outputFilename = strings.TrimSuffix(flag.Arg(0), ".paa") + ".png"
	}

	if *all {
		for _, mipmap := range paaImg.Mipmaps {
			mipmapFilename := fmt.Sprintf("%s_(%dx%d).png", strings.TrimSuffix(*outputFilename, ".png"), mipmap.Width, mipmap.Height)
			err = writeMipmap(in, mipmap, mipmapFilename)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(1)
			}
			fmt.Println("Wrote", mipmapFilename)
		}
	} else {
		err = writeMipmap(in, paaImg.Mipmaps[0], *outputFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println("Wrote", *outputFilename)
	}
}

func writeMipmap(in *os.File, mipmap paa.Mipmap, outputFilename string) error {
	rgba, err := mipmap.Image(in)
	if err != nil {
		return fmt.Errorf("error: failed to decode mipmap image: %w", err)
	}

	sink, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("error: failed to open output file: %w", err)
	}
	defer func() {
		if err := sink.Close(); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to close output file:", err)
		}
	}()

	err = png.Encode(sink, rgba)
	if err != nil {
		return fmt.Errorf("error: failed to encode output png: %w", err)
	}

	return nil
}
