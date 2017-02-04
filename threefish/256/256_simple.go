package threefish

func encrypt256Simple(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	key[4] = c240 ^ key[0] ^ key[1] ^ key[2] ^ key[3]
	tweak[2] = tweak[0] ^ tweak[1]
	s0 := state[0]
	s1 := state[1]
	s2 := state[2]
	s3 := state[3]

	// Round 0

	// Key Schedule
	s0 += key[0]
	s1 += key[1] + tweak[0]
	s2 += key[2] + tweak[1]
	s3 += key[3] + 0

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 1

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 2

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 3

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 4

	// Key Schedule
	s0 += key[1]
	s1 += key[2] + tweak[1]
	s2 += key[3] + tweak[2]
	s3 += key[4] + 1

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 5

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 6

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 7

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 8

	// Key Schedule
	s0 += key[2]
	s1 += key[3] + tweak[2]
	s2 += key[4] + tweak[0]
	s3 += key[0] + 2

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 9

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 10

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 11

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 12

	// Key Schedule
	s0 += key[3]
	s1 += key[4] + tweak[0]
	s2 += key[0] + tweak[1]
	s3 += key[1] + 3

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 13

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 14

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 15

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 16

	// Key Schedule
	s0 += key[4]
	s1 += key[0] + tweak[1]
	s2 += key[1] + tweak[2]
	s3 += key[2] + 4

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 17

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 18

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 19

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 20

	// Key Schedule
	s0 += key[0]
	s1 += key[1] + tweak[2]
	s2 += key[2] + tweak[0]
	s3 += key[3] + 5

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 21

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 22

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 23

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 24

	// Key Schedule
	s0 += key[1]
	s1 += key[2] + tweak[0]
	s2 += key[3] + tweak[1]
	s3 += key[4] + 6

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 25

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 26

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 27

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 28

	// Key Schedule
	s0 += key[2]
	s1 += key[3] + tweak[1]
	s2 += key[4] + tweak[2]
	s3 += key[0] + 7

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 29

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 30

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 31

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 32

	// Key Schedule
	s0 += key[3]
	s1 += key[4] + tweak[2]
	s2 += key[0] + tweak[0]
	s3 += key[1] + 8

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 33

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 34

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 35

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 36

	// Key Schedule
	s0 += key[4]
	s1 += key[0] + tweak[0]
	s2 += key[1] + tweak[1]
	s3 += key[2] + 9

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 37

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 38

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 39

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 40

	// Key Schedule
	s0 += key[0]
	s1 += key[1] + tweak[1]
	s2 += key[2] + tweak[2]
	s3 += key[3] + 10

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 41

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 42

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 43

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 44

	// Key Schedule
	s0 += key[1]
	s1 += key[2] + tweak[2]
	s2 += key[3] + tweak[0]
	s3 += key[4] + 11

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 45

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 46

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 47

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 48

	// Key Schedule
	s0 += key[2]
	s1 += key[3] + tweak[0]
	s2 += key[4] + tweak[1]
	s3 += key[0] + 12

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 49

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 50

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 51

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 52

	// Key Schedule
	s0 += key[3]
	s1 += key[4] + tweak[1]
	s2 += key[0] + tweak[2]
	s3 += key[1] + 13

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 53

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 54

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 55

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 56

	// Key Schedule
	s0 += key[4]
	s1 += key[0] + tweak[2]
	s2 += key[1] + tweak[0]
	s3 += key[2] + 14

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 57

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 58

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 59

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 60

	// Key Schedule
	s0 += key[0]
	s1 += key[1] + tweak[0]
	s2 += key[2] + tweak[1]
	s3 += key[3] + 15

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 61

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 62

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 63

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 64

	// Key Schedule
	s0 += key[1]
	s1 += key[2] + tweak[1]
	s2 += key[3] + tweak[2]
	s3 += key[4] + 16

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 14) | (s1 >> (64 - 14)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 16) | (s3 >> (64 - 16)) ^ s2

	// Round 65

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 52) | (s3 >> (64 - 52)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 57) | (s1 >> (64 - 57)) ^ s2

	// Round 66

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 23) | (s1 >> (64 - 23)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 40) | (s3 >> (64 - 40)) ^ s2

	// Round 67

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 5) | (s3 >> (64 - 5)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 37) | (s1 >> (64 - 37)) ^ s2

	// Round 68

	// Key Schedule
	s0 += key[2]
	s1 += key[3] + tweak[2]
	s2 += key[4] + tweak[0]
	s3 += key[0] + 17

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 25) | (s1 >> (64 - 25)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 33) | (s3 >> (64 - 33)) ^ s2

	// Round 69

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 46) | (s3 >> (64 - 46)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 12) | (s1 >> (64 - 12)) ^ s2

	// Round 70

	// Mix 0 with 1
	s0 += s1
	s1 = (s1 << 58) | (s1 >> (64 - 58)) ^ s0

	// Mix 2 with 3
	s2 += s3
	s3 = (s3 << 22) | (s3 >> (64 - 22)) ^ s2

	// Round 71

	// Mix 0 with 3
	s0 += s3
	s3 = (s3 << 32) | (s3 >> (64 - 32)) ^ s0

	// Mix 2 with 1
	s2 += s1
	s1 = (s1 << 32) | (s1 >> (64 - 32)) ^ s2

	// Round 72

	// Key Schedule
	s0 += key[3]
	s1 += key[4] + tweak[0]
	s2 += key[0] + tweak[1]
	s3 += key[1] + 18

	state[0] = s0
	state[1] = s1
	state[2] = s2
	state[3] = s3
}

func decrypt256Simple(state *[4]uint64, key *[5]uint64, tweak *[3]uint64) {
	key[4] = c240 ^ key[0] ^ key[1] ^ key[2] ^ key[3]
	tweak[2] = tweak[0] ^ tweak[1]
	s0 := state[0]
	s1 := state[1]
	s2 := state[2]
	s3 := state[3]

	// Round 72

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4] + tweak[0]
	s2 -= key[0] + tweak[1]
	s3 -= key[1] + 18

	// Round 71

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 70

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 69

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 68

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3] + tweak[2]
	s2 -= key[4] + tweak[0]
	s3 -= key[0] + 17

	// Round 67

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 66

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 65

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 64

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2] + tweak[1]
	s2 -= key[3] + tweak[2]
	s3 -= key[4] + 16

	// Round 63

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 62

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 61

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 60

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1] + tweak[0]
	s2 -= key[2] + tweak[1]
	s3 -= key[3] + 15

	// Round 59

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 58

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 57

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 56

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[4]
	s1 -= key[0] + tweak[2]
	s2 -= key[1] + tweak[0]
	s3 -= key[2] + 14

	// Round 55

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 54

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 53

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 52

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4] + tweak[1]
	s2 -= key[0] + tweak[2]
	s3 -= key[1] + 13

	// Round 51

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 50

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 49

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 48

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3] + tweak[0]
	s2 -= key[4] + tweak[1]
	s3 -= key[0] + 12

	// Round 47

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 46

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 45

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 44

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2] + tweak[2]
	s2 -= key[3] + tweak[0]
	s3 -= key[4] + 11

	// Round 43

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 42

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 41

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 40

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1] + tweak[1]
	s2 -= key[2] + tweak[2]
	s3 -= key[3] + 10

	// Round 39

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 38

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 37

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 36

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[4]
	s1 -= key[0] + tweak[0]
	s2 -= key[1] + tweak[1]
	s3 -= key[2] + 9

	// Round 35

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 34

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 33

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 32

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4] + tweak[2]
	s2 -= key[0] + tweak[0]
	s3 -= key[1] + 8

	// Round 31

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 30

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 29

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 28

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3] + tweak[1]
	s2 -= key[4] + tweak[2]
	s3 -= key[0] + 7

	// Round 27

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 26

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 25

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 24

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2] + tweak[0]
	s2 -= key[3] + tweak[1]
	s3 -= key[4] + 6

	// Round 23

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 22

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 21

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 20

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1] + tweak[2]
	s2 -= key[2] + tweak[0]
	s3 -= key[3] + 5

	// Round 19

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 18

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 17

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 16

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[4]
	s1 -= key[0] + tweak[1]
	s2 -= key[1] + tweak[2]
	s3 -= key[2] + 4

	// Round 15

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 14

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 13

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 12

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[3]
	s1 -= key[4] + tweak[0]
	s2 -= key[0] + tweak[1]
	s3 -= key[1] + 3

	// Round 11

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 10

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 9

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 8

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[2]
	s1 -= key[3] + tweak[2]
	s2 -= key[4] + tweak[0]
	s3 -= key[0] + 2

	// Round 7

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 32) | (s1 << (64 - 32))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 32) | (s3 << (64 - 32))
	s0 -= s3

	// Round 6

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 22) | (s3 << (64 - 22))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 58) | (s1 << (64 - 58))
	s0 -= s1

	// Round 5

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 12) | (s1 << (64 - 12))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 46) | (s3 << (64 - 46))
	s0 -= s3

	// Round 4

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 33) | (s3 << (64 - 33))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 25) | (s1 << (64 - 25))
	s0 -= s1

	// Key Schedule
	s0 -= key[1]
	s1 -= key[2] + tweak[1]
	s2 -= key[3] + tweak[2]
	s3 -= key[4] + 1

	// Round 3

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 37) | (s1 << (64 - 37))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 5) | (s3 << (64 - 5))
	s0 -= s3

	// Round 2

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 40) | (s3 << (64 - 40))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 23) | (s1 << (64 - 23))
	s0 -= s1

	// Round 1

	// Mix 2 with 1
	s1 ^= s2
	s1 = (s1 >> 57) | (s1 << (64 - 57))
	s2 -= s1

	// Mix 0 with 3
	s3 ^= s0
	s3 = (s3 >> 52) | (s3 << (64 - 52))
	s0 -= s3

	// Round 0

	// Mix 2 with 3
	s3 ^= s2
	s3 = (s3 >> 16) | (s3 << (64 - 16))
	s2 -= s3

	// Mix 0 with 1
	s1 ^= s0
	s1 = (s1 >> 14) | (s1 << (64 - 14))
	s0 -= s1

	// Key Schedule
	s0 -= key[0]
	s1 -= key[1] + tweak[0]
	s2 -= key[2] + tweak[1]
	s3 -= key[3] + 0

	state[0] = s0
	state[1] = s1
	state[2] = s2
	state[3] = s3
}
