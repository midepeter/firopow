package firopow

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
	}
}

//func Sum(b Block) (*big.Int, error) {}

//func Verify(block Block) (bool, error) {}
