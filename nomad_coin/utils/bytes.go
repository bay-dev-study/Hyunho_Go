package utils

import (
	"bytes"
	"encoding/gob"
)

func ObjectToBytes(i interface{}) ([]byte, error) {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	err := encoder.Encode(i)
	return aBuffer.Bytes(), err
}

func ObjectFromBytes(i interface{}, data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(i)
}
