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
	PrevHash   string
	Height     uint64
	Difficulty string
}

func HashSum(b Block) (*big.Int, error) {
	var seedHead uint64 = 2048
	var seed [25]uint32
	cache := make([]uint32, 1024/4)
	var dataset []uint32

	seed, _ = progpow.Hash_seed(b.Header, b.Nonce)
	fmt.Println(seed)
	ethash.GenerateCache(cache, b.Height, seed)
	ethash.CacheSize(b.Nonce)
	epoch := ethash.CalcEpoch(b.Height)
	datasetSize := ethash.CalculateDatasetSize(epoch)

	look := func(data []uint32, index uint32) progpow.LookupFunc {
		//var h hash.Hash
		keccak512hasher := ethash.NewKeccak512hasher()
		lookup := func(data []uint32, index uint32) []uint32 {
			return ethash.GenerateDataset(data, uint32(len(data)), index, keccak512hasher)
		}
		dataset = lookup(cache, 5)
		fmt.Println(dataset)
		return lookup
	}

	l1 := make([]uint32, 1280)
	ethash.GenerateL1Cache(l1, cache)
	fmt.Println("iteml1", l1)
	mix_hash := progpow.Hash_mix(b.Height, seedHead, datasetSize, look(cache, uint32(100)), l1)
	final_hash := progpow.Hash_final(seed, mix_hash)
	final_int := binary.BigEndian.Uint64(final_hash)
	fmt.Println(final_hash)
	return big.NewInt(int64(final_int)), nil
}

func main() {
	block := Block{
		Header: []byte("c56347a2929c51d721584bfa18f2a88ddd56f6c30b52cef4b1b14ce7f36e54e4"),
		Nonce:  446508,
		Height: 446508,
	}
	ans, _ := HashSum(block)
	fmt.Println(ans)
}

//func Verify(block Block) (bool, error) {}
