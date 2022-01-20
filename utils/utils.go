package utils

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strconv"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))

//const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

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

func Uint64ToBytesLE(val uint64) []byte {
	data := make([]byte, 8)

	binary.BigEndian.PutUint64(data, val)

	return data
}

func XORBytes(dst, a, b []byte) int {
	// if supportsUnaligned {
	return fastXORBytes(dst, a, b)
	// }
	// return safeXORBytes(dst, a, b)
}

func fastXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

/* func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
} */

func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, errors.New("error")
	}
	if !has0xPrefix(input) {
		return nil, errors.New("invalid input")
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = mapError(err)
	}
	return b, err
}

// MustDecode decodes a hex string with 0x prefix. It panics for invalid input.
func MustDecode(input string) []byte {
	dec, err := Decode(input)
	if err != nil {
		panic(err)
	}
	return dec
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return errors.New("out of range")
		case strconv.ErrSyntax:
			return errors.New("error syntax")
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return errors.New("error in syntax")
	}
	if err == hex.ErrLength {
		return errors.New("Error in length")
	}
	return err
}
