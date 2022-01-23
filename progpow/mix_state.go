package progpow

type mixState struct {
	DstCounter uint32
	SrcCounter uint32
	Rng        *kiss99
	SrcSeq     [progpowRegs]uint32
	DstSeq     [progpowRegs]uint32
}

func (m mixState) nextDst() uint32 {
	val := m.DstSeq[m.DstCounter%progpowRegs]
	m.DstCounter += 1
	return val
}

func (m mixState) nextSrc() uint32 {
	val := m.SrcSeq[m.SrcCounter%progpowRegs]
	m.SrcCounter += 1
	return val
}

func (s mixState) rng() uint32 {
	return s.Rng.Next()
}

func fill_mix(seed uint64, size uint32) *mixState {
	var z, w, jsr, jcong uint32

	z = Fnv1a(fnvoffSetBasis, uint32(seed))
	w = Fnv1a(z, uint32(seed>>32))
	jsr = Fnv1a(w, uint32(seed))
	jcong = Fnv1a(jsr, uint32(seed>>32))

	rng := New(z, w, jsr, jcong)

	var srcSeq [progpowRegs]uint32
	var dstSeq [progpowRegs]uint32

	for i := uint32(0); i < progpowRegs; i++ {
		dstSeq[i] = i
		srcSeq[i] = i
	}

	//Using Fisher-Yates Shuffle
	for i := uint32(progpowRegs); i > 1; i-- {
		dstInd := rng.Next() % i
		dstSeq[i-1], dstSeq[dstInd] = dstSeq[dstInd], dstSeq[i-1]

		srcInd := rng.Next() % i
		srcSeq[i-1], srcSeq[srcInd] = srcSeq[srcInd], srcSeq[i-1]
	}

	return &mixState{0, 0, rng, srcSeq, dstSeq}
}
