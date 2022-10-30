package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

func ErrHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func GetNowUnixTimestamp() int {
	return int(time.Now().Unix())
}

func HashObject(i interface{}) string {
	bytesFromObject, err := ObjectToBytes(i)
	ErrHandler(err)
	return fmt.Sprintf("%x", sha256.Sum256(bytesFromObject))
}
