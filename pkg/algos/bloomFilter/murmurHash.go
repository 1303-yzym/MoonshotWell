package bloomFilter

import (
	"hash"

	"github.com/duke-git/lancet/v2/random"
	"github.com/twmb/murmur3"
)

// GenMurmurHashFunctions
// 这是一个更快的hash函数
func GenMurmurHashFunctions(numFunctions int) []hash.Hash64 {
	var hashFunc []hash.Hash64
	for i := 0; i < numFunctions; i++ {
		hashFunc = append(hashFunc, murmur3.SeedNew64(uint64(random.RandInt(1, 99999))))
	}

	return hashFunc
}
