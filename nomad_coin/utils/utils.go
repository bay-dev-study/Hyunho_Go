package utils

import "log"

func ErrHandler(err error) {
	if err != nil {
		log.Panic(err)
	}
}
