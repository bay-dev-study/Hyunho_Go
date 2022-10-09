package dictionary

import (
	"errors"
	"fmt"
)

type Dictionary map[string]string

var (
	ErrorKeyNotFound   error = errors.New("key not found")
	ErrorAlreadyExists error = errors.New("key already exists")
)

func (dct *Dictionary) Search(key string) (string, error) {
	val, exists := (*dct)[key]
	if !exists {
		return "", ErrorKeyNotFound
	}
	return val, nil
}

func (dct *Dictionary) Add(key, val string) error {
	_, exists := (*dct)[key]
	if exists {
		return ErrorAlreadyExists
	}
	(*dct)[key] = val
	return nil
}

func (dct *Dictionary) Update(key, val string) error {
	_, exists := (*dct)[key]
	if !exists {
		return ErrorKeyNotFound
	}
	(*dct)[key] = val
	return nil
}

func (dct *Dictionary) Delete(key string) {
	delete(*dct, key)
}

func (dct *Dictionary) Print() {
	fmt.Println("--------------------------")
	for key, val := range *dct {
		fmt.Printf("%s : %s\n", key, val)
	}
	fmt.Println("--------------------------")
}
