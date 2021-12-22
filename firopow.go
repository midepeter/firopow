package firopow

import (
	"math/big"
)

type Block struct {
	Target     string
	Header     string
	Nonce      uint64
	Difficulty uint64
	PrevHash   string
	Height     uint64
	SeedHash   string
}

func Sum(block Block) (*big.Int, error) {}

func Verify(block Block) (bool, error) {}
