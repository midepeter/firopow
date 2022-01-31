package main

// import (
// 	"crypto/sha256"
// 	"fmt"
// 	"testing"
// )

// func TestHashSum(t *testing.T) {
// 	testBlocks := []struct {
// 		Target     string
// 		Header     string
// 		Nonce      uint64
// 		Height     uint64
// 		Difficulty float64
// 	}{
// 		{
// 			Target:     " ",
// 			Header:     "e8d20ca775c2d32a91d097c3b57836674c2b5cf0d2b8c1652b6b3e9cdc9e6b05",
// 			Nonce:      445397,
// 			Height:     445397,
// 			Difficulty: 8352.47625956,
// 		},
// 	}

// 	for _, v := range testBlocks[:] {
// 		ans := sha256.Sum256([]byte(v.Header))
// 		fmt.Sprintln(ans)
// 		sum, err := HashSum(v)
// 		if err != nil {
// 			fmt.Println("unable to hash the sum of the bloc")
// 		}
// 		fmt.Printf("The hash sum is %v", sum)
// 	}
// }
