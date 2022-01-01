package progpow

type kiss99 struct {
	z     uint32
	w     uint32
	jsr   uint32
	jcong uint32
}

func New(z, w, jsr, jcong uint32) *kiss99 {
	return &kiss99{
		z,
		w,
		jsr,
		jcong,
	}
}

//Pseudorandom number generator
func (k *kiss99) Next() uint32 {
	k.z = 36969*(k.z&65535) + (k.z >> 16)
	k.w = 18000*(k.w&65535) + (k.w >> 16)

	MWC := (k.z << 16) + k.w

	k.jsr = k.jsr ^ (k.jsr << 17)
	k.jsr = k.jsr ^ (k.jsr >> 13)
	k.jsr = k.jsr ^ (k.jsr << 5)

	k.jcong = 69069*k.jcong + 1234567

	return (MWC ^ k.jcong) + k.jsr
}
