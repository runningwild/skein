package skein

import "unsafe"

// Does this work on little endian systems?  I just don't know...

// inplaceConvert16BytesToUInt64 returns a pointer to a [2]uint64 that covers the exact bytes in b.
func inplaceConvert16BytesToUInt64(b []byte) *[2]uint64 {
	if len(b) != 16 {
		return nil
	}
	return (*[2]uint64)(unsafe.Pointer(&b[0]))
}

// inplaceConvert32BytesToUInt64 returns a pointer to a [4]uint64 that covers the exact bytes in b.
func inplaceConvert32BytesToUInt64(b []byte) *[4]uint64 {
	if len(b) != 32 {
		return nil
	}
	return (*[4]uint64)(unsafe.Pointer(&b[0]))
}

// inplaceConvert64BytesToUInt64 returns a pointer to a [8]uint64 that covers the exact bytes in b.
func inplaceConvert64BytesToUInt64(b []byte) *[8]uint64 {
	if len(b) != 64 {
		return nil
	}
	return (*[8]uint64)(unsafe.Pointer(&b[0]))
}

// inplaceConvert128BytesToUInt64 returns a pointer to a [16]uint64 that covers the exact bytes in b.
func inplaceConvert128BytesToUInt64(b []byte) *[16]uint64 {
	if len(b) != 128 {
		return nil
	}
	return (*[16]uint64)(unsafe.Pointer(&b[0]))
}
