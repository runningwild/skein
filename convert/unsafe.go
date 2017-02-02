package convert

import (
	"reflect"
	"unsafe"
)

// Does this work on big endian systems?  I don't think so...

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

// Inplace64BytesToUInt64 returns a pointer to a [8]uint64 that covers the exact bytes in b.
func Inplace64BytesToUInt64(b []byte) *[8]uint64 {
	if len(b) != 64 {
		return nil
	}
	return (*[8]uint64)(unsafe.Pointer(&b[0]))
}

// Inplace128BytesToUInt64 returns a pointer to a [16]uint64 that covers the exact bytes in b.
func Inplace128BytesToUInt64(b []byte) *[16]uint64 {
	if len(b) != 128 {
		return nil
	}
	return (*[16]uint64)(unsafe.Pointer(&b[0]))
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

// InplaceBytesToUInt64 returns a slice of uint64 that covers the exact data in v.
func InplaceBytesToUInt64(b []byte) []uint64 {
	if len(b)%8 != 0 {
		return nil
	}
	a := reflect.NewAt(reflect.ArrayOf(len(b)/8, reflect.TypeOf(uint64(0))), unsafe.Pointer(&b[0]))
	return a.Elem().Slice(0, len(b)/8).Interface().([]uint64)
}
