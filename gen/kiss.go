package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"text/template"
)

var (
	size = flag.Int("size", 256, "what size of threefish to write (256, 512, 1024)")
)

func main() {
	flag.Parse()
	var params SkeinParams
	switch *size {
	case 256:
		params = skein256Params
	case 512:
		params = skein512Params
	case 1024:
		params = skein1024Params
	default:
		fmt.Printf("--size must be one of 256, 512, or 1024.\n")
		os.Exit(1)
	}
	enc := template.New("encrypt")
	if _, err := enc.Parse(encryptTemplate); err != nil {
		panic(err)
	}
	dec := template.New("decrypt")
	if _, err := dec.Parse(decryptTemplate); err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	buf.WriteString("package skein\n")
	enc.Execute(buf, &params)
	dec.Execute(buf, &params)

	fmt.Printf("%s\n", buf.Bytes())
}

type SkeinParams struct {
	KeyBits     int
	KeyBytes    int
	KeyInt64s   int
	StateInt64s int
	Rounds      int
	Permutation []int
	Rotations   [8][]int
}

func (p *SkeinParams) EncryptParams() []roundParams {
	var rp []roundParams
	order := make([]int, len(p.Permutation))
	for i := range order {
		order[i] = i
	}
	for i := 0; i < p.Rounds; i++ {
		var mixes [][3]int
		for j := 0; j < len(order); j += 2 {
			mixes = append(mixes, [3]int{order[j], order[j+1], p.Rotations[i%8][j/2]})
		}
		rp = append(rp, roundParams{
			Round:       i,
			Mixes:       mixes,
			KeySchedule: i%4 == 0,
			Schedule:    keySchedParams(p.KeyBytes, i),
		})
		for i := range order {
			order[i] = p.Permutation[order[i]]
		}
	}
	rp = append(rp, roundParams{Round: p.Rounds, KeySchedule: true, Schedule: keySchedParams(p.KeyBytes, p.Rounds)})
	return rp
}

func (p *SkeinParams) DecryptParams() []roundParams {
	var rp []roundParams
	order := make([]int, len(p.Permutation))
	for i := range order {
		order[i] = i
	}
	for i := 0; i < p.Rounds; i++ {
		var mixes [][3]int
		for j := 0; j < len(order); j += 2 {
			mixes = append(mixes, [3]int{order[j], order[j+1], p.Rotations[i%8][j/2]})
		}
		rp = append(rp, roundParams{
			Round:       i,
			Mixes:       mixes,
			KeySchedule: i%4 == 0,
			Schedule:    keySchedParams(p.KeyBytes, i),
		})
		for i := range order {
			order[i] = p.Permutation[order[i]]
		}
	}
	rp = append(rp, roundParams{Round: p.Rounds, KeySchedule: true, Schedule: keySchedParams(p.KeyBytes, p.Rounds)})
	for i := 0; i < len(rp)/2; i++ {
		swap := len(rp) - i - 1
		rp[i], rp[swap] = rp[swap], rp[i]
	}
	for i := range rp {
		for j := 0; j < len(rp[i].Mixes)/2; j++ {
			swap := len(rp[i].Mixes) - j - 1
			rp[i].Mixes[j], rp[i].Mixes[swap] = rp[i].Mixes[swap], rp[i].Mixes[j]
		}
	}
	return rp
}

func keySchedParams(keyBytes, round int) [][3]int {
	var sched [][3]int
	for j := 0; j < keyBytes/8; j++ {
		var v = [3]int{(round/4 + j) % (keyBytes/8 + 1), -1, -1}
		sched = append(sched, v)
	}
	sched[len(sched)-3][1] = (round / 4) % 3
	sched[len(sched)-2][1] = (round/4 + 1) % 3
	sched[len(sched)-1][2] = round / 4
	return sched
}

type roundParams struct {
	Round int
	Mixes [][3]int // first two are the indexes, third is the rotation

	KeySchedule bool
	Schedule    [][3]int
}

const encryptTemplate = `
func encrypt{{.KeyBits}}Simple(state *[{{.StateInt64s}}]uint64, key *[{{.KeyInt64s}}]uint64, tweak *[3]uint64) {
	key[{{.StateInt64s}}] = c240{{range $index, $dummy := .Permutation}} ^ key[{{$index}}]{{end}}
	tweak[2] = tweak[0] ^ tweak[1]
	{{range $index, $dummy := .Permutation}}s{{$index}} := state[{{$index}}]
	{{end}}{{range $params := .EncryptParams}}
	// Round {{$params.Round}}
	{{if $params.KeySchedule}}
	// Key Schedule{{range $index, $inject := $params.Schedule}}
	s{{$index}} += key[{{index $inject 0}}] {{if ge (index $inject 1) 0}}+ tweak[{{index $inject 1}}]{{end}}{{if ge (index $inject 2) 0}}+ {{index $inject 2}}{{end}}{{end}}
	{{end}}{{range $index, $mix := $params.Mixes}}
	// Mix {{index $mix 0}} with {{index $mix 1}}
	s{{index $mix 0}} += s{{index $mix 1}}
	s{{index $mix 1}} = (s{{index $mix 1}} << {{index $mix 2}}) | (s{{index $mix 1}} >> (64 - {{index $mix 2}})) ^ s{{index $mix 0}}
	{{end}}{{end}}{{range $index, $dummy := .Permutation}}
	state[{{$index}}] = s{{$index}}{{end}}
}
`
const decryptTemplate = `
func decrypt{{.KeyBits}}Simple(state *[{{.StateInt64s}}]uint64, key *[{{.KeyInt64s}}]uint64, tweak *[3]uint64) {
	key[{{.StateInt64s}}] = c240{{range $index, $dummy := .Permutation}} ^ key[{{$index}}]{{end}}
	tweak[2] = tweak[0] ^ tweak[1]
	{{range $index, $dummy := .Permutation}}s{{$index}} := state[{{$index}}]
	{{end}}{{range $params := .DecryptParams}}
	// Round {{$params.Round}}
	{{range $index, $mix := $params.Mixes}}
	// Mix {{index $mix 0}} with {{index $mix 1}}
	s{{index $mix 1}} ^= s{{index $mix 0}}
	s{{index $mix 1}} = (s{{index $mix 1}} >> {{index $mix 2}}) | (s{{index $mix 1}} << (64 - {{index $mix 2}}))
	s{{index $mix 0}} -= s{{index $mix 1}}
	{{end}}{{if $params.KeySchedule}}
	// Key Schedule{{range $index, $inject := $params.Schedule}}
	s{{$index}} -= key[{{index $inject 0}}] {{if ge (index $inject 1) 0}}+ tweak[{{index $inject 1}}]{{end}}{{if ge (index $inject 2) 0}}+ {{index $inject 2}}{{end}}{{end}}
	{{end}}{{end}}{{range $index, $dummy := .Permutation}}
	state[{{$index}}] = s{{$index}}{{end}}
}
`

var skein256Params = SkeinParams{
	KeyBytes:    32,
	KeyBits:     256,
	StateInt64s: 4,
	KeyInt64s:   5,
	Rounds:      72,
	Permutation: []int{0, 3, 2, 1},
	Rotations: [8][]int{
		{14, 16},
		{52, 57},
		{23, 40},
		{5, 37},
		{25, 33},
		{46, 12},
		{58, 22},
		{32, 32},
	},
}
var skein512Params = SkeinParams{
	KeyBytes:    64,
	KeyBits:     512,
	StateInt64s: 8,
	KeyInt64s:   9,
	Rounds:      72,
	Permutation: []int{2, 1, 4, 7, 6, 5, 0, 3},
	Rotations: [8][]int{
		{46, 36, 19, 37},
		{33, 27, 14, 42},
		{17, 49, 36, 39},
		{44, 9, 54, 56},
		{39, 30, 34, 24},
		{13, 50, 10, 17},
		{25, 29, 39, 43},
		{8, 35, 56, 22},
	},
}
var skein1024Params = SkeinParams{
	KeyBytes:    128,
	KeyBits:     1024,
	StateInt64s: 16,
	KeyInt64s:   17,
	Rounds:      80,
	Permutation: []int{0, 9, 2, 13, 6, 11, 4, 15, 10, 7, 12, 3, 14, 5, 8, 1},
	Rotations: [8][]int{
		{24, 13, 8, 47, 8, 17, 22, 37},
		{38, 19, 10, 55, 49, 18, 23, 52},
		{33, 4, 51, 13, 34, 41, 59, 17},
		{5, 20, 48, 41, 47, 28, 16, 25},
		{41, 9, 37, 31, 12, 47, 44, 30},
		{16, 34, 56, 51, 4, 53, 42, 41},
		{31, 44, 47, 46, 19, 42, 44, 25},
		{9, 48, 35, 52, 23, 31, 37, 20},
	},
}
