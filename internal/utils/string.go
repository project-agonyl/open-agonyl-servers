package utils

import (
	"bytes"
	"math/rand"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func ReadStringFromBytes(buffer []byte) string {
	nullIndex := bytes.IndexByte(buffer, 0)
	if nullIndex == -1 {
		return string(buffer)
	}

	return string(buffer[:nullIndex])
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
