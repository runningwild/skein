package ubi

import "github.com/runningwild/skein/convert"

type Hasher struct {
	ubi *UBI
	buf []byte
	key []byte
	n   uint64
	it  *Iterator
}

func (h *Hasher) Write(data []byte) (n int, err error) {
	written := len(data)

	// If the buffer plus b won't fill up more than a full block, just add it to the buffer and return.
	if len(h.buf)+len(data) <= h.ubi.blockBytes {
		h.buf = append(h.buf, data...)
		return written, nil
	}

	// If we have something in the buffer we need to make a full block using that first.
	if len(h.buf) > 0 {
		amt := h.ubi.blockBytes - len(h.buf)
		h.it.Block(append(h.buf, data[0:amt]...))
		h.buf = h.buf[0:0]
		data = data[amt:]
	}

	// Iterate until we don't have more than a full block.
	for len(data) > h.ubi.blockBytes {
		h.it.Block(data[0:h.ubi.blockBytes])
		data = data[h.ubi.blockBytes:]
	}

	// Store anything up to a single remaining full block for later.
	h.buf = append(h.buf, data...)

	return written, nil
}

func (h *Hasher) Sum(b []byte) []byte {
	return h.sumInternal(b, 0)
}

func (h *Hasher) SumBits(b []byte, lbb int) []byte {
	return h.sumInternal(b, lbb)
}

func (h *Hasher) sumInternal(b []byte, lbb int) []byte {
	Gn := h.it.Finish(h.buf, lbb)

	buf := make([]byte, int(h.n)/8+h.ubi.blockBytes)
	view := buf[:]
	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (h.n + uint64(h.ubi.blockSize) - 1) / uint64(h.ubi.blockSize)
	for c[0] < iterations {
		copy(view, h.ubi.UBI(Gn, cb, tweakTypeOut))
		view = view[h.ubi.blockBytes:]
		c[0]++
	}
	if uint64(len(buf)*8) > h.n {
		buf = buf[0 : int(h.n+7)/8]
		if h.n&0x07 != 0 {
			// This masks away the upper bits that we don't care about, in the event that we asked for a
			// number of bits that doesn't evenly divide a byte.
			buf[len(buf)-1] = buf[len(buf)-1] & ((1 << uint(h.n&0x07)) - 1)
		}
	}
	return append(b, buf...)
}

func (h *Hasher) Reset() {
	Gn := h.ubi.GetInitialChainingValue(h.key, h.n)
	h.it = h.ubi.Iterate(Gn, tweakTypeMsg)
	h.buf = h.buf[0:0]
}

func (h *Hasher) Size() int {
	return (int(h.n) + 7) / 8
}

func (h *Hasher) BlockSize() int {
	return h.ubi.blockSize
}
