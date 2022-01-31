package main

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"firo/firopow-go/ethash"
	"firo/firopow-go/progpow"
)

type Block struct {
	Target     string
	Header     []byte
	Nonce      uint64
	Height     uint64
	Difficulty float64
}

func HashSum(b Block) (*big.Int, error) {
	var seedHead uint64 = 2048
	var seed [25]uint32
	cache := make([]uint32, 4056*4)

	seed, _ = progpow.Hash_seed(b.Header, b.Nonce)
	ethash.GenerateCache(cache, 3, seed)
	ethash.CacheSize(b.Nonce)

	epoch := ethash.CalcEpoch(b.Height)
	datasetSize := ethash.CalculateDatasetSize(epoch)

	look := func(data []uint32, index uint32) progpow.LookupFunc {
		keccak512hasher := ethash.NewKeccak512hasher()
		lookup := func(data []uint32, index uint32) []uint32 {
			return ethash.GenerateDataset(data, uint32(len(data)), index, keccak512hasher)
		}
		return lookup
	}

	l1 := make([]uint32, 4096*4)
	ethash.GenerateL1Cache(l1, cache)
	mix_hash := progpow.Hash_mix(b.Height, seedHead, datasetSize, look(cache, 4), l1)
	final_hash := progpow.Hash_final(seed, mix_hash)
	final_int := binary.BigEndian.Uint64(final_hash)
	return big.NewInt(int64(final_int)), nil
}

//func Verify(Difficulty float64,  hash *big.Int) (bool, error) {}

func main() {
	block := Block{
		Header: []byte("1fc36bd5d1bff8d134e24a997cfa43cbb2a0b956379bdc0c8df444f2553f6b7d"),
		Nonce:  1300,
		Height: 1300,
	}
	ans, _ := HashSum(block)
	fmt.Println(ans)
}
