package bds

// Adapter satisfies asterix.BDSDecoder via structural typing — no import needed.
type Adapter struct{}

func (a *Adapter) DecodeBDS(bdsRegisterAddress byte, data []byte) (any, error) {
	bds1 := (bdsRegisterAddress >> 4) & 0x0F
	bds2 := bdsRegisterAddress & 0x0F
	return Decode(bds1, bds2, data)
}
