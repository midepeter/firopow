package progpow

import (
	"encoding/binary"
	"firo/firopow-go/keccak"
	"firo/firopow-go/utils"
	"math/bits"
)

var (
	PeriodLength        = 10
	DagLoads            = 4
	CacheBytes          = 16 * 1024
	LaneCount           = 16
	RegisterCount       = 32
	RoundCount          = 64
	RoundCacheAccesses  = 11
	RoundMathOperations = 18
)

const (
	fnvOffsetBasis uint32 = 0x811c9dc5
	progpowRegs    uint32 = 32
	progpowLanes   uint32 = 16
	fnvoffSetBasis uint32 = 0x811c9dc5
)

type LookupFunc func(index uint32) []uint32

func rotl32(a, b uint32) uint32 {
	return a<<(b&31) | a>>((32-b)&31)
}

func rotr32(a, b uint32) uint32 {
	return a<<((32-b)&31) | a>>(b&31)
}

func clz32(a uint32) uint32 {
	return uint32(bits.LeadingZeros32(a))
}

func popcount32(a uint32) uint32 {
	return uint32(bits.OnesCount32(a))
}

func mul_hi32(a, b uint32) uint32 {
	return uint32((uint64(a) * uint64(b)) >> 32)
}

func merge(a, b, selector uint32) uint32 {
	x := ((selector >> 16) % 31) + 1

	switch selector % 4 {
	case 0:
		return (a * 33) + b
	case 1:
		return (a ^ b) * 33
	case 2:
		return rotl32(a, x) ^ b
	case 3:
		return rotr32(a, x) ^ b
	}

	return 0
}

func math(a, b, r uint32) uint32 {
	switch r % 11 {
	case 0:
		return a + b
	case 1:
		return a * b
	case 2:
		return mul_hi32(a, b)
	case 3:
		if a > b {
			return b
		}
		return a
	case 4:
		return rotl32(a, b)
	case 5:
		return rotr32(a, b)
	case 6:
		return a & b
	case 7:
		return a | b
	case 8:
		return a ^ b
	case 9:
		return clz32(a) + clz32(b)
	case 10:
		return popcount32(a) + popcount32(b)
	}

	return 0
}

func ProgPOWInit(hash []byte, nonce uint64) ([25]uint32, uint64) {
	var seed [25]uint32
	for i := 0; i < 8; i += 1 {
		seed[i] = binary.LittleEndian.Uint32(hash[i*4 : i*4+4])
	}

	seed[8] = uint32(nonce)
	seed[9] = uint32(nonce >> 32)

	keccak.KeccakF800(&seed)

	seedHead := uint64(seed[0]) + (uint64(seed[1]) << 32)

	return seed, seedHead
}

func initMix(seed uint64, numLanes, numRegs int) [][]uint32 {
	z := Fnv1a(fnvoffSetBasis, uint32(seed))
	w := Fnv1a(z, uint32(seed>>32))

	mix := make([][]uint32, numLanes)

	for lane := range mix {
		jsr := Fnv1a(w, uint32(lane))
		jcong := Fnv1a(jsr, uint32(lane<<32))

		rng := New(z, w, jsr, jcong)

		mix[lane] = make([]uint32, numRegs)
		for reg := range mix[lane] {
			mix[lane][reg] = rng.Next()
		}
	}

	return mix
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func round(seed uint64, r uint32, mix [][]uint32, datasetSize uint64, lookup LookupFunc, l1 []uint32) [][]uint32 {
	state := fill_mix_init(seed, uint32(RegisterCount))
	numItems := uint32(datasetSize / (2 * 128))
	itemIndex := mix[r%uint32(LaneCount)][0] % numItems

	item := lookup(itemIndex)

	numWordsPerLane := len(item) / LaneCount
	maxOperations := max(RoundMathOperations, RoundMathOperations)
	for i := 0; i < maxOperations; i++ {
		if i < RoundCacheAccesses {
			dst := state.nextDst()
			src := state.nextSrc()
			sel := state.nextRng()

			for l := 0; l < LaneCount; l++ {
				offset := mix[l][src] % (uint32(CacheBytes) / 4)
				mix[l][dst] = merge(mix[l][src], l1[offset], sel)
			}
		}

		if i < RoundMathOperations {
			srcRand := state.nextRng() % (uint32(RegisterCount) * uint32(RegisterCount-1))
			src1 := srcRand % uint32(RegisterCount)
			src2 := srcRand / uint32(RegisterCount)

			if src1 >= src2 {
				src2 += 1
			}

			dst := state.nextDst()
			sel2 := state.nextSrc()
			sel1 := state.nextRng()

			for l := 0; l < LaneCount; l++ {
				data := math(mix[l][src1], mix[l][src2], sel1)
				mix[l][dst] = merge(mix[l][dst], data, sel2)
			}
		}
	}

	dsts := make([]uint32, numWordsPerLane)
	sels := make([]uint32, numWordsPerLane)
	for i := 0; i < numWordsPerLane; i++ {
		if i == 0 {
			dsts[i] = 0
		} else {
			dsts[i] = state.nextDst()
		}

		sels[i] = state.nextRng()
	}

	for l := 0; l < LaneCount; l++ {
		offset := ((uint32(l) ^ r) % uint32(LaneCount)) * uint32(numWordsPerLane)
		for i := 0; i < numWordsPerLane; i++ {
			word := item[offset+uint32(i)]
			mix[l][dsts[i]] = merge(mix[l][dsts[i]], word, sels[i])
		}
	}
	return mix
}

func Hash(cfg, height, seed, datasetSize uint64, lookup LookupFunc, l1 []uint32) []byte {
	mix := initMix(seed, LaneCount, RegisterCount)

	number := height / uint64(PeriodLength)
	for i := 0; i < RoundCount; i++ {
		mix = round(number, uint32(i), mix, datasetSize, lookup, l1)
	}

	laneHash := make([]uint32, LaneCount)
	for l := range laneHash {
		laneHash[l] = fnvOffsetBasis

		for i := 0; i < RegisterCount; i++ {
			laneHash[l] = Fnv1a(laneHash[l], mix[l][i])
		}
	}

	numWords := 8
	mixHash := make([]uint32, numWords)
	for i := 0; i < numWords; i++ {
		mixHash[i] = fnvOffsetBasis
	}

	for l := 0; l < LaneCount; l++ {
		mixHash[l%numWords] = Fnv1a(mixHash[l%numWords], laneHash[l])
	}

	return utils.Uint32ArrayToBytesLE(mixHash)
}

func finalize(seed [25]uint32, mixHash []byte) []byte {
	var state [25]uint32
	for i := 0; i < 8; i++ {
		state[i] = seed[i]
		state[i+8] = binary.LittleEndian.Uint32(mixHash[i*4 : i*4+4])
	}

	keccak.KeccakF800(&state)

	digest := utils.Uint32ArrayToBytesLE(state[:8])

	return digest
}
