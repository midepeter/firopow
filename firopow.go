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
	fmt.Println("Seeeeeed", seed)
	fmt.Println("seed hash", SeedHash)
	epoch := ethash.CalcEpoch(Height)
	size := ethash.CacheSize(Height)
	fmt.Println("Cache size: ", size)
	cache := make([]uint32, size/4)

	seedByte := ethash.SeedHash(38*1300 + 1)

	fmt.Println("seedByte: ", seedByte)
	ethash.GenerateCache(cache, epoch, seedByte)
	datasetSize := ethash.DatasetSize(epoch)

	lookup := func(index uint32) []uint32 {
		keccak512hasher := ethash.NewKeccak512hasher()
		return ethash.GenerateDatasetItemUint(cache, index, 4, keccak512hasher)
	}

	l1 := make([]uint32, 4096*4)
	ethash.GenerateL1Cache(l1, cache)
	fmt.Println("Dataset size", datasetSize)
	mix_hash := progpow.Hash_mix(Height, SeedHash, datasetSize, lookup, l1)
	fmt.Println("The value of the mix hash is", hex.EncodeToString(mix_hash))
	final_hash := progpow.Hash_final(seed, mix_hash)
	fmt.Println("The encoded value of the hash is ", hex.EncodeToString(final_hash))
	final_int := binary.BigEndian.Uint64(final_hash)
	fmt.Println("Cache", len(cache))
	return big.NewInt(int64(final_int)), nil
}

func Verify(target *big.Int, hashSum *big.Int) (bool, error) {
	r := hashSum.CmpAbs(target)
	if r == -1 {
		return true, nil
	}
	return false, errors.New("could not verify block successfully")
}
