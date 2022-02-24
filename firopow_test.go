package main

import (
	"encoding/hex"
	"math"
	"math/big"
	"strings"
	"testing"
)

func MustDecodeHex(inp string) []byte {
	inp = strings.Replace(inp, "0x", "", -1)
	out, err := hex.DecodeString(inp)
	if err != nil {
		panic(err)
	}

	return out
}

func TestVerify(t *testing.T) {
	testBlocks := []struct {
		Header     []byte
		Nonce      uint64
		Height     uint64
		Difficulty float64
	}{
		// {
		// 	Header:     "5a085fb8be7e0f10cbeb45a1deda25abfef270e12a203c7dfb020aac0723fa7c",
		// 	Nonce:      0xf3e95657f2470e38,
		// 	Height:     265000,
		// 	Difficulty: 5676.51423654,
		// },
		{
			Header:     MustDecodeHex("2c128024a0274ec45f773fa878e0f9efc309ebc4864e63346931fb0a80ec9f1e"),
			Nonce:      0xf63c14518f7a9067,
			Height:     48653,
			Difficulty: 58493.72502553,
		},
	}

	for _, v := range testBlocks {
		sum, err := HashSum([]byte(v.Header), v.Nonce, v.Height)
		if err != nil {
			t.Errorf("Unable to get the hashsum of the block %v", err)
		}
		out := math.Pow(2, 256)
		div := (out / v.Difficulty)
		target := big.NewInt(int64(div))
		res, err := Verify(target, sum)
		if !res || err != nil {
			t.Errorf("Unable to verify block %v", err)
		}
	}
}
