package ucrypto

import (
	"testing"
)

func BenchmarkName(b *testing.B) {
	str := "fwf"
	b.Run("V1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Md5V1(str)
		}
	})
	b.Run("V2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Md5V2(str)
		}
	})
	b.Run("V3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Md5V3(str)
		}
	})
}

func TestNewMd5(t *testing.T) {

}
