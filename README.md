[![Go Reference](https://pkg.go.dev/badge/github.com/jmhobbs/go-paa.svg)](https://pkg.go.dev/github.com/jmhobbs/go-paa)
[![Lint & Test](https://github.com/jmhobbs/go-paa/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/jmhobbs/go-paa/actions/workflows/lint-and-test.yml)
[![codecov](https://codecov.io/github/jmhobbs/go-paa/graph/badge.svg?token=sB2axgNro5)](https://codecov.io/github/jmhobbs/go-paa)

# go-paa

A (WIP) Go library for working with PAA files from Bohemia Interactive.

Currently, it can only read DXT1, DXT3 and DXT5 format PAA files.

## Usage

### Web

A PAA to PNG conversion tool is available in your browser at [tools.dzhosts.com](https://tools.dzhosts.com/paa-to-png/). This is the easiest way to use it, and it should work on any platform.

### paa-unpack

The `paa-unpack` command line tool with unpack the highest resolution image from the PAA file into a PNG file.

```bash
$ paa-unpack -h
usage: paa-unpack [options] <input-file>

options:
  -output string
        output filename (default: <input-file>.png)

$ paa-unpack testdata/test-pattern.paa
Converted testdata/test-pattern.paa to testdata/test-pattern.png
```

### paa-inspect

Dump PAA information.

```bash
$ paa-inspect -h
usage: paa-inspect <input-file>

$ paa-inspect testdata/test-pattern.paa
Type: DXT1
AVGC: R=182, G=170, B=169, A=255 <#B6AAA9FF>

[Mipmaps]
  0: 512x256, size=639, compressed=true
  1: 256x128, size=373, compressed=true
  2: 128x64, size=4096, compressed=false
  3: 64x32, size=1024, compressed=false
  4: 32x16, size=256, compressed=false
  5: 16x8, size=64, compressed=false
  6: 8x4, size=16, compressed=false
```

# References

https://community.bistudio.com/wiki/PAA_File_Format
