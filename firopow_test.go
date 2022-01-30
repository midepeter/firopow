package main

// import (
// 	"fmt"
// 	"testing"
// )

// func TestHashSum(t *testing.T) {
// 	testBlocks := []struct {
// 		Target     string
// 		Header     string
// 		Nonce      uint64
// 		PrevHash   string
// 		Height     uint64
// 		Difficulty string
// 	}{
// 		{
// 			Target:     " ",
// 			Header:     "e8d20ca775c2d32a91d097c3b57836674c2b5cf0d2b8c1652b6b3e9cdc9e6b05",
// 			Nonce:      445397,
// 			PrevHash:   "6c6e049e3387948c29eb9c08121b13d040cd5c059f8563d191d0459f3fa34875",
// 			Height:     445397,
// 			Difficulty: "8352.47625956",
// 		},
// 	}

// 	for _, v := range testBlocks {
// 		sum, err := HashSum(v)
// 		if err != nil {
// 			fmt.Println("unable to hash the sum of the bloc")
// 		}
// 		fmt.Printf("The hash sum is %v", sum)
// 	}
// }
