package util

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

func DecodeGBK(text string) (string, error) {
	result := make([]byte, len(text)*2)
	decoder := simplifiedchinese.GB18030.NewDecoder()
	transformsize, _, err := decoder.Transform(result, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(result[:transformsize]), nil
}
