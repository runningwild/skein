package threefish

func encrypt512Simple(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	key[8] = c240 ^ key[0] ^ key[1] ^ key[2] ^ key[3] ^ key[4] ^ key[5] ^ key[6] ^ key[7]
	tweak[2] = tweak[0] ^ tweak[1]
	s0 := state[0]
	s1 := state[1]
	s2 := state[2]
	s3 := state[3]
	s4 := state[4]
	s5 := state[5]
	s6 := state[6]
	s7 := state[7]

	// Round 0

	// Key Schedule
	s0 += key[0]
	s1 += key[1]
	s2 += key[2]
	s3 += key[3]
	s4 += key[4]
	s5 += key[5] + tweak[0]
	s6 += key[6] + tweak[1]
	s7 += key[7] + 0

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 1

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 2

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 3

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 4

	// Key Schedule
	s0 += key[1]
	s1 += key[2]
	s2 += key[3]
	s3 += key[4]
	s4 += key[5]
	s5 += key[6] + tweak[1]
	s6 += key[7] + tweak[2]
	s7 += key[8] + 1

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 5

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 6

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 7

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 8

	// Key Schedule
	s0 += key[2]
	s1 += key[3]
	s2 += key[4]
	s3 += key[5]
	s4 += key[6]
	s5 += key[7] + tweak[2]
	s6 += key[8] + tweak[0]
	s7 += key[0] + 2

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 9

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 10

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 11

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 12

	// Key Schedule
	s0 += key[3]
	s1 += key[4]
	s2 += key[5]
	s3 += key[6]
	s4 += key[7]
	s5 += key[8] + tweak[0]
	s6 += key[0] + tweak[1]
	s7 += key[1] + 3

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 13

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 14

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 15

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 16

	// Key Schedule
	s0 += key[4]
	s1 += key[5]
	s2 += key[6]
	s3 += key[7]
	s4 += key[8]
	s5 += key[0] + tweak[1]
	s6 += key[1] + tweak[2]
	s7 += key[2] + 4

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 17

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 18

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 19

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 20

	// Key Schedule
	s0 += key[5]
	s1 += key[6]
	s2 += key[7]
	s3 += key[8]
	s4 += key[0]
	s5 += key[1] + tweak[2]
	s6 += key[2] + tweak[0]
	s7 += key[3] + 5

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 21

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 22

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 23

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 24

	// Key Schedule
	s0 += key[6]
	s1 += key[7]
	s2 += key[8]
	s3 += key[0]
	s4 += key[1]
	s5 += key[2] + tweak[0]
	s6 += key[3] + tweak[1]
	s7 += key[4] + 6

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 25

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 26

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 27

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 28

	// Key Schedule
	s0 += key[7]
	s1 += key[8]
	s2 += key[0]
	s3 += key[1]
	s4 += key[2]
	s5 += key[3] + tweak[1]
	s6 += key[4] + tweak[2]
	s7 += key[5] + 7

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 29

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 30

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 31

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 32

	// Key Schedule
	s0 += key[8]
	s1 += key[0]
	s2 += key[1]
	s3 += key[2]
	s4 += key[3]
	s5 += key[4] + tweak[2]
	s6 += key[5] + tweak[0]
	s7 += key[6] + 8

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 33

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 34

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 35

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 36

	// Key Schedule
	s0 += key[0]
	s1 += key[1]
	s2 += key[2]
	s3 += key[3]
	s4 += key[4]
	s5 += key[5] + tweak[0]
	s6 += key[6] + tweak[1]
	s7 += key[7] + 9

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 37

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 38

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 39

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 40

	// Key Schedule
	s0 += key[1]
	s1 += key[2]
	s2 += key[3]
	s3 += key[4]
	s4 += key[5]
	s5 += key[6] + tweak[1]
	s6 += key[7] + tweak[2]
	s7 += key[8] + 10

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 41

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 42

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 43

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 44

	// Key Schedule
	s0 += key[2]
	s1 += key[3]
	s2 += key[4]
	s3 += key[5]
	s4 += key[6]
	s5 += key[7] + tweak[2]
	s6 += key[8] + tweak[0]
	s7 += key[0] + 11

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 45

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 46

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 47

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 48

	// Key Schedule
	s0 += key[3]
	s1 += key[4]
	s2 += key[5]
	s3 += key[6]
	s4 += key[7]
	s5 += key[8] + tweak[0]
	s6 += key[0] + tweak[1]
	s7 += key[1] + 12

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 49

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 50

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 51

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 52

	// Key Schedule
	s0 += key[4]
	s1 += key[5]
	s2 += key[6]
	s3 += key[7]
	s4 += key[8]
	s5 += key[0] + tweak[1]
	s6 += key[1] + tweak[2]
	s7 += key[2] + 13

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 53

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 54

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 55

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 56

	// Key Schedule
	s0 += key[5]
	s1 += key[6]
	s2 += key[7]
	s3 += key[8]
	s4 += key[0]
	s5 += key[1] + tweak[2]
	s6 += key[2] + tweak[0]
	s7 += key[3] + 14

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 57

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 58

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 59

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 60

	// Key Schedule
	s0 += key[6]
	s1 += key[7]
	s2 += key[8]
	s3 += key[0]
	s4 += key[1]
	s5 += key[2] + tweak[0]
	s6 += key[3] + tweak[1]
	s7 += key[4] + 15

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 61

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 62

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 63

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 64

	// Key Schedule
	s0 += key[7]
	s1 += key[8]
	s2 += key[0]
	s3 += key[1]
	s4 += key[2]
	s5 += key[3] + tweak[1]
	s6 += key[4] + tweak[2]
	s7 += key[5] + 16

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 46) | (s1 >> (64 - 46)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 36) | (s3 >> (64 - 36)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 19) | (s5 >> (64 - 19)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 37) | (s7 >> (64 - 37)) ^ s6

	// Round 65

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 33) | (s1 >> (64 - 33)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 27) | (s7 >> (64 - 27)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 14) | (s5 >> (64 - 14)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 42) | (s3 >> (64 - 42)) ^ s0

	// Round 66

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 17) | (s1 >> (64 - 17)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 49) | (s3 >> (64 - 49)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 36) | (s5 >> (64 - 36)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 39) | (s7 >> (64 - 39)) ^ s2

	// Round 67

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 44) | (s1 >> (64 - 44)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 9) | (s7 >> (64 - 9)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 54) | (s5 >> (64 - 54)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 56) | (s3 >> (64 - 56)) ^ s4

	// Round 68

	// Key Schedule
	s0 += key[8]
	s1 += key[0]
	s2 += key[1]
	s3 += key[2]
	s4 += key[3]
	s5 += key[4] + tweak[2]
	s6 += key[5] + tweak[0]
	s7 += key[6] + 17

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 39) | (s1 >> (64 - 39)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 30) | (s3 >> (64 - 30)) ^ s2

	// Mix 4 with 5
	s4 += s5
	s5 = (s5 << 34) | (s5 >> (64 - 34)) ^ s4

	// Mix 6 with 7
	s6 += s7
	s7 = (s7 << 24) | (s7 >> (64 - 24)) ^ s6

	// Round 69

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 13) | (s1 >> (64 - 13)) ^ s2

	// Mix 4 with 7
	s4 += s7
	s7 = (s7 << 50) | (s7 >> (64 - 50)) ^ s4

	// Mix 6 with 5
	s6 += s5
	s5 = (s5 << 10) | (s5 >> (64 - 10)) ^ s6

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 17) | (s3 >> (64 - 17)) ^ s0

	// Round 70

	// Mix 4 with 1
	s4 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s4

	// Mix 6 with 3
	s6 += s3
	s3 = (s3 << 29) | (s3 >> (64 - 29)) ^ s6

	// Mix 0 with 5
	s0 += s5
	s5 = (s5 << 39) | (s5 >> (64 - 39)) ^ s0

	// Mix 2 with 7
	s2 += s7
	s7 = (s7 << 43) | (s7 >> (64 - 43)) ^ s2

	// Round 71

	// Mix 6 with 1
	s6 += s1
	s1 = (s1 << 8) | (s1 >> (64 - 8)) ^ s6

	// Mix 0 with 7
	s0 += s7
	s7 = (s7 << 35) | (s7 >> (64 - 35)) ^ s0

	// Mix 2 with 5
	s2 += s5
	s5 = (s5 << 56) | (s5 >> (64 - 56)) ^ s2

	// Mix 4 with 3
	s4 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s4

	// Round 72

	// Key Schedule
	s0 += key[0]
	s1 += key[1]
	s2 += key[2]
	s3 += key[3]
	s4 += key[4]
	s5 += key[5] + tweak[0]
	s6 += key[6] + tweak[1]
	s7 += key[7] + 18

	state[0] = s0
	state[1] = s1
	state[2] = s2
	state[3] = s3
	state[4] = s4
	state[5] = s5
	state[6] = s6
	state[7] = s7
}

func decrypt512Simple(state *[8]uint64, key *[9]uint64, tweak *[3]uint64) {
	key[8] = c240 ^ key[0] ^ key[1] ^ key[2] ^ key[3] ^ key[4] ^ key[5] ^ key[6] ^ key[7]
	tweak[2] = tweak[0] ^ tweak[1]
	s0 := state[0]
	s1 := state[1]
	s2 := state[2]
	s3 := state[3]
	s4 := state[4]
	s5 := state[5]
	s6 := state[6]
	s7 := state[7]

	// Round 72

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1]
	s2 -= key[2]
	s3 -= key[3]
	s4 -= key[4]
	s5 -= key[5] + tweak[0]
	s6 -= key[6] + tweak[1]
	s7 -= key[7] + 18

	// Round 71

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 70

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 69

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 68

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[8]
	s1 -= key[0]
	s2 -= key[1]
	s3 -= key[2]
	s4 -= key[3]
	s5 -= key[4] + tweak[2]
	s6 -= key[5] + tweak[0]
	s7 -= key[6] + 17

	// Round 67

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 66

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 65

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 64

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[7]
	s1 -= key[8]
	s2 -= key[0]
	s3 -= key[1]
	s4 -= key[2]
	s5 -= key[3] + tweak[1]
	s6 -= key[4] + tweak[2]
	s7 -= key[5] + 16

	// Round 63

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 62

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 61

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 60

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[6]
	s1 -= key[7]
	s2 -= key[8]
	s3 -= key[0]
	s4 -= key[1]
	s5 -= key[2] + tweak[0]
	s6 -= key[3] + tweak[1]
	s7 -= key[4] + 15

	// Round 59

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 58

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 57

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 56

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[5]
	s1 -= key[6]
	s2 -= key[7]
	s3 -= key[8]
	s4 -= key[0]
	s5 -= key[1] + tweak[2]
	s6 -= key[2] + tweak[0]
	s7 -= key[3] + 14

	// Round 55

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 54

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 53

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 52

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[4]
	s1 -= key[5]
	s2 -= key[6]
	s3 -= key[7]
	s4 -= key[8]
	s5 -= key[0] + tweak[1]
	s6 -= key[1] + tweak[2]
	s7 -= key[2] + 13

	// Round 51

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 50

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 49

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 48

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4]
	s2 -= key[5]
	s3 -= key[6]
	s4 -= key[7]
	s5 -= key[8] + tweak[0]
	s6 -= key[0] + tweak[1]
	s7 -= key[1] + 12

	// Round 47

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 46

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 45

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 44

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3]
	s2 -= key[4]
	s3 -= key[5]
	s4 -= key[6]
	s5 -= key[7] + tweak[2]
	s6 -= key[8] + tweak[0]
	s7 -= key[0] + 11

	// Round 43

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 42

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 41

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 40

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2]
	s2 -= key[3]
	s3 -= key[4]
	s4 -= key[5]
	s5 -= key[6] + tweak[1]
	s6 -= key[7] + tweak[2]
	s7 -= key[8] + 10

	// Round 39

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 38

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 37

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 36

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1]
	s2 -= key[2]
	s3 -= key[3]
	s4 -= key[4]
	s5 -= key[5] + tweak[0]
	s6 -= key[6] + tweak[1]
	s7 -= key[7] + 9

	// Round 35

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 34

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 33

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 32

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[8]
	s1 -= key[0]
	s2 -= key[1]
	s3 -= key[2]
	s4 -= key[3]
	s5 -= key[4] + tweak[2]
	s6 -= key[5] + tweak[0]
	s7 -= key[6] + 8

	// Round 31

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 30

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 29

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 28

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[7]
	s1 -= key[8]
	s2 -= key[0]
	s3 -= key[1]
	s4 -= key[2]
	s5 -= key[3] + tweak[1]
	s6 -= key[4] + tweak[2]
	s7 -= key[5] + 7

	// Round 27

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 26

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 25

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 24

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[6]
	s1 -= key[7]
	s2 -= key[8]
	s3 -= key[0]
	s4 -= key[1]
	s5 -= key[2] + tweak[0]
	s6 -= key[3] + tweak[1]
	s7 -= key[4] + 6

	// Round 23

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 22

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 21

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 20

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[5]
	s1 -= key[6]
	s2 -= key[7]
	s3 -= key[8]
	s4 -= key[0]
	s5 -= key[1] + tweak[2]
	s6 -= key[2] + tweak[0]
	s7 -= key[3] + 5

	// Round 19

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 18

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 17

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 16

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[4]
	s1 -= key[5]
	s2 -= key[6]
	s3 -= key[7]
	s4 -= key[8]
	s5 -= key[0] + tweak[1]
	s6 -= key[1] + tweak[2]
	s7 -= key[2] + 4

	// Round 15

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 14

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 13

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 12

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4]
	s2 -= key[5]
	s3 -= key[6]
	s4 -= key[7]
	s5 -= key[8] + tweak[0]
	s6 -= key[0] + tweak[1]
	s7 -= key[1] + 3

	// Round 11

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 10

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 9

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 8

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3]
	s2 -= key[4]
	s3 -= key[5]
	s4 -= key[6]
	s5 -= key[7] + tweak[2]
	s6 -= key[8] + tweak[0]
	s7 -= key[0] + 2

	// Round 7

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 56) | (s5 << (64 - 56))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 35) | (s7 << (64 - 35))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 8) | (s1 << (64 - 8))
	s6 -= s1

	// Round 6

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 43) | (s7 << (64 - 43))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 39) | (s5 << (64 - 39))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 29) | (s3 << (64 - 29))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s4 -= s1

	// Round 5

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 17) | (s3 << (64 - 17))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 10) | (s5 << (64 - 10))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 50) | (s7 << (64 - 50))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 13) | (s1 << (64 - 13))
	s2 -= s1

	// Round 4

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 24) | (s7 << (64 - 24))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 34) | (s5 << (64 - 34))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 30) | (s3 << (64 - 30))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 39) | (s1 << (64 - 39))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2]
	s2 -= key[3]
	s3 -= key[4]
	s4 -= key[5]
	s5 -= key[6] + tweak[1]
	s6 -= key[7] + tweak[2]
	s7 -= key[8] + 1

	// Round 3

	// Mix 4 with 3
	s3 ^= s4
	s3 = (s3 >> 56) | (s3 << (64 - 56))
	s4 -= s3

	// Mix 2 with 5
	s5 ^= s2
	s5 = (s5 >> 54) | (s5 << (64 - 54))
	s2 -= s5

	// Mix 0 with 7
	s7 ^= s0
	s7 = (s7 >> 9) | (s7 << (64 - 9))
	s0 -= s7

	// Mix 6 with 1
	s1 ^= s6
	s1 = (s1 >> 44) | (s1 << (64 - 44))
	s6 -= s1

	// Round 2

	// Mix 2 with 7
	s7 ^= s2
	s7 = (s7 >> 39) | (s7 << (64 - 39))
	s2 -= s7

	// Mix 0 with 5
	s5 ^= s0
	s5 = (s5 >> 36) | (s5 << (64 - 36))
	s0 -= s5

	// Mix 6 with 3
	s3 ^= s6
	s3 = (s3 >> 49) | (s3 << (64 - 49))
	s6 -= s3

	// Mix 4 with 1
	s1 ^= s4
	s1 = (s1 >> 17) | (s1 << (64 - 17))
	s4 -= s1

	// Round 1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 42) | (s3 << (64 - 42))
	s0 -= s3

	// Mix 6 with 5
	s5 ^= s6
	s5 = (s5 >> 14) | (s5 << (64 - 14))
	s6 -= s5

	// Mix 4 with 7
	s7 ^= s4
	s7 = (s7 >> 27) | (s7 << (64 - 27))
	s4 -= s7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 33) | (s1 << (64 - 33))
	s2 -= s1

	// Round 0

	// Mix 6 with 7
	s7 ^= s6
	s7 = (s7 >> 37) | (s7 << (64 - 37))
	s6 -= s7

	// Mix 4 with 5
	s5 ^= s4
	s5 = (s5 >> 19) | (s5 << (64 - 19))
	s4 -= s5

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 36) | (s3 << (64 - 36))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 46) | (s1 << (64 - 46))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1]
	s2 -= key[2]
	s3 -= key[3]
	s4 -= key[4]
	s5 -= key[5] + tweak[0]
	s6 -= key[6] + tweak[1]
	s7 -= key[7] + 0

	state[0] = s0
	state[1] = s1
	state[2] = s2
	state[3] = s3
	state[4] = s4
	state[5] = s5
	state[6] = s6
	state[7] = s7
}
