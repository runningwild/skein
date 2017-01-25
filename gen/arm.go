package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {
	t := template.New("sched")
	if _, err := t.Parse(keySchedEncrypt); err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, &skeinParams{
		KeyLenPlusOne:  5,
		KeySchedOffset: 32, // 32 bytes of state
	})
	fmt.Printf("%s\n", buf.Bytes())
}

type skeinParams struct {
	KeyLenPlusOne  int
	KeySchedOffset int
}

// TODO: skeinParams should only require KeyLen and generate the rest on its own
// TODO: add tweak values and the round number.  round can probably be added last,
//		 but either way it would be nice to do it without storing the state twice.

const keySchedEncrypt = `
TEXT	·encrypt256(SB), $-4-24

keySched:
	MOVW	stateIndex+4(FP),	R7	// State index
	MOVW 	round+0(FP),	R8	// Round number
	MOVW	key+16(FP),	R4	// Pointer to key data
	MOVW	state+12(FP),	R5	// Pointer to state

	MOVW	(R7),	R7	// Load the state index
	MOVW	$1,	R6	// Loop counter.  Starts at 1 because we end at KeyLenPlusOne
	MOVW	R7<<3,	R7	// Get offset in uint64s
	ADD	R4,	R7	// R7 is now the addr of the first key data to use

	// Loop goes one uint64 at a time.
	keySchedLoop:
		// Load state
		MOVW	0(R5),	R0	// Load first word of the state
		MOVW	0(R7), 	R2	// Load first word of the key
		MOVW	4(R5),	R1	// Load second word of the state
		MOVW	4(R7),	R3	// Load second word of the key
		ADD	$8,	R7	// Increment key pointer
		ADD.S	R2,	R0	// Add first word of key to first word of state
		MOVW	R0,	0(R5)	// Store first word of state
		ADC	R3,	R1	// Add second word of key to second word of state
		MOVW	R1,	4(R5)	// Store second word of state
		ADD	$1,	R6	// Increment loop counter
		TEQ	${{.KeyLenPlusOne}},	R6	// Check if we're done with the loop
		B.EQ	keySchedDone
		ADD	$8,	R5	// Increment state pointer
		B	keySchedLoop

	keySchedDone:

	// Add in the round number while we still have the last value of the state in registers
	ADD.S	R8,	R0	// Add the round number to the last value in the state
	ADC	$0,	R1	// Carry
	MOVW	R0,	0(R5)	// Store first word of state
	MOVW	R1,	4(R5)	// Store second word of state

	// Store updated state pointer
	MOVW	stateIndex+4(FP),	R2
	MOVW	(R2),	R3
	ADD	$1,	R3
	MOVW	${{.KeyLenPlusOne}},	R0
	CMP	R3,	R0	// Compare against the natural key length
	SUB.LE	${{.KeyLenPlusOne}},	R3	// If we're over, subtract out the key length so we stay in bounds
	MOVW	R3, (R2)

	// Now we have to add in the tweak values to the previous third- and second-to-last values in the state.
	MOVW	tweakIndex+8(FP),	R8	// Tweak index
	MOVW	tweak+20(FP),	R6	// Pointer to tweak
	MOVW	$2,	R0
	MOVW	(R8),	R7
	CMP	R7,	R0	// Compare against the natural tweak length
	SUB.LE	$3,	R7	// If we're almost over, subtract out the tweak length so we stay in bounds
	ADD	$1,	R7, R4	// This is the udpated tweak index
	MOVW	R4,	(R8)	// Store the updated tweak index (i.e. the index we'll use on the next call to this function)
	MOVW	R7<<3,	R7	// Get offset in uint64s
	ADD	R6,	R7	// R7 is now the addr of the first tweak data to use

	// Add the first tweak value to the third-to-last element of state
	MOVW	-16(R5),	R0	// Load first word of third-to-last element of the state
	MOVW	-12(R5),	R1	// Load second word of third-to-last element of the state
	MOVW	0(R7),	R2	// Load first word of the tweak
	MOVW	4(R7),	R3	// Load second word of the tweak
	ADD.S	R2,	R0	// Add first word of tweak to first word of state
	ADC	R3,	R1	// Add second word of tweak to second word of state
	MOVW	R0,	-16(R5)	// Store first word of state
	MOVW	R1,	-12(R5)	// Store second word of state

	// Add the second tweak value to the second-to-last element of state
	MOVW	-8(R5),	R0	// Load first word of second-to-last element of the state
	MOVW	-4(R5),	R1	// Load second word of second-to-last element of the state
	MOVW	8(R7),	R2	// Load first word of the tweak
	MOVW	12(R7),	R3	// Load second word of the tweak
	ADD.S	R2,	R0	// Add first word of tweak to first word of state
	ADC	R3,	R1	// Add second word of tweak to second word of state
	MOVW	R0,	-8(R5)	// Store first word of state
	MOVW	R1,	-4(R5)	// Store second word of state

	RET
`

var rotations256 = [8][]uint{
	{14, 16},
	{52, 57},
	{23, 40},
	{5, 37},
	{25, 33},
	{46, 12},
	{58, 22},
	{32, 32},
}

var rotations512 = [8][]uint{
	{46, 36, 19, 37},
	{33, 27, 14, 42},
	{17, 49, 36, 39},
	{44, 9, 54, 56},
	{39, 30, 34, 24},
	{13, 50, 10, 17},
	{25, 29, 39, 43},
	{8, 35, 56, 22},
}

var rotations1024 = [8][]uint{
	{24, 13, 8, 47, 8, 17, 22, 37},
	{38, 19, 10, 55, 49, 18, 23, 52},
	{33, 4, 51, 13, 34, 41, 59, 17},
	{5, 20, 48, 41, 47, 28, 16, 25},
	{41, 9, 37, 31, 12, 47, 44, 30},
	{16, 34, 56, 51, 4, 53, 42, 41},
	{31, 44, 47, 46, 19, 42, 44, 25},
	{9, 48, 35, 52, 23, 31, 37, 20},
}
