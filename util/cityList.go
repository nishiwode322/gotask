package util

type CityList []string

func (c CityList) Len() int {
	return len(c)
}

func (c CityList) Less(i, j int) bool {
	first, _ := DecodeUTF(c[i])
	second, _ := DecodeUTF(c[j])
	sLen := len(second)
	for index, character := range first {
		if index > sLen-1 {
			return false
		}
		if character != second[index] {
			return character < second[index]
		}
	}
	return true
}

func (c CityList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
