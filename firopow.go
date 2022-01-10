package firopow

import (
	"math/big"

	"firo/firopow-go/progpow"
)

type Block struct {
	Target     string
	Header     string
	Nonce      uint64
	Difficulty uint64
	PrevHash   string
	Height     uint64
	SeedHash   [25]uint32
}

func newBlock() Block {
	return Block{
		Target:   " ",
		Header:   " ",
		Nonce:    0x85f22c9b3cd2f123,
		PrevHash: " ",
		Height:   1,
		SeedHash: "",
	}
}
func Sum(b Block) (*big.Int, error) {
	hash_byte := progpow.Hash(b.Height, b.SeedHash, uint64(16*1024))
	final_hash := progpow.Final_hash(b.SeedHash, hash_byte)
	big.Int(final_hash)
}

func Verify(block Block) (bool, error) {}
