package utils

import "bytes"

func ReadStringFromBytes(buffer []byte) string {
	nullIndex := bytes.IndexByte(buffer, 0)
	if nullIndex == -1 {
		return string(buffer)
	}

	return string(buffer[:nullIndex])
}
