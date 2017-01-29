package skein

import "unsafe"

// inplaceCovertBytesToUInt64 returns a pointer to a [4]uint64 that covers the exact bytes in b.
func inplaceCovertBytesToUInt64(b []byte) *[4]uint64 {
	if len(b) != 32 {
		return nil
	}
	// Does this work on little endian systems?  I just don't know...
	return (*[4]uint64)(unsafe.Pointer(&b[0]))
}
