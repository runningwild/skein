// func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
TEXT	·decrypt256(SB), $-4-12
RET

// func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
TEXT	·encrypt256(SB), $-4-12
	// Extend the key
	MOVW	key+4(FP), R12
	MOVW	$2851871266,R0
	MOVW	(R12),		R1
	EOR		R1,			R0
	MOVW	8(R12),		R1
	EOR		R1,			R0
	MOVW	16(R12),	R1
	EOR		R1,			R0
	MOVW	24(R12),	R1
	EOR		R1,			R0
	MOVW	R0,			32(R12)

	MOVW	$466688986,	R0
	MOVW	4(R12),		R1
	EOR		R1,			R0
	MOVW	12(R12),	R1
	EOR		R1,			R0
	MOVW	20(R12),	R1
	EOR		R1,			R0
	MOVW	28(R12),	R1
	EOR		R1,			R0
	MOVW	R0,			36(R12)

	// Extend the tweak
	MOVW	tweak+8(FP),R12
	MOVW	(R12),		R0
	MOVW	8(R12),		R1
	EOR		R0,			R1
	MOVW	R1,			16(R12)
	MOVW	4(R12),		R0
	MOVW	12(R12),	R1
	EOR		R0,			R1
	MOVW	R1,			20(R12)

	// Load the full state
	MOVW	state(FP),	R12
	MOVW	(R12),		R0
	MOVW	4(R12),		R1
	MOVW	8(R12),		R2
	MOVW	12(R12),	R3
	MOVW	16(R12),	R4
	MOVW	20(R12),	R5
	MOVW	24(R12),	R6
	MOVW	28(R12),	R7


	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 0 (round number)
	ADD.S	$0,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 1 (round number)
	ADD.S	$1,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 2 (round number)
	ADD.S	$2,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 3 (round number)
	ADD.S	$3,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 4 (round number)
	ADD.S	$4,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 5 (round number)
	ADD.S	$5,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 6 (round number)
	ADD.S	$6,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 7 (round number)
	ADD.S	$7,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 8 (round number)
	ADD.S	$8,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 9 (round number)
	ADD.S	$9,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 10 (round number)
	ADD.S	$10,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 11 (round number)
	ADD.S	$11,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 12 (round number)
	ADD.S	$12,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 13 (round number)
	ADD.S	$13,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 14 (round number)
	ADD.S	$14,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 15 (round number)
	ADD.S	$15,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11

	ADC		R11,		R5	// state[3] += 16 (round number)
	ADD.S	$16,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 14
	MOVW	R2<<14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>18,	R12
	ORR		R12,	R11
	MOVW	R3<<14,	R3
	ORR		R2>>18,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 16
	MOVW	R6<<16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>16,	R12
	ORR		R12,	R11
	MOVW	R7<<16,	R7
	ORR	R6>>16,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 20
	MOVW	R6>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<20,	R12
	ORR		R12,	R11
	MOVW	R7>>12,	R7
	ORR		R6<<20,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 25
	MOVW	R2>>7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<25,	R12
	ORR		R12,	R11
	MOVW	R3>>7,	R3
	ORR		R2<<25,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 23
	MOVW	R2<<23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>9,	R12
	ORR		R12,	R11
	MOVW	R3<<23,	R3
	ORR		R2>>9,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 8
	MOVW	R6>>24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<8,	R12
	ORR		R12,	R11
	MOVW	R7>>24,	R7
	ORR		R6<<8,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1

	// rot << 5
	MOVW	R6<<5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>27,	R12
	ORR		R12,	R11
	MOVW	R7<<5,	R7
	ORR		R6>>27,	R7
	MOVW	R11,	R6



	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5



	// rot >> 5
	MOVW	R2>>27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<5,	R12
	ORR		R12,	R11
	MOVW	R3>>27,	R3
	ORR		R2<<5,	R3
	MOVW	R11,	R2

	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11

	ADC		R11,		R5	// state[3] += 17 (round number)
	ADD.S	$17,			R6
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// rot << 25
	MOVW	R2<<25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>7,	R12
	ORR		R12,	R11
	MOVW	R3<<25,	R3
	ORR		R2>>7,	R3
	MOVW	R11,	R2



	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5



	// rot >> 1
	MOVW	R6>>31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<1,	R12
	ORR		R12,	R11
	MOVW	R7>>31,	R7
	ORR		R6<<1,	R7
	MOVW	R11,	R6

	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1



	// rot >> 14
	MOVW	R6>>18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<14,	R12
	ORR		R12,	R11
	MOVW	R7>>18,	R7
	ORR		R6<<14,	R7
	MOVW	R11,	R6

	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5

	// rot << 12
	MOVW	R2<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>20,	R12
	ORR		R12,	R11
	MOVW	R3<<12,	R3
	ORR	R2>>20,	R3
	MOVW	R11,	R2



	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1



	// rot >> 26
	MOVW	R2>>6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<26,	R12
	ORR		R12,	R11
	MOVW	R3>>6,	R3
	ORR		R2<<26,	R3
	MOVW	R11,	R2

	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R4			// y0 = x0 + x1
	ADC		R7,		R5

	// rot << 22
	MOVW	R6<<22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>10,	R12
	ORR		R12,	R11
	MOVW	R7<<22,	R7
	ORR	R6>>10,	R7
	MOVW	R11,	R6



	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R6,		R0			// y0 = x0 + x1
	ADC		R7,		R1


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Mix state[2] and state[3]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R4			// y0 = x0 + x1
	ADC		R3,		R5


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[0]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[1]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[2]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[3]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[1]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11

	ADC		R11,		R5	// state[3] += 18 (round number)
	ADD.S	$18,			R6
	ADC		$0,			R7




	// Store the full state
	MOVW	state(FP),	R12
	MOVW	R0,			(R12)
	MOVW	R1,			4(R12)
	MOVW	R2,			8(R12)
	MOVW	R3,			12(R12)
	MOVW	R4,			16(R12)
	MOVW	R5,			20(R12)
	MOVW	R6,			24(R12)
	MOVW	R7,			28(R12)

	RET
