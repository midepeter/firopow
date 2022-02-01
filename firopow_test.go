package main

import (
	"math"
	"math/big"
	"testing"
)

func TestVerify(t *testing.T) {
	testBlocks := []struct {
		Header     string
		Nonce      uint64
		Height     uint64
		Difficulty float64
	}{
		{
			Header:     "1d695cf4cd0eee3dcb981486238f900f13efad15793cc35b326eaa33579bbd06",
			Nonce:      447435,
			Height:     447435,
			Difficulty: 8125.83461514,
		},
		{
			Header:     "7ae08b7120eac2738e2279c927cf0db9909505d158b13a2f18640abdfb640beb",
			Nonce:      447433,
			Height:     447433,
			Difficulty: 8125.83461514,
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
