package bloomFilter

import (
	"fmt"
	"hash"
	"log"
)

// BloomFilter 布隆过滤器
type BloomFilter struct {
	// 位图
	BitSet BitSet
	// 哈希方法
	Hashes []hash.Hash64
}

// NewBloomFilter New
func NewBloomFilter(byteSize uint64, hashes ...hash.Hash64) *BloomFilter {
	return &BloomFilter{BitSet: NewBitSet(byteSize), Hashes: hashes}
}

func DefaultHashes(size int) []hash.Hash64 {
	hashes := []hash.Hash64{
		Hash64, MD5, SHA1, SHA256, SHA512,
	}
	if !(size > 0 && size <= len(hashes)) {
		panic("hashes size long")
	}

	return hashes[:size]
}

// Push 添加数据
func (f *BloomFilter) Push(content []byte) {
	var byteLen = f.BitSet.Size()
	if byteLen < 1 {
		fmt.Printf("da")

		return
	}

	for _, h := range f.Hashes {
		h.Reset()

		_, err := h.Write(content)
		if err != nil {
			log.Print(err.Error())
		}

		var res = h.Sum64()
		// Get the byte.
		var yByte = res % byteLen
		// Get the bit position in the byte.
		var (
			yBit = (res / byteLen) & 7
			now  = f.BitSet[yByte] | 1<<yBit
		)

		if now != f.BitSet[yByte] {
			f.BitSet[yByte] = now
		}
	}
}

// IsExists 判断数据是否存在
func (f *BloomFilter) IsExists(content []byte) bool {
	var byteLen = f.BitSet.Size()
	for _, h := range f.Hashes {
		h.Reset()

		_, err := h.Write(content)
		if err != nil {
			log.Print(err.Error())
		}

		var (
			res   = h.Sum64()
			yByte = res % byteLen
			yBit  = (res / byteLen) & 7
		)

		if f.BitSet[yByte]|1<<yBit != f.BitSet[yByte] {
			return false
		}
	}

	return true
}

func (f *BloomFilter) Reset() {
	f.BitSet.Reset()
}
