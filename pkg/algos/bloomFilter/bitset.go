package bloomFilter

type BitSet []byte

func NewBitSet(size uint64) BitSet {
	return make([]byte, size)
}

// Reset 重置所有值为0
func (b BitSet) Reset() {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func (b BitSet) Size() uint64 {
	return uint64(len(b))
}

// IsEmpty 判断是否为空
func (b BitSet) IsEmpty() bool {
	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			return false
		}
	}

	return true
}
