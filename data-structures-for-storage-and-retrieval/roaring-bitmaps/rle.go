package bitmap

func compress(b *uncompressedBitmap) []uint64 {
	var compressedData []uint64
	var curRun uint64
	var numBits int
	var prevRun uint64
	var length uint64

	for _, el := range b.data {
		for j := 0; j < wordSize; j++ {
			if numBits == wordSize-2 {
				numBits = 0

				if curRun == 0 || curRun == 1 {
					switch {
					case length == 0:
						length = 1
						prevRun = curRun
					case length > 0 && curRun == prevRun:
						length++
					case length > 0 && curRun != prevRun:
						fillData := length
						fillData |= (1 << 63)
						fillData |= (prevRun << 62)

						compressedData = append(compressedData, fillData)
						length = 0
					}
				} else {
					if length > 0 {
						fillData := length
						fillData |= (1 << 63)
						fillData |= (prevRun << 62)

						compressedData = append(compressedData, fillData)
					}

					compressedData = append(compressedData, curRun)
					length = 0
				}

				curRun = 0
				continue
			}

			bit := el & (1 << j)
			curRun = ((2 << 63) - 1) & bit
			numBits++
		}
	}

	// fmt.Printf("% b", compressedData)

	return compressedData
}

func decompress(compressed []uint64) *uncompressedBitmap {
	var data []uint64
	var uncompressedEl uint64
	i := wordSize - 1

	for _, el := range compressed {
		if (1<<63)&el > 0 {
			bitValue := (1 << 62) & el
			length := ((2 << 62) - 1) & el

			for j := uint64(0); j < length; j++ {
				for k := 61; k >= 0; k-- {
					uncompressedEl |= (bitValue << i)
					i--

					if i < 0 {
						data = append(data, uncompressedEl)
						uncompressedEl = 0
						i = wordSize - 1
					}
				}
			}

			// fmt.Printf("Fill element: %064b\n", el)
		} else {
			// fmt.Printf("Literal element: %064b\n", el)
			for j := 62; j >= 0; j-- {
				bitValue := (1 << j) & el
				uncompressedEl |= (bitValue << i)
				i--

				if i < 0 {
					data = append(data, uncompressedEl)
					uncompressedEl = 0
					i = wordSize - 1
				}
			}
		}
	}

	return &uncompressedBitmap{
		data: data,
	}
}
