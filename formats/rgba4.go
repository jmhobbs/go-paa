package formats

// Decode RGBA 4:4:4:4 into RGBA 8:8:8:8
func DecodeRGBA4(data []byte, width, height uint) ([]byte, error) {
	rgba := make([]byte, width*height*4)
	for y := range height {
		inputOffset := y * width * 2
		outputOffset := y * width * 4
		for x := range width {
			rgba[outputOffset+x*4+0] = data[inputOffset+x*2+0] >> 4 & 0x0F
			rgba[outputOffset+x*4+1] = data[inputOffset+x*2+0] & 0x0F
			rgba[outputOffset+x*4+2] = data[inputOffset+x*2+1] >> 4 & 0x0F
			rgba[outputOffset+x*4+3] = data[inputOffset+x*2+1] & 0x0F
		}
	}
	return rgba, nil
}
