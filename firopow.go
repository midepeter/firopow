package main

import (
	"encoding/binary"
	"fmt"
	"hash"
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
	var cache []uint32
	var dataset []uint32

	seed, _ = progpow.Hash_seed(b.Header, b.Nonce)
	ethash.GenerateCache(cache, b.Height, seed)
	ethash.CacheSize(b.Nonce)

	look := func(data []uint32, index uint32) progpow.LookupFunc {
		var h hash.Hash
		hasher := ethash.MakeHasher(h)
		lookup := func(data []uint32, ndex uint32) []uint32 {
			return ethash.GenerateDataset(data, uint32(len(data)), 4, hasher)
		}
		return lookup
	}
	mix_hash := progpow.Hash_mix(b.Height, seedHead, uint64(len(dataset)), look(cache, 4), dataset)
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
