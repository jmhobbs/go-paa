package formats

// Decode RGBA 5:5:5:1 into RGBA 8:8:8:8
func DecodeRGBA5(data []byte, width, height uint) ([]byte, error) {
	rgba := make([]byte, width*height*4)
	for y := range height {
		inputOffset := y * width * 2
		outputOffset := y * width * 4
		for x := range width {
			rgba[outputOffset+x*4+0] = (data[inputOffset+x*2+0] & 0x7C) >> 3
			rgba[outputOffset+x*4+1] = (data[inputOffset+x*2+0]&0x07)<<2 | (data[inputOffset+x*2+1]&0xE0)>>6
			rgba[outputOffset+x*4+2] = (data[inputOffset+x*2+1] & 0x3E) >> 1
			if (data[inputOffset+x*2+1] & 0x01) == 0x01 {
				rgba[outputOffset+x*4+3] = 0xFF
			} else {
				rgba[outputOffset+x*4+3] = 0x00
			}
		}
	}
	return rgba, nil
}
