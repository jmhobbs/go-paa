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

# References

https://community.bistudio.com/wiki/PAA_File_Format
