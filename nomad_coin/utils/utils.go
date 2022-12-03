package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var LogFn = log.Panic

func ErrHandler(err error) {
	if err != nil {
		LogFn(err)
	}
}

var timeNowFn = time.Now

func GetNowUnixTimestamp() int {
	return int(timeNowFn().Unix())
}

func HashObject(i interface{}) string {
	bytesFromObject, err := ObjectToBytes(i)
	ErrHandler(err)
	return fmt.Sprintf("%x", sha256.Sum256(bytesFromObject))
}

func Splitter(s, sep string, index int) string {
	s_slice := strings.Split(s, sep)
	if index >= len(s_slice) {
		return ""
	}
	return s_slice[index]
}

func ToJson(i interface{}) []byte {
	encodedBytes, err := json.Marshal(&i)
	ErrHandler(err)
	return encodedBytes
}
