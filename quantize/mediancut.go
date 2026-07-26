package quantize

import "slices"

func MedianCut(buf [][3]byte) [3]byte {
	if len(buf) == 0 {
		return [3]byte{0, 0, 0}
	}
	if len(buf) == 1 {
		return buf[0]
	}

	var (
		rmin uint8 = 255
		bmin uint8 = 255
		gmin uint8 = 255
		rmax uint8 = 0
		bmax uint8 = 0
		gmax uint8 = 0
	)
	for i := range len(buf) {
		rmin = min(rmin, buf[i][0])
		rmax = max(rmax, buf[i][0])
		gmin = min(gmin, buf[i][1])
		gmax = max(gmax, buf[i][1])
		bmin = min(bmin, buf[i][2])
		bmax = max(bmax, buf[i][2])
	}
	rrange := rmax - rmin
	bbrange := bmax - bmin
	grange := gmax - gmin

	var index int
	if rrange >= bbrange && rrange >= grange {
		index = 0
	} else if bbrange >= rrange && bbrange >= grange {
		index = 2
	} else {
		index = 1
	}

	slices.SortFunc(buf, comparator(index))

	return MedianCut(buf[:len(buf)/2])
}

func comparator(index int) func(a, b [3]byte) int {
	return func(a, b [3]byte) int {
		return int(b[index]) - int(a[index])
	}
}
