package main

import (
	"bytes"
	"flag"
	"fmt"
	"text/template"
)

var (
	rounds = flag.Int("rounds", 72, "num rounds")
)

func main() {
	flag.Parse()
	var ht, kt, mt *template.Template
	{
		ht = template.New("encryptHeader")
		_, err := ht.Parse(header)
		if err != nil {
			panic(err)
		}
		kt = template.New("encryptKeySched")
		_, err = kt.Parse(keySchedEncrypt)
		if err != nil {
			panic(err)
		}
		mt = template.New("encryptMix")
		_, err = mt.Parse(mixEncrypt)
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(nil)
		ht.Execute(buf, doParams(0, true))
		for r := 0; r < *rounds; r++ {
			p := doParams(r, true)
			if r%4 == 0 {
				kt.Execute(buf, p)
			}
			mt.Execute(buf, p)
		}
		if *rounds%4 == 0 {
			kt.Execute(buf, doParams(*rounds, true))
		}
		buf.WriteString(footer)

		fmt.Printf("%s\n", buf.Bytes())
	}
	{
		ht = template.New("decryptHeader")
		_, err := ht.Parse(header)
		if err != nil {
			panic(err)
		}
		kt = template.New("decryptKeySched")
		_, err = kt.Parse(keySchedDecrypt)
		if err != nil {
			panic(err)
		}
		mt = template.New("decryptMix")
		_, err = mt.Parse(mixDecrypt)
		if err != nil {
			panic(err)
		}

		buf := bytes.NewBuffer(nil)
		ht.Execute(buf, doParams(*rounds, false))
		if *rounds%4 == 0 {
			kt.Execute(buf, doParams(*rounds, false))
		}
		for r := *rounds - 1; r >= 0; r-- {
			p := doParams(r, false)
			mt.Execute(buf, p)
			if r%4 == 0 {
				kt.Execute(buf, p)
			}
		}
		buf.WriteString(footer)

		fmt.Printf("%s\n", buf.Bytes())
	}
}

type roundParams struct {
	Pre                                    string
	K0a, K0b, K1a, K1b, K2a, K2b, K3a, K3b int
	T0a, T0b, T1a, T1b                     int
	R, D                                   int
	M0, M1, M2, M3, M4, M5, M6, M7         int
	Rot0, Rot0X, Rot1, Rot1X               int
	Low0, Low1, Mid0, Mid1, High0, High1   bool
}

func doParams(r int, encrypt bool) *roundParams {
	var rp roundParams
	if encrypt {
		rp.Pre = "en"
	} else {
		rp.Pre = "de"
	}
	kr := r / 4
	rp.K0a = (kr % 5) * 8
	rp.K0b = (kr%5)*8 + 4
	rp.K1a = ((kr + 1) % 5) * 8
	rp.K1b = ((kr+1)%5)*8 + 4
	rp.K2a = ((kr + 2) % 5) * 8
	rp.K2b = ((kr+2)%5)*8 + 4
	rp.K3a = ((kr + 3) % 5) * 8
	rp.K3b = ((kr+3)%5)*8 + 4

	rp.T0a = (kr % 3) * 8
	rp.T0b = (kr%3)*8 + 4
	rp.T1a = ((kr + 1) % 3) * 8
	rp.T1b = ((kr+1)%3)*8 + 4

	rp.M0 = 0
	rp.M1 = 1
	rp.M2 = 2
	rp.M3 = 3
	rp.M4 = 4
	rp.M5 = 5
	rp.M6 = 6
	rp.M7 = 7
	if r%2 == 1 {
		rp.M2 = 6
		rp.M3 = 7
		rp.M6 = 2
		rp.M7 = 3
	}

	switch r % 8 {
	case 0:
		rp.Rot0 = 14
		rp.Rot1 = 16
	case 1:
		rp.Rot0 = 52
		rp.Rot1 = 57
	case 2:
		rp.Rot0 = 23
		rp.Rot1 = 40
	case 3:
		rp.Rot0 = 5
		rp.Rot1 = 37
	case 4:
		rp.Rot0 = 25
		rp.Rot1 = 33
	case 5:
		rp.Rot0 = 46
		rp.Rot1 = 12
	case 6:
		rp.Rot0 = 58
		rp.Rot1 = 22
	case 7:
		rp.Rot0 = 32
		rp.Rot1 = 32
	}
	if !encrypt {
		rp.Rot0 = 64 - rp.Rot0
		rp.Rot1 = 64 - rp.Rot1
	}

	if rp.Rot0 == 32 {
		rp.Mid0 = true
	} else if rp.Rot0 < 32 {
		rp.Low0 = true
		rp.Rot0X = 32 - rp.Rot0
	} else {
		rp.High0 = true
		rp.Rot0 -= 32
		rp.Rot0X = 32 - rp.Rot0
	}
	if rp.Rot1 == 32 {
		rp.Mid1 = true
	} else if rp.Rot1 < 32 {
		rp.Low1 = true
		rp.Rot1X = 32 - rp.Rot1
	} else {
		rp.High1 = true
		rp.Rot1 -= 32
		rp.Rot1X = 32 - rp.Rot1
	}

	rp.R = r
	rp.D = r / 4
	return &rp
}

const keySchedEncrypt = `
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	{{.K0a}}(R12),		R8		// state[0] += key[{{.D}}%5]
	MOVW	{{.K0b}}(R12),		R11
	ADD.S	R8,		R0
	MOVW	{{.K1a}}(R12),		R8		// state[1] += key[({{.D}}+1)%5]
	ADC		R11,		R1
	MOVW	{{.K1b}}(R12),	R11
	ADD.S	R8,		R2
	MOVW	{{.K2a}}(R12),	R8		// state[2] += key[({{.D}}+2)%5]
	ADC		R11,		R3
	MOVW	{{.K2b}}(R12),	R11
	ADD.S	R8,		R4
	MOVW	{{.K3a}}(R12),	R8		// state[3] += key[({{.D}}+3)%5]
	ADC		R11,		R5
	MOVW	{{.K3b}}(R12),	R11
	MOVW	tweak+8(FP),R12
	ADD.S	R8,		R6
	ADC		R11,		R7

	MOVW	{{.T0a}}(R12),		R8		// state[1] += tweak[{{.D}}%3]
	MOVW	{{.T0b}}(R12),		R11
	ADD.S	R8,		R2
	MOVW	{{.T1a}}(R12),		R8		// state[2] += tweak[({{.D}}+1)%3]
	ADC		R11,		R3
	MOVW	{{.T1b}}(R12),	R11
	ADD.S	R8,		R4
	ADC		R11,		R5

	ADD.S	${{.D}},			R6		// state[3] += {{.D}} (round number)
	ADC		$0,			R7
`
const keySchedDecrypt = `
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	{{.K0a}}(R12),		R8		// state[0] += key[{{.D}}%5]
	MOVW	{{.K0b}}(R12),		R11
	SUB.S	R8,		R0
	MOVW	{{.K1a}}(R12),		R8		// state[1] += key[({{.D}}+1)%5]
	SBC		R11,		R1
	MOVW	{{.K1b}}(R12),	R11
	SUB.S	R8,		R2
	MOVW	{{.K2a}}(R12),	R8		// state[2] += key[({{.D}}+2)%5]
	SBC		R11,		R3
	MOVW	{{.K2b}}(R12),	R11
	SUB.S	R8,		R4
	MOVW	{{.K3a}}(R12),	R8		// state[3] += key[({{.D}}+3)%5]
	SBC		R11,		R5
	MOVW	{{.K3b}}(R12),	R11
	MOVW	tweak+8(FP),R12
	SUB.S	R8,		R6
	MOVW	{{.T0a}}(R12),		R8		// state[1] += tweak[{{.D}}%3]
	SBC		R11,		R7

	MOVW	{{.T0b}}(R12),		R11
	SUB.S	R8,		R2
	MOVW	{{.T1a}}(R12),		R8		// state[2] += tweak[({{.D}}+1)%3]
	SBC		R11,		R3
	MOVW	{{.T1b}}(R12),	R11
	SUB.S	R8,		R4
	SBC		R11,		R5

	SUB.S	${{.D}},			R6		// state[3] += {{.D}} (round number)
	SBC		$0,			R7
`

const mixEncrypt = `
	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R{{.M2}},		R{{.M0}}			// y0 = x0 + x1
	ADC		R{{.M3}},		R{{.M1}}
{{if .Low0}}
	// Low0: rot << {{.Rot0}}
	MOVW	R{{.M2}}<<{{.Rot0}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M3}}>>{{.Rot0X}},	R11
	MOVW	R{{.M3}}<<{{.Rot0}},	R{{.M3}}
	ORR		R{{.M2}}>>{{.Rot0X}},	R{{.M3}}
	MOVW	R11,	R{{.M2}}
{{end}}
{{if .Mid0}}
	MOVW	R{{.M2}},	R11						// for rot==32 we can just swap
	MOVW	R{{.M3}},	R{{.M2}}
	MOVW	R11,	R{{.M3}}
{{end}}
{{if .High0}}
	// High0: rot >> {{.Rot0}}
	MOVW	R{{.M2}}>>{{.Rot0X}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M3}}<<{{.Rot0}},	R11
	MOVW	R{{.M3}}>>{{.Rot0X}},	R{{.M3}}
	ORR		R{{.M2}}<<{{.Rot0}},	R{{.M3}}
	MOVW	R11,	R{{.M2}}
{{end}}
	EOR		R{{.M0}},		R{{.M2}}			// y1 = y1 ^ y0
	EOR		R{{.M1}},		R{{.M3}}

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R{{.M6}},		R{{.M4}}			// y0 = x0 + x1
	ADC		R{{.M7}},		R{{.M5}}
{{if .Low1}}
	// Low1: rot << {{.Rot1}}
	MOVW	R{{.M6}}<<{{.Rot1}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M7}}>>{{.Rot1X}},	R11
	MOVW	R{{.M7}}<<{{.Rot1}},	R{{.M7}}
	ORR	R{{.M6}}>>{{.Rot1X}},	R{{.M7}}
	MOVW	R11,	R{{.M6}}
{{end}}
{{if .Mid1}}
	MOVW	R{{.M6}},	R11						// for rot==32 we can just swap
	MOVW	R{{.M7}},	R{{.M6}}
	MOVW	R11,	R{{.M7}}
{{end}}
{{if .High1}}
	// High1: rot >> {{.Rot1}}
	MOVW	R{{.M6}}>>{{.Rot1X}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M7}}<<{{.Rot1}},	R11
	MOVW	R{{.M7}}>>{{.Rot1X}},	R{{.M7}}
	ORR		R{{.M6}}<<{{.Rot1}},	R{{.M7}}
	MOVW	R11,	R{{.M6}}
{{end}}
	EOR		R{{.M4}},		R{{.M6}}			// y1 = y1 ^ y0
	EOR		R{{.M5}},		R{{.M7}}
		`

const mixDecrypt = `
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R{{.M4}},		R{{.M6}}			// y1 = y1 ^ y0
	EOR		R{{.M5}},		R{{.M7}}
{{if .Low1}}
	// Low1: rot << {{.Rot1}}
	MOVW	R{{.M6}}<<{{.Rot1}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M7}}>>{{.Rot1X}},	R11
	MOVW	R{{.M7}}<<{{.Rot1}},	R{{.M7}}
	ORR	R{{.M6}}>>{{.Rot1X}},	R{{.M7}}
	MOVW	R11,	R{{.M6}}
{{end}}
{{if .Mid1}}
	MOVW	R{{.M6}},	R11						// for rot==32 we can just swap
	MOVW	R{{.M7}},	R{{.M6}}
	MOVW	R11,	R{{.M7}}
{{end}}
{{if .High1}}
	// High1: rot >> {{.Rot1}}
	MOVW	R{{.M6}}>>{{.Rot1X}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M7}}<<{{.Rot1}},	R11
	MOVW	R{{.M7}}>>{{.Rot1X}},	R{{.M7}}
	ORR		R{{.M6}}<<{{.Rot1}},	R{{.M7}}
	MOVW	R11,	R{{.M6}}
{{end}}
	SUB.S	R{{.M6}},		R{{.M4}}			// y0 = x0 - x1
	SBC		R{{.M7}},		R{{.M5}}

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R{{.M0}},		R{{.M2}}			// y1 = y1 ^ y0
	EOR		R{{.M1}},		R{{.M3}}
{{if .Low0}}
	// Low0: rot << {{.Rot0}}
	MOVW	R{{.M2}}<<{{.Rot0}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M3}}>>{{.Rot0X}},	R11
	MOVW	R{{.M3}}<<{{.Rot0}},	R{{.M3}}
	ORR		R{{.M2}}>>{{.Rot0X}},	R{{.M3}}
	MOVW	R11,	R{{.M2}}
{{end}}
{{if .Mid0}}
	MOVW	R{{.M2}},	R11						// for rot==32 we can just swap
	MOVW	R{{.M3}},	R{{.M2}}
	MOVW	R11,	R{{.M3}}
{{end}}
{{if .High0}}
	// High0: rot >> {{.Rot0}}
	MOVW	R{{.M2}}>>{{.Rot0X}},	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	ORR	R{{.M3}}<<{{.Rot0}},	R11
	MOVW	R{{.M3}}>>{{.Rot0X}},	R{{.M3}}
	ORR		R{{.M2}}<<{{.Rot0}},	R{{.M3}}
	MOVW	R11,	R{{.M2}}
{{end}}
	SUB.S	R{{.M2}},		R{{.M0}}			// y0 = x0 - x1
	SBC		R{{.M3}},		R{{.M1}}
		`
const header = `
// func {{.Pre}}crypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
TEXT    ·{{.Pre}}crypt256(SB), $-4-12
    // Extend the key
    MOVW    key+4(FP), R12
    MOVW    $2851871266,R0
    MOVW    (R12),          R1
    MOVW    8(R12),         R2
    MOVW    16(R12),        R3
    MOVW    24(R12),        R4
    EOR             R1,                     R0
    EOR             R2,                     R0
    EOR             R3,                     R0
    EOR             R4,                     R0
    MOVW    R0,                     32(R12)

    MOVW    $466688986,     R0
    MOVW    4(R12),         R1
    MOVW    12(R12),        R2
    MOVW    20(R12),        R3
    MOVW    28(R12),        R4
    EOR             R1,                     R0
    EOR             R2,                     R0
    EOR             R3,                     R0
    EOR             R4,                     R0
    MOVW    R0,                     36(R12)

    // Extend the tweak
    MOVW    tweak+8(FP),R12
    MOVW    (R12),          R0
    MOVW    8(R12),         R2
    MOVW    4(R12),         R1
    MOVW    12(R12),        R3
    EOR             R0,                     R2
    MOVW    R2,                     16(R12)
    EOR             R1,                     R3
    MOVW    R3,                     20(R12)

    // Load the full state
    MOVW    state(FP),      R12
    MOVW    (R12),          R0
    MOVW    4(R12),         R1
    MOVW    8(R12),         R2
    MOVW    12(R12),        R3
    MOVW    16(R12),        R4
    MOVW    20(R12),        R5
    MOVW    24(R12),        R6
    MOVW    28(R12),        R7
`

const footer = `

    // Store the full state
    MOVW    state(FP),      R12
    MOVW    R0,                     (R12)
    MOVW    R1,                     4(R12)
    MOVW    R2,                     8(R12)
    MOVW    R3,                     12(R12)
    MOVW    R4,                     16(R12)
    MOVW    R5,                     20(R12)
    MOVW    R6,                     24(R12)
    MOVW    R7,                     28(R12)

    RET
`
