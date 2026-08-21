package bloomFilter

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// 获取期望的存储数量和错误率，需要的位图大小
func TestExpectationsBitmapSize(t *testing.T) {
	elSize := 100000000
	errSize := 0.00001
	byteSet := ExpectationsBitmapSize(elSize, errSize)
	t.Logf("byteSet:%d,memSzie:%dkb,%.3fmb", byteSet, byteSet/1024, float64(byteSet)/1024/1024)
}

func TestNewBloomFilter(t *testing.T) {
	bloomFilter := NewBloomFilter(10240*8, DefaultHashes(4)...)
	bloomFilter.Push([]byte("da"))
	if bloomFilter.BitSet.IsEmpty() {
		t.Log("IsNil")
	}
	t.Log(bloomFilter.IsExists([]byte("dDAa")))
}

func TestMurmurHash(t *testing.T) {
	bloomFilter := NewBloomFilter(10240*8, GenMurmurHashFunctions(4)...)
	bloomFilter.Push([]byte("da"))
	if bloomFilter.BitSet.IsEmpty() {
		t.Log("IsNil")
	}
	t.Log(bloomFilter.IsExists([]byte("da")))
}

// fill some nums into f
func fillNums(filter *BloomFilter, begin int, end int) {
	for i := begin; i < end+1; i++ {
		filter.Push([]byte("id" + strconv.Itoa(i)))
	}
	fmt.Printf("已填入%d-%d的数据\n", begin, end)
}

func TestMemoryFilter(t *testing.T) {
	// Initial a memory filter mem page 4096
	memFilter := NewBloomFilter(1024*4*32, GenMurmurHashFunctions(6)...)
	// Push 0-9999 numbers to the filter.
	fillNums(memFilter, 0, 99999)
	// Check whether 10000-1000000 exist in the filter or not.
	// 查看10000-1000000 是否存在于过滤器中
	memFilter.Push([]byte(strconv.Itoa(10000)))
	var counter = 0
	for i := 100000; i < 10000000; i++ {
		if memFilter.IsExists([]byte(strconv.Itoa(i))) {
			counter++
			t.Log("hashed:", i)
		}
	}
	t.Log("碰撞个数：", counter)
}

func getRandomNums(top int, quantity int) map[int]struct{} {
	rand.Seed(time.Now().Unix())
	var m = make(map[int]struct{})
	for i := 0; i < quantity; i++ {
		r := rand.Intn(top)
		_, exists := m[r]
		for exists {
			r = rand.Intn(top)
			_, exists = m[r]
		}
		m[r] = struct{}{}
	}
	return m
}

func TestFalsePositiveRate(t *testing.T) {
	bitsCount := 1024 * 4 * 32
	existedCount := 100000
	hashCount := 6
	top := 100000000
	bloomFilter := NewBloomFilter(uint64(bitsCount), GenMurmurHashFunctions(hashCount)...)
	rate := GetFalsePositiveRate(uint(bitsCount), uint(existedCount), uint(hashCount))
	t.Log("理论误判率 Theoretical false positive rate:", rate)
	nums := getRandomNums(top, existedCount)
	for m := range nums {
		bloomFilter.Push([]byte(strconv.Itoa(m)))
	}
	count := 0
	for i := 0; i < top; i++ {
		if bloomFilter.IsExists([]byte(strconv.Itoa(i))) {
			count++
		}
	}
	realRate := float64(count-existedCount) / float64(top-existedCount)
	t.Logf("实际误判率 Real false positive rate:%f", realRate)
}
