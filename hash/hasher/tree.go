package hasher

import (
	"bytes"
	"sync"

	"github.com/runningwild/skein/convert"
	"github.com/runningwild/skein/ubi"
)

func NewTreeHasher(u *ubi.UBI, N int, Yl, Yf, Ym byte) *TreeHasher {
	h := &TreeHasher{
		ubi:        u,
		gn:         u.GetInitialChainingValue(nil, uint64(N), Yl, Yf, Ym),
		n:          uint64(N),
		leafBytes:  u.TBC().BlockSize() / 8 * (1 << Yl),
		nodeBytes:  u.TBC().BlockSize() / 8 * (1 << Yf),
		blockBytes: u.TBC().BlockSize() / 8,
		yl:         Yl,
		yf:         Yf,
		ym:         Ym,
	}
	h.Reset()
	return h
}

type TreeHasher struct {
	ubi        *ubi.UBI
	levels     []*treeLevel
	gn         []byte
	key        []byte
	n          uint64
	blockBytes int
	leafBytes  int
	nodeBytes  int
	yl, yf, ym byte

	it *ubi.Iterator

	// rootIt is used exclusively for the root of the tree.
	rootIt *ubi.Iterator
}

type treeLevel struct {
	m     []byte
	tweak [2]uint64
	size  int
}

func (h *TreeHasher) addLevel() {
	h.levels = append(h.levels, &treeLevel{
		tweak: [2]uint64{0, uint64(ubi.TypeMsg) | (uint64(1+len(h.levels)) << (112 - 64))},
	})
	if len(h.levels) == 1 {
		h.levels[len(h.levels)-1].size = h.leafBytes
	} else {
		h.levels[len(h.levels)-1].size = h.nodeBytes
	}
}

// returns true if anything was compressed
func (h *TreeHasher) bubbleUp(level int) bool {
	if level == int(h.ym)-1 {
		// We can't bubble up past the max level, so just eat the data with the root iterator.
		if h.rootIt == nil {
			h.rootIt = h.ubi.Iterate(h.gn, [2]uint64{0, uint64(ubi.TypeMsg) | uint64(h.ym)<<(112-64)})
		}
		l := h.levels[level]
		// Only feed a block to the iterator if there is more than a full block remaining, if there
		// is exactly one block we will have no way of knowing that it is the final block and the
		// iterator will not be able to mark the appropriate flag in the tweak.
		for len(l.m) > h.blockBytes {
			h.rootIt.Block(l.m[0:h.blockBytes])
			l.m = l.m[h.blockBytes:]
		}
		return false
	}
	low := h.levels[level]
	if len(low.m) < low.size {
		return false
	}
	if len(h.levels) <= level+1 {
		h.addLevel()
	}
	high := h.levels[level+1]

	if true {
		// This is where the parallelism is for now.
		var blocks, outputs [][]byte
		for start := 0; start+low.size <= len(low.m); start += low.size {
			blocks = append(blocks, low.m[start:start+low.size])
			outputs = append(outputs, nil)
		}
		var wg sync.WaitGroup
		for i := range blocks {
			wg.Add(1)
			go func(index int, tweak [2]uint64) {
				defer wg.Done()
				outputs[index] = h.ubi.UBI(h.gn, blocks[i], tweak)
			}(i, low.tweak)
			low.tweak[0] += uint64(low.size)
		}
		wg.Wait()
		low.m = low.m[low.size*len(blocks):]
		high.m = bytes.Join(append([][]byte{high.m}, outputs...), nil)
	} else {
		for len(low.m) >= low.size {
			high.m = append(high.m, h.ubi.UBI(h.gn, low.m[0:low.size], low.tweak)...)
			low.m = low.m[low.size:]
			low.tweak[0] += uint64(low.size)
		}
	}
	return true
}

func (h *TreeHasher) Write(data []byte) (n int, err error) {
	written := len(data)

	// If the buffer plus b won't fill up more than a full leaf block, just add it to the buffer
	// and return.
	if len(h.levels[0].m)+len(data) <= h.leafBytes {
		h.levels[0].m = append(h.levels[0].m, data...)
		return written, nil
	}

	// If we have something in the buffer we need to make a full block using that first.
	if len(h.levels[0].m) > 0 {
		amt := h.levels[0].size - len(h.levels[0].m)
		h.levels[0].m = append(h.levels[0].m, data[0:amt]...)
		h.bubbleUp(0)
		data = data[amt:]
	}
	if len(h.levels[0].m) != 0 {
		panic("you screwed up!")
	}

	// Now we can bubble up the entire input message
	h.levels[0].m = data
	for cur := 0; h.bubbleUp(cur); cur++ {
		h.bubbleUp(cur)
	}

	// Store anything up to a single remaining full block for later.
	rem := make([]byte, len(h.levels[0].m))
	copy(rem, h.levels[0].m)
	h.levels[0].m = rem

	return written, nil
}

func (h *TreeHasher) Sum(b []byte) []byte {
	height := 0
	for height = len(h.levels) - 1; height > 0 && len(h.levels[height].m) == 0; height-- {
	}
	var carry []byte
	for i := 0; i <= height && i-1 < int(h.ym); i++ {
		if len(carry) > 0 {
			carry = h.ubi.UBI(h.gn, carry, h.levels[i-1].tweak)
		}
		buf := make([]byte, 0, len(h.levels[i].m)+len(carry))
		buf = append(buf, h.levels[i].m...)
		carry = append(buf, carry...)
	}
	var g []byte
	if h.rootIt != nil {
		rootIt := h.rootIt
		if len(carry) > h.blockBytes {
			rootIt = h.rootIt.Clone()
		}
		for len(carry) > h.blockBytes {
			rootIt.Block(carry[0:h.blockBytes])
			carry = carry[h.blockBytes:]
		}
		g = rootIt.Finish(carry, 0)

	} else if height > 0 && len(carry) == h.blockBytes {
		g = carry
	} else {
		tweak := [2]uint64{0, uint64(ubi.TypeMsg) | (uint64(len(h.levels)) << (112 - 64))}
		g = h.ubi.UBI(h.gn, carry, tweak)
	}
	buf := make([]byte, int(h.n)/8+h.blockBytes)
	view := buf[:]

	// put c in an array so we can convert it to bytes to pass to UBI.
	var c [1]uint64
	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
	iterations := (h.n + uint64(h.ubi.TBC().BlockSize()) - 1) / uint64(h.ubi.TBC().BlockSize())
	for c[0] < iterations {
		copy(view, h.ubi.UBI(g, cb, [2]uint64{0, uint64(ubi.TypeOut)}))
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

// func (h *TreeHasher) SumBits(b []byte, lbb int) []byte {
// 	return h.sumInternal(b, lbb)
// }

// func (h *TreeHasher) sumInternal(b []byte, lbb int) []byte {
// 	Gn := h.it.Finish(h.buf, lbb)

// 	buf := make([]byte, int(h.n)/8+h.blockBytes)
// 	view := buf[:]
// 	// put c in an array so we can convert it to bytes to pass to UBI.
// 	var c [1]uint64
// 	cb := convert.Inplace1Uint64ToBytes(c[:])[:]
// 	iterations := (h.n + uint64(h.ubi.TBC().BlockSize()) - 1) / uint64(h.ubi.TBC().BlockSize())
// 	for c[0] < iterations {
// 		copy(view, h.ubi.UBI(Gn, cb, [2]uint64{0, uint64(ubi.TypeOut)}))
// 		view = view[h.blockBytes:]
// 		c[0]++
// 	}
// 	if uint64(len(buf)*8) > h.n {
// 		buf = buf[0 : int(h.n+7)/8]
// 		if h.n&0x07 != 0 {
// 			// This masks away the upper bits that we don't care about, in the event that we asked for a
// 			// number of bits that doesn't evenly divide a byte.
// 			buf[len(buf)-1] = buf[len(buf)-1] & ((1 << uint(h.n&0x07)) - 1)
// 		}
// 	}
// 	return append(b, buf...)
// }

func (h *TreeHasher) Reset() {
	h.it = h.ubi.Iterate(h.gn, [2]uint64{0, uint64(ubi.TypeMsg)})
	h.rootIt = nil
	h.levels = nil
	h.addLevel()
}

func (h *TreeHasher) Size() int {
	return (int(h.n) + 7) / 8
}

func (h *TreeHasher) BlockSize() int {
	return h.ubi.TBC().BlockSize()
}
