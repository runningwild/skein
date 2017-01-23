
// func encrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
TEXT    ·encrypt256(SB), $-4-12
    // Extend the key
    MOVW    key+4(FP), R12
    MOVW    $2851871266,R0
    MOVW    (R12),          R1
    EOR             R1,                     R0
    MOVW    8(R12),         R1
    EOR             R1,                     R0
    MOVW    16(R12),        R1
    EOR             R1,                     R0
    MOVW    24(R12),        R1
    EOR             R1,                     R0
    MOVW    R0,                     32(R12)

    MOVW    $466688986,     R0
    MOVW    4(R12),         R1
    EOR             R1,                     R0
    MOVW    12(R12),        R1
    EOR             R1,                     R0
    MOVW    20(R12),        R1
    EOR             R1,                     R0
    MOVW    28(R12),        R1
    EOR             R1,                     R0
    MOVW    R0,                     36(R12)

    // Extend the tweak
    MOVW    tweak+8(FP),R12
    MOVW    (R12),          R0
    MOVW    8(R12),         R1
    EOR             R0,                     R1
    MOVW    R1,                     16(R12)
    MOVW    4(R12),         R0
    MOVW    12(R12),        R1
    EOR             R0,                     R1
    MOVW    R1,                     20(R12)

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

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0%5]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(0+1)%5]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(0+2)%5]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(0+3)%5]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(0+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$0,			R6		// state[3] += 0 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	8(R12),		R11		// state[0] += key[1%5]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(1+1)%5]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(1+2)%5]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(1+3)%5]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[1%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(1+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$1,			R6		// state[3] += 1 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	16(R12),		R11		// state[0] += key[2%5]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(2+1)%5]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(2+2)%5]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(2+3)%5]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[2%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(2+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$2,			R6		// state[3] += 2 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	24(R12),		R11		// state[0] += key[3%5]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(3+1)%5]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(3+2)%5]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(3+3)%5]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[3%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(3+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$3,			R6		// state[3] += 3 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	32(R12),		R11		// state[0] += key[4%5]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(4+1)%5]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(4+2)%5]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(4+3)%5]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[4%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(4+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$4,			R6		// state[3] += 4 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	0(R12),		R11		// state[0] += key[5%5]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(5+1)%5]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(5+2)%5]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(5+3)%5]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[5%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(5+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$5,			R6		// state[3] += 5 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	8(R12),		R11		// state[0] += key[6%5]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(6+1)%5]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(6+2)%5]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(6+3)%5]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[6%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(6+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$6,			R6		// state[3] += 6 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	16(R12),		R11		// state[0] += key[7%5]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(7+1)%5]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(7+2)%5]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(7+3)%5]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[7%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(7+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$7,			R6		// state[3] += 7 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	24(R12),		R11		// state[0] += key[8%5]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(8+1)%5]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(8+2)%5]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(8+3)%5]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[8%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(8+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$8,			R6		// state[3] += 8 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	32(R12),		R11		// state[0] += key[9%5]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(9+1)%5]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(9+2)%5]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(9+3)%5]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[9%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(9+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$9,			R6		// state[3] += 9 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	0(R12),		R11		// state[0] += key[10%5]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(10+1)%5]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(10+2)%5]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(10+3)%5]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[10%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(10+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$10,			R6		// state[3] += 10 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	8(R12),		R11		// state[0] += key[11%5]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(11+1)%5]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(11+2)%5]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(11+3)%5]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[11%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(11+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$11,			R6		// state[3] += 11 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	16(R12),		R11		// state[0] += key[12%5]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(12+1)%5]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(12+2)%5]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(12+3)%5]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[12%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(12+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$12,			R6		// state[3] += 12 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	24(R12),		R11		// state[0] += key[13%5]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(13+1)%5]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(13+2)%5]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(13+3)%5]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[13%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(13+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$13,			R6		// state[3] += 13 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	32(R12),		R11		// state[0] += key[14%5]
	ADD.S	R11,		R0
	MOVW	36(R12),		R11
	ADC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(14+1)%5]
	ADD.S	R11,		R2
	MOVW	4(R12),	R11
	ADC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(14+2)%5]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(14+3)%5]
	ADD.S	R11,		R6
	MOVW	20(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[14%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(14+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$14,			R6		// state[3] += 14 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	0(R12),		R11		// state[0] += key[15%5]
	ADD.S	R11,		R0
	MOVW	4(R12),		R11
	ADC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(15+1)%5]
	ADD.S	R11,		R2
	MOVW	12(R12),	R11
	ADC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(15+2)%5]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(15+3)%5]
	ADD.S	R11,		R6
	MOVW	28(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[15%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(15+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$15,			R6		// state[3] += 15 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	8(R12),		R11		// state[0] += key[16%5]
	ADD.S	R11,		R0
	MOVW	12(R12),		R11
	ADC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(16+1)%5]
	ADD.S	R11,		R2
	MOVW	20(R12),	R11
	ADC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(16+2)%5]
	ADD.S	R11,		R4
	MOVW	28(R12),	R11
	ADC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(16+3)%5]
	ADD.S	R11,		R6
	MOVW	36(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[16%3]
	ADD.S	R11,		R2
	MOVW	12(R12),		R11
	ADC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(16+1)%3]
	ADD.S	R11,		R4
	MOVW	20(R12),	R11
	ADC		R11,		R5

	ADD.S	$16,			R6		// state[3] += 16 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 14
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

	// Low1: rot << 16
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



	// High0: rot >> 20
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



	// High1: rot >> 25
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

	// Low0: rot << 23
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



	// High1: rot >> 8
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

	// Low0: rot << 5
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



	// High1: rot >> 5
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
	MOVW	16(R12),		R11		// state[0] += key[17%5]
	ADD.S	R11,		R0
	MOVW	20(R12),		R11
	ADC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(17+1)%5]
	ADD.S	R11,		R2
	MOVW	28(R12),	R11
	ADC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(17+2)%5]
	ADD.S	R11,		R4
	MOVW	36(R12),	R11
	ADC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(17+3)%5]
	ADD.S	R11,		R6
	MOVW	4(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[17%3]
	ADD.S	R11,		R2
	MOVW	20(R12),		R11
	ADC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(17+1)%3]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5

	ADD.S	$17,			R6		// state[3] += 17 (round number)
	ADC		$0,			R7

	// Mix state[0] and state[1]
	// y0 = x0 + x1
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y1 = y1 ^ y0
	ADD.S	R2,		R0			// y0 = x0 + x1
	ADC		R3,		R1

	// Low0: rot << 25
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



	// High1: rot >> 1
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



	// High0: rot >> 14
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

	// Low1: rot << 12
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



	// High0: rot >> 26
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

	// Low1: rot << 22
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
	MOVW	24(R12),		R11		// state[0] += key[18%5]
	ADD.S	R11,		R0
	MOVW	28(R12),		R11
	ADC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(18+1)%5]
	ADD.S	R11,		R2
	MOVW	36(R12),	R11
	ADC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(18+2)%5]
	ADD.S	R11,		R4
	MOVW	4(R12),	R11
	ADC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(18+3)%5]
	ADD.S	R11,		R6
	MOVW	12(R12),	R11
	ADC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[18%3]
	ADD.S	R11,		R2
	MOVW	4(R12),		R11
	ADC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(18+1)%3]
	ADD.S	R11,		R4
	MOVW	12(R12),	R11
	ADC		R11,		R5

	ADD.S	$18,			R6		// state[3] += 18 (round number)
	ADC		$0,			R7


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


// func decrypt256(state *[4]uint64, key *[5]uint64, tweak *[3]uint64)
TEXT    ·decrypt256(SB), $-4-12
    // Extend the key
    MOVW    key+4(FP), R12
    MOVW    $2851871266,R0
    MOVW    (R12),          R1
    EOR             R1,                     R0
    MOVW    8(R12),         R1
    EOR             R1,                     R0
    MOVW    16(R12),        R1
    EOR             R1,                     R0
    MOVW    24(R12),        R1
    EOR             R1,                     R0
    MOVW    R0,                     32(R12)

    MOVW    $466688986,     R0
    MOVW    4(R12),         R1
    EOR             R1,                     R0
    MOVW    12(R12),        R1
    EOR             R1,                     R0
    MOVW    20(R12),        R1
    EOR             R1,                     R0
    MOVW    28(R12),        R1
    EOR             R1,                     R0
    MOVW    R0,                     36(R12)

    // Extend the tweak
    MOVW    tweak+8(FP),R12
    MOVW    (R12),          R0
    MOVW    8(R12),         R1
    EOR             R0,                     R1
    MOVW    R1,                     16(R12)
    MOVW    4(R12),         R0
    MOVW    12(R12),        R1
    EOR             R0,                     R1
    MOVW    R1,                     20(R12)

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

	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[18%5]
	SUB.S	R11,		R0
	MOVW	28(R12),		R11
	SBC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(18+1)%5]
	SUB.S	R11,		R2
	MOVW	36(R12),	R11
	SBC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(18+2)%5]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(18+3)%5]
	SUB.S	R11,		R6
	MOVW	12(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[18%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(18+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$18,			R6		// state[3] += 18 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[17%5]
	SUB.S	R11,		R0
	MOVW	20(R12),		R11
	SBC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(17+1)%5]
	SUB.S	R11,		R2
	MOVW	28(R12),	R11
	SBC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(17+2)%5]
	SUB.S	R11,		R4
	MOVW	36(R12),	R11
	SBC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(17+3)%5]
	SUB.S	R11,		R6
	MOVW	4(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[17%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(17+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$17,			R6		// state[3] += 17 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[16%5]
	SUB.S	R11,		R0
	MOVW	12(R12),		R11
	SBC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(16+1)%5]
	SUB.S	R11,		R2
	MOVW	20(R12),	R11
	SBC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(16+2)%5]
	SUB.S	R11,		R4
	MOVW	28(R12),	R11
	SBC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(16+3)%5]
	SUB.S	R11,		R6
	MOVW	36(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[16%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(16+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$16,			R6		// state[3] += 16 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[15%5]
	SUB.S	R11,		R0
	MOVW	4(R12),		R11
	SBC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(15+1)%5]
	SUB.S	R11,		R2
	MOVW	12(R12),	R11
	SBC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(15+2)%5]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(15+3)%5]
	SUB.S	R11,		R6
	MOVW	28(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[15%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(15+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$15,			R6		// state[3] += 15 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[14%5]
	SUB.S	R11,		R0
	MOVW	36(R12),		R11
	SBC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(14+1)%5]
	SUB.S	R11,		R2
	MOVW	4(R12),	R11
	SBC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(14+2)%5]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(14+3)%5]
	SUB.S	R11,		R6
	MOVW	20(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[14%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(14+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$14,			R6		// state[3] += 14 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[13%5]
	SUB.S	R11,		R0
	MOVW	28(R12),		R11
	SBC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(13+1)%5]
	SUB.S	R11,		R2
	MOVW	36(R12),	R11
	SBC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(13+2)%5]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(13+3)%5]
	SUB.S	R11,		R6
	MOVW	12(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[13%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(13+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$13,			R6		// state[3] += 13 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[12%5]
	SUB.S	R11,		R0
	MOVW	20(R12),		R11
	SBC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(12+1)%5]
	SUB.S	R11,		R2
	MOVW	28(R12),	R11
	SBC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(12+2)%5]
	SUB.S	R11,		R4
	MOVW	36(R12),	R11
	SBC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(12+3)%5]
	SUB.S	R11,		R6
	MOVW	4(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[12%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(12+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$12,			R6		// state[3] += 12 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[11%5]
	SUB.S	R11,		R0
	MOVW	12(R12),		R11
	SBC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(11+1)%5]
	SUB.S	R11,		R2
	MOVW	20(R12),	R11
	SBC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(11+2)%5]
	SUB.S	R11,		R4
	MOVW	28(R12),	R11
	SBC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(11+3)%5]
	SUB.S	R11,		R6
	MOVW	36(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[11%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(11+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$11,			R6		// state[3] += 11 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[10%5]
	SUB.S	R11,		R0
	MOVW	4(R12),		R11
	SBC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(10+1)%5]
	SUB.S	R11,		R2
	MOVW	12(R12),	R11
	SBC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(10+2)%5]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(10+3)%5]
	SUB.S	R11,		R6
	MOVW	28(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[10%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(10+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$10,			R6		// state[3] += 10 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[9%5]
	SUB.S	R11,		R0
	MOVW	36(R12),		R11
	SBC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(9+1)%5]
	SUB.S	R11,		R2
	MOVW	4(R12),	R11
	SBC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(9+2)%5]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(9+3)%5]
	SUB.S	R11,		R6
	MOVW	20(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[9%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(9+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$9,			R6		// state[3] += 9 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[8%5]
	SUB.S	R11,		R0
	MOVW	28(R12),		R11
	SBC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(8+1)%5]
	SUB.S	R11,		R2
	MOVW	36(R12),	R11
	SBC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(8+2)%5]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(8+3)%5]
	SUB.S	R11,		R6
	MOVW	12(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[8%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(8+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$8,			R6		// state[3] += 8 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[7%5]
	SUB.S	R11,		R0
	MOVW	20(R12),		R11
	SBC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(7+1)%5]
	SUB.S	R11,		R2
	MOVW	28(R12),	R11
	SBC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(7+2)%5]
	SUB.S	R11,		R4
	MOVW	36(R12),	R11
	SBC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(7+3)%5]
	SUB.S	R11,		R6
	MOVW	4(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[7%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(7+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$7,			R6		// state[3] += 7 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[6%5]
	SUB.S	R11,		R0
	MOVW	12(R12),		R11
	SBC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(6+1)%5]
	SUB.S	R11,		R2
	MOVW	20(R12),	R11
	SBC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(6+2)%5]
	SUB.S	R11,		R4
	MOVW	28(R12),	R11
	SBC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(6+3)%5]
	SUB.S	R11,		R6
	MOVW	36(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[6%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(6+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$6,			R6		// state[3] += 6 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[5%5]
	SUB.S	R11,		R0
	MOVW	4(R12),		R11
	SBC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(5+1)%5]
	SUB.S	R11,		R2
	MOVW	12(R12),	R11
	SBC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(5+2)%5]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(5+3)%5]
	SUB.S	R11,		R6
	MOVW	28(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[5%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(5+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$5,			R6		// state[3] += 5 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	32(R12),		R11		// state[0] += key[4%5]
	SUB.S	R11,		R0
	MOVW	36(R12),		R11
	SBC		R11,		R1
	MOVW	0(R12),		R11		// state[1] += key[(4+1)%5]
	SUB.S	R11,		R2
	MOVW	4(R12),	R11
	SBC		R11,		R3
	MOVW	8(R12),	R11		// state[2] += key[(4+2)%5]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5
	MOVW	16(R12),	R11		// state[3] += key[(4+3)%5]
	SUB.S	R11,		R6
	MOVW	20(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[4%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(4+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$4,			R6		// state[3] += 4 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	24(R12),		R11		// state[0] += key[3%5]
	SUB.S	R11,		R0
	MOVW	28(R12),		R11
	SBC		R11,		R1
	MOVW	32(R12),		R11		// state[1] += key[(3+1)%5]
	SUB.S	R11,		R2
	MOVW	36(R12),	R11
	SBC		R11,		R3
	MOVW	0(R12),	R11		// state[2] += key[(3+2)%5]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5
	MOVW	8(R12),	R11		// state[3] += key[(3+3)%5]
	SUB.S	R11,		R6
	MOVW	12(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[3%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(3+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$3,			R6		// state[3] += 3 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	16(R12),		R11		// state[0] += key[2%5]
	SUB.S	R11,		R0
	MOVW	20(R12),		R11
	SBC		R11,		R1
	MOVW	24(R12),		R11		// state[1] += key[(2+1)%5]
	SUB.S	R11,		R2
	MOVW	28(R12),	R11
	SBC		R11,		R3
	MOVW	32(R12),	R11		// state[2] += key[(2+2)%5]
	SUB.S	R11,		R4
	MOVW	36(R12),	R11
	SBC		R11,		R5
	MOVW	0(R12),	R11		// state[3] += key[(2+3)%5]
	SUB.S	R11,		R6
	MOVW	4(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	16(R12),		R11		// state[1] += tweak[2%3]
	SUB.S	R11,		R2
	MOVW	20(R12),		R11
	SBC		R11,		R3
	MOVW	0(R12),		R11		// state[2] += tweak[(2+1)%3]
	SUB.S	R11,		R4
	MOVW	4(R12),	R11
	SBC		R11,		R5

	SUB.S	$2,			R6		// state[3] += 2 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3


	MOVW	R2,	R11						// for rot==32 we can just swap
	MOVW	R3,	R2
	MOVW	R11,	R3


	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7


	MOVW	R6,	R11						// for rot==32 we can just swap
	MOVW	R7,	R6
	MOVW	R11,	R7


	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 10
	MOVW	R6>>22,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<10,	R12
	ORR		R12,	R11
	MOVW	R7>>22,	R7
	ORR		R6<<10,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3

	// Low0: rot << 6
	MOVW	R2<<6,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>26,	R12
	ORR		R12,	R11
	MOVW	R3<<6,	R3
	ORR		R2>>26,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3



	// High1: rot >> 20
	MOVW	R2>>12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<20,	R12
	ORR		R12,	R11
	MOVW	R3>>12,	R3
	ORR		R2<<20,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 18
	MOVW	R6<<18,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>14,	R12
	ORR		R12,	R11
	MOVW	R7<<18,	R7
	ORR		R6>>14,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 31
	MOVW	R6<<31,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>1,	R12
	ORR		R12,	R11
	MOVW	R7<<31,	R7
	ORR	R6>>1,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 7
	MOVW	R2>>25,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<7,	R12
	ORR		R12,	R11
	MOVW	R3>>25,	R3
	ORR		R2<<7,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	8(R12),		R11		// state[0] += key[1%5]
	SUB.S	R11,		R0
	MOVW	12(R12),		R11
	SBC		R11,		R1
	MOVW	16(R12),		R11		// state[1] += key[(1+1)%5]
	SUB.S	R11,		R2
	MOVW	20(R12),	R11
	SBC		R11,		R3
	MOVW	24(R12),	R11		// state[2] += key[(1+2)%5]
	SUB.S	R11,		R4
	MOVW	28(R12),	R11
	SBC		R11,		R5
	MOVW	32(R12),	R11		// state[3] += key[(1+3)%5]
	SUB.S	R11,		R6
	MOVW	36(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	8(R12),		R11		// state[1] += tweak[1%3]
	SUB.S	R11,		R2
	MOVW	12(R12),		R11
	SBC		R11,		R3
	MOVW	16(R12),		R11		// state[2] += tweak[(1+1)%3]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5

	SUB.S	$1,			R6		// state[3] += 1 (round number)
	SBC		$0,			R7

	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 27
	MOVW	R2<<27,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>5,	R12
	ORR		R12,	R11
	MOVW	R3<<27,	R3
	ORR	R2>>5,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7



	// High0: rot >> 27
	MOVW	R6>>5,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<27,	R12
	ORR		R12,	R11
	MOVW	R7>>5,	R7
	ORR		R6<<27,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7

	// Low1: rot << 24
	MOVW	R6<<24,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>8,	R12
	ORR		R12,	R11
	MOVW	R7<<24,	R7
	ORR	R6>>8,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 9
	MOVW	R2>>23,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<9,	R12
	ORR		R12,	R11
	MOVW	R3>>23,	R3
	ORR		R2<<9,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R2			// y1 = y1 ^ y0
	EOR		R5,		R3

	// Low1: rot << 7
	MOVW	R2<<7,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3>>25,	R12
	ORR		R12,	R11
	MOVW	R3<<7,	R3
	ORR	R2>>25,	R3
	MOVW	R11,	R2



	SUB.S	R2,		R4			// y0 = x0 - x1
	SBC		R3,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R6			// y1 = y1 ^ y0
	EOR		R1,		R7

	// Low0: rot << 12
	MOVW	R6<<12,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7>>20,	R12
	ORR		R12,	R11
	MOVW	R7<<12,	R7
	ORR		R6>>20,	R7
	MOVW	R11,	R6



	SUB.S	R6,		R0			// y0 = x0 - x1
	SBC		R7,		R1
		
	// Mix state[2] and state[3]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R4,		R6			// y1 = y1 ^ y0
	EOR		R5,		R7



	// High1: rot >> 16
	MOVW	R6>>16,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R7<<16,	R12
	ORR		R12,	R11
	MOVW	R7>>16,	R7
	ORR		R6<<16,	R7
	MOVW	R11,	R6

	SUB.S	R6,		R4			// y0 = x0 - x1
	SBC		R7,		R5

	// Mix state[0] and state[1]
	// y1 = y1 ^ y0
	// y1 = (x1 << r) | (x1 >> (64 - r))
	// y0 = x0 - x1
	EOR		R0,		R2			// y1 = y1 ^ y0
	EOR		R1,		R3



	// High0: rot >> 18
	MOVW	R2>>14,	R11			// y1 = (x1 << r) | (x1 >> (64 - r))
	MOVW	R3<<18,	R12
	ORR		R12,	R11
	MOVW	R3>>14,	R3
	ORR		R2<<18,	R3
	MOVW	R11,	R2

	SUB.S	R2,		R0			// y0 = x0 - x1
	SBC		R3,		R1
		
	// Key Schedule
	MOVW	key+4(FP),	R12
	MOVW	0(R12),		R11		// state[0] += key[0%5]
	SUB.S	R11,		R0
	MOVW	4(R12),		R11
	SBC		R11,		R1
	MOVW	8(R12),		R11		// state[1] += key[(0+1)%5]
	SUB.S	R11,		R2
	MOVW	12(R12),	R11
	SBC		R11,		R3
	MOVW	16(R12),	R11		// state[2] += key[(0+2)%5]
	SUB.S	R11,		R4
	MOVW	20(R12),	R11
	SBC		R11,		R5
	MOVW	24(R12),	R11		// state[3] += key[(0+3)%5]
	SUB.S	R11,		R6
	MOVW	28(R12),	R11
	SBC		R11,		R7

	MOVW	tweak+8(FP),R12
	MOVW	0(R12),		R11		// state[1] += tweak[0%3]
	SUB.S	R11,		R2
	MOVW	4(R12),		R11
	SBC		R11,		R3
	MOVW	8(R12),		R11		// state[2] += tweak[(0+1)%3]
	SUB.S	R11,		R4
	MOVW	12(R12),	R11
	SBC		R11,		R5

	SUB.S	$0,			R6		// state[3] += 0 (round number)
	SBC		$0,			R7


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

