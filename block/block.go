package block

type Block struct {
	Header     string
	SeedHash   string
	Height     uint64
	Nonce      uint64
	Difficulty string
}

func AddBlock(block Block) bool {
	if block.Header == " " {
		return false
	}
	return true
}
