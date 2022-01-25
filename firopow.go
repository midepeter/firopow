package main

import (
	"encoding/binary"
	"math/big"

	"firo/firopow-go/ethash"
	"firo/firopow-go/progpow"
)

type Block struct {
	Target     string
	Header     string
	Nonce      uint64
	PrevHash   string
	Height     uint64
	Difficulty string
}

func HashSum(b Block) (*big.Int, error) {
	seed_hash := ethash.SeedHash(b.Nonce)
	var seed uint64 = 2048
	var cache []uint32
	var dataset []uint32

	ethash.GenerateCache(cache, b.Height, seed_hash)
	ethash.GenerateDataset(dataset, b.Height, cache)

	lookup := func(dst uint32) []byte {
		res := make([]byte, 1024)
		binary.LittleEndian.PutUint32(res, dataset[dst])
		return res
	}
	mix_hash := progpow.Hash_mix(b.Height, seed, 2048, lookup, dataset)
	seedhash := progpow.Hash_seed([]byte(b.Header), b.Nonce)
	final_hash := progpow.Hash_final(seedhash, mix_hash)
	final_int := binary.BigEndian.Uint64(final_hash)
	return big.NewInt(int64(final_int)), nil
}

//func Verify(block Block) (bool, error) {}
