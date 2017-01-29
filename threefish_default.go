// +build !arm

package skein

func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	encrypt256Simple(state, key, tweak)
}
func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	decrypt256Simple(state, key, tweak)
}
func encrypt512(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	encrypt512Simple(state, key, tweak)
}
func decrypt512(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	decrypt512Simple(state, key, tweak)
}
func encrypt1024(state *[16]uint64, key *[17]uint64, tweak *[3]uint64) {
	encrypt1024Simple(state, key, tweak)
}
func decrypt1024(state *[16]uint64, key *[17]uint64, tweak *[3]uint64) {
	decrypt1024Simple(state, key, tweak)
}
