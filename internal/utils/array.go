package utils

import "math/rand"

func GetRandomElement[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}
