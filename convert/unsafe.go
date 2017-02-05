package convert

import (
	"reflect"
	"unsafe"
)

// Does this work on big endian systems?  I don't think so...

// Inplace8BytesToUInt64 returns a pointer to a [1]uint64 that covers the exact bytes in b.
func Inplace8BytesToUInt64(b []byte) *[1]uint64 {
	if len(b) != 8 {
		return nil
	}
	return (*[1]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace16BytesToUInt64 returns a pointer to a [2]uint64 that covers the exact bytes in b.
func Inplace16BytesToUInt64(b []byte) *[2]uint64 {
	if len(b) != 16 {
		return nil
	}
	return (*[2]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace32BytesToUInt64 returns a pointer to a [4]uint64 that covers the exact bytes in b.
func Inplace32BytesToUInt64(b []byte) *[4]uint64 {
	if len(b) != 32 {
		return nil
	}
	return (*[4]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace40BytesToUInt64 returns a pointer to a [5]uint64 that covers the exact bytes in b.
func Inplace40BytesToUInt64(b []byte) *[5]uint64 {
	if len(b) != 40 {
		return nil
	}
	return (*[5]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace64BytesToUInt64 returns a pointer to a [8]uint64 that covers the exact bytes in b.
func Inplace64BytesToUInt64(b []byte) *[8]uint64 {
	if len(b) != 64 {
		return nil
	}
	return (*[8]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace72BytesToUInt64 returns a pointer to a [9]uint64 that covers the exact bytes in b.
func Inplace72BytesToUInt64(b []byte) *[9]uint64 {
	if len(b) != 72 {
		return nil
	}
	return (*[9]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace128BytesToUInt64 returns a pointer to a [16]uint64 that covers the exact bytes in b.
func Inplace128BytesToUInt64(b []byte) *[16]uint64 {
	if len(b) != 128 {
		return nil
	}
	return (*[16]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace136BytesToUInt64 returns a pointer to a [17]uint64 that covers the exact bytes in b.
func Inplace136BytesToUInt64(b []byte) *[17]uint64 {
	if len(b) != 136 {
		return nil
	}
	return (*[17]uint64)(unsafe.Pointer(&b[0]))
}

// InplaceBytesToUInt64 returns a slice of uint64 that covers the exact data in v.
func InplaceBytesToUint64(b []byte) []uint64 {
	if len(b)%8 != 0 || len(b) == 0 {
		return nil
	}
	a := reflect.NewAt(reflect.ArrayOf(len(b)/8, reflect.TypeOf(uint64(0))), unsafe.Pointer(&b[0]))
	return a.Elem().Slice(0, len(b)/8).Interface().([]uint64)
}

// Inplace1Uint64ToBytes returns a pointer to a [8]byte that covers the exact data in v.
func Inplace1Uint64ToBytes(v []uint64) *[8]byte {
	if len(v) != 1 {
		return nil
	}
	return (*[8]byte)(unsafe.Pointer(&v[0]))
}

// Inplace2Uint64ToBytes returns a pointer to a [16]byte that covers the exact data in v.
func Inplace2Uint64ToBytes(v []uint64) *[16]byte {
	if len(v) != 2 {
		return nil
	}
	return (*[16]byte)(unsafe.Pointer(&v[0]))
}

// Inplace4Uint64ToBytes returns a pointer to a [32]byte that covers the exact data in v.
func Inplace4Uint64ToBytes(v []uint64) *[32]byte {
	if len(v) != 4 {
		return nil
	}
	return (*[32]byte)(unsafe.Pointer(&v[0]))
}

// Inplace8Uint64ToBytes returns a pointer to a [64]byte that covers the exact data in v.
func Inplace8Uint64ToBytes(v []uint64) *[64]byte {
	if len(v) != 8 {
		return nil
	}
	return (*[64]byte)(unsafe.Pointer(&v[0]))
}

// Inplace16Uint64ToBytes returns a pointer to a [128]byte that covers the exact data in v.
func Inplace16Uint64ToBytes(v []uint64) *[128]byte {
	if len(v) != 16 {
		return nil
	}
	return (*[128]byte)(unsafe.Pointer(&v[0]))
}

// InplaceUint64ToBytes returns a slice of uint64 that covers the exact data in v.
func InplaceUint64ToBytes(v []uint64) []byte {
	a := reflect.NewAt(reflect.ArrayOf(len(v)*8, reflect.TypeOf(byte(0))), unsafe.Pointer(&v[0]))
	return a.Elem().Slice(0, len(v)*8).Interface().([]byte)
}

func Xor(a, b, c []byte) {
	if len(a) != len(b) || len(b) != len(c) || len(a)&0x07 != 0 {
		panic("Xor requires all slices have the same length that is a multiple of 8")
	}
	for i := 0; i < len(a); i += 8 {
		*(*uint64)(unsafe.Pointer(&a[i])) = (*(*uint64)(unsafe.Pointer(&b[i]))) ^ (*(*uint64)(unsafe.Pointer(&c[i])))
	}
}
