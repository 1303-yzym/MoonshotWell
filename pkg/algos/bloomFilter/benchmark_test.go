package bloomFilter

import (
	"hash"
	"log"
	"testing"

	"github.com/twmb/murmur3"
)

func testHash(hash hash.Hash64) {
	hash.Reset()
	_, err := hash.Write([]byte("content"))
	if err != nil {
		log.Print(err)
	}
	hash.Sum64()
}

func BenchmarkHashFunc(b *testing.B) {
	benchmarks := []struct {
		name     string
		HashFunc hash.Hash64
	}{
		{
			name:     "Hash-murmur",
			HashFunc: murmur3.SeedNew64(1),
		},
		{
			name:     "Hash-MD5",
			HashFunc: NewMD5(),
		},
		{
			name:     "Hash-SHA1",
			HashFunc: NewSHA1(),
		},
		{
			name:     "SHA256",
			HashFunc: NewSHA256(),
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testHash(bm.HashFunc)
			}
		})
	}
}

func BenchmarkMurMur(b *testing.B) {
	hash64 := murmur3.New64()
	hash64.Write([]byte("Hello, World!"))
	for i := 0; i < b.N; i++ {
		_ = hash64.Sum64()
	}
}
