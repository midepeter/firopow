package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
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

func HashSum(Header []byte, Nonce uint64, Height uint64) (*big.Int, error) {
	var seed [25]uint32

	seed, SeedHash := progpow.Hash_seed(Header, Nonce)

	epoch := ethash.CalcEpoch(Height)
	size := ethash.CacheSize(epoch)
	fmt.Println("Cache size: ", size)
	cache := make([]uint32, size/4)

	seedByte := ethash.SeedHash(Height)

	ethash.GenerateCache(cache, epoch, seedByte)
	datasetSize := ethash.DatasetSize(epoch)

	look := func(data []uint32, index uint32) progpow.LookupFunc {
		keccak512hasher := ethash.NewKeccak512hasher()
		lookup := func(data []uint32, index uint32) []uint32 {
			return ethash.GenerateDataset(data, uint32(len(data)), index, keccak512hasher)
		}
		return lookup
	}

	l1 := make([]uint32, 4096*4)
	ethash.GenerateL1Cache(l1, cache)
	fmt.Println("Dataset size", datasetSize)
	mix_hash := progpow.Hash_mix(Height, SeedHash, datasetSize, look(cache, uint32(epoch)), l1)
	fmt.Println("The value of the mix hash is", hex.EncodeToString(mix_hash))
	final_hash := progpow.Hash_final(seed, mix_hash)
	fmt.Println("The encoded value of the hash is ", hex.EncodeToString(final_hash))
	final_int := binary.BigEndian.Uint64(final_hash)
	return big.NewInt(int64(final_int)), nil
}

func Verify(target *big.Int, hashSum *big.Int) (bool, error) {
	r := hashSum.CmpAbs(target)
	if r == -1 {
		return true, nil
	}
	return false, errors.New("could not verify block successfully")
}
