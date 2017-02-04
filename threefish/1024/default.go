package threefish

func encrypt1024(state *[16]uint64, key *[17]uint64, tweak *[3]uint64) {
	encrypt1024Simple(state, key, tweak)
}
func decrypt1024(state *[16]uint64, key *[17]uint64, tweak *[3]uint64) {
	decrypt1024Simple(state, key, tweak)
}
