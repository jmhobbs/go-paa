package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jmhobbs/go-paa"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: paa-inspect <input-file>")
	}
	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Fprintln(os.Stderr, "error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

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

	fmt.Printf("Type: %s\n", paaImg.Type)
	if paaImg.AVGC != nil {
		fmt.Printf("AVGC: %s\n", paaImg.AVGC)
	}
	if paaImg.SWIZ != nil {
		fmt.Printf("SWIZ: %s\n", paaImg.SWIZ)
	}
	fmt.Println("")
	fmt.Println("[Mipmaps]")
	for i, m := range paaImg.Mipmaps {
		fmt.Printf("  %d: %dx%d, size=%d, compressed=%v\n", i, m.Width, m.Height, m.Size, m.Compressed)
	}
}
