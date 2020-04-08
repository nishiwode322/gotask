package util

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func DecodeGBK(text string) (string, error) {
	result := make([]byte, len(text)*2)
	decoder := simplifiedchinese.GB18030.NewDecoder()
	transformSize, _, err := decoder.Transform(result, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(result[:transformSize]), nil
}

func DecodeUTF(text string) ([]byte, error) {
	gb := simplifiedchinese.All[0]
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(text)), gb.NewEncoder()))
}
