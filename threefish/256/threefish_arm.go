package threefish

// Hand-rolled asm for these
func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
