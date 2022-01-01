package progpow

const (
	fnvPrime uint32 = 0x01000193
)

//Has flaw that can be used by ASICs and FPGAs to decrease the level of computation used but still used in DAG  generation
func Fnv1(u, v uint32) uint32 {
	return (u * fnvPrime) ^ v
}

// Has better distribution properties then Fnv1 and in use the main progpow loop
func Fnv1a(u, v uint32) uint32 {
	return (u ^ v) * fnvPrime
}

func FnvHash(mix []uint32, data []uint32) {
	for i := 0; i < len(mix); i++ {
		mix[i] = Fnv1(mix[i], data[i])
	}
}
