package utils

import (
	"encoding/binary"
)

func Convertuint32ArrTobyte(arr []uint32) []byte {
	buf := make([]byte, len(arr)*4)

	for i, v := range arr {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}

	return buf
}

func Uint32ArrayToBytesLE(arr []uint32) []byte {
	buf := make([]byte, len(arr)*4)

	for i, v := range arr {
		binary.BigEndian.PutUint32(buf[i*4:], v)
	}

	return buf
}
