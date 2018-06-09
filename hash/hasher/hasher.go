package hasher

import (
	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/ubi"
)

func Hash(u *ubi.UBI, M []byte, lastByteBits int, N uint64) []byte {
	return u.Skein(nil, []ubi.Tuple{{ubi.TypeMsg, M, lastByteBits}}, N)
}

func TreeHash(u *ubi.UBI, M []byte, lastByteBits int, N uint64, Yl, Yf, Ym byte) []byte {
	return u.SkeinTree(nil, []ubi.Tuple{{ubi.TypeMsg, M, lastByteBits}}, N, Yl, Yf, Ym)
}

func MAC(u *ubi.UBI, K []byte, M []byte, lastByteBits int, N uint64) []byte {
	return u.Skein(K, []ubi.Tuple{{ubi.TypeMsg, M, lastByteBits}}, N)
}

func NewHasher(u *ubi.UBI, N int) *Hasher {
	h := &Hasher{
		ubi:        u,
		buf:        make([]byte, u.TBC().BlockSize()/8)[0:0],
		n:          uint64(N),
		blockBytes: u.TBC().BlockSize() / 8,
	}
	h.Reset()
	return h
}

func NewMACer(u *ubi.UBI, key []byte, N int) *Hasher {
	h := &Hasher{
		ubi:        u,
		buf:        make([]byte, u.TBC().BlockSize()/8)[0:0],
		key:        make([]byte, len(key)),
		n:          uint64(N),
		blockBytes: u.TBC().BlockSize() / 8,
	}
	copy(h.key, key)
	h.Reset()
	return h
}

type Hasher struct {
	ubi        *ubi.UBI
	buf        []byte
	key        []byte
	n          uint64
	blockBytes int
	it         *ubi.Iterator
}

func (h *Hasher) Write(data []byte) (n int, err error) {
	written := len(data)

	// If the buffer plus b won't fill up more than a full block, just add it to the buffer and return.
	if len(h.buf)+len(data) <= h.blockBytes {
		h.buf = append(h.buf, data...)
		return written, nil
	}

	// If we have something in the buffer we need to make a full block using that first.
	if len(h.buf) > 0 {
		amt := h.blockBytes - len(h.buf)
		h.it.Block(append(h.buf, data[0:amt]...))
		h.buf = h.buf[0:0]
		data = data[amt:]
	}

	// Iterate until we don't have more than a full block.
	for len(data) > h.blockBytes {
		h.it.Block(data[0:h.blockBytes])
		data = data[h.blockBytes:]
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

	buf := make([]byte, int(h.n)/8+h.blockBytes)
	view := buf[:]
	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (h.n + uint64(h.ubi.TBC().BlockSize()) - 1) / uint64(h.ubi.TBC().BlockSize())
	for c[0] < iterations {
		copy(view, h.ubi.UBI(Gn, cb, [2]uint64{0, uint64(ubi.TypeOut)}))
		view = view[h.blockBytes:]
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
	Gn := h.ubi.GetInitialChainingValue(h.key, h.n, 0, 0, 0)
	h.it = h.ubi.Iterate(Gn, [2]uint64{0, uint64(ubi.TypeMsg)})
	h.buf = h.buf[0:0]
}

func (h *Hasher) Size() int {
	return (int(h.n) + 7) / 8
}

func (h *Hasher) BlockSize() int {
	return h.ubi.TBC().BlockSize()
}
