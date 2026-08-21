package bloomFilter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/fnv"

	"github.com/twmb/murmur3"
)

var (
	Hash64    = fnv.New64()
	MD5       = NewMD5()
	SHA1      = NewSHA1()
	SHA256    = NewSHA256()
	SHA512    = NewSHA512()
	MurMur32  = murmur3.New32()
	MurMur64  = murmur3.New64()
	MurMur128 = murmur3.New128()
)

type MD5Hash64 struct {
	hash.Hash
}

func (m MD5Hash64) Sum64() uint64 {
	b := m.Sum(nil)

	return bytes16ToUint64(b)
}

func NewMD5() hash.Hash64 {
	m := &MD5Hash64{
		md5.New(),
	}

	return m
}

type SHA1Hash64 struct {
	hash.Hash
}

func (h SHA1Hash64) Sum64() uint64 {
	b := h.Sum(nil)

	return bytes16ToUint64(b)
}

func NewSHA1() hash.Hash64 {
	m := &MD5Hash64{
		sha1.New(),
	}

	return m
}

type SHA256Hash64 struct {
	hash.Hash
}

func NewSHA256() hash.Hash64 {
	m := &MD5Hash64{
		sha256.New(),
	}

	return m
}

func (h SHA256Hash64) Sum64() uint64 {
	b := h.Sum(nil)

	return bytes16ToUint64(b)
}

type SHA512Hash64 struct {
	hash.Hash
}

func NewSHA512() hash.Hash64 {
	m := &MD5Hash64{
		sha512.New(),
	}

	return m
}

func (h SHA512Hash64) Sum64() uint64 {
	b := h.Sum(nil)

	return bytes16ToUint64(b)
}

func bytes16ToUint64(b []byte) uint64 {
	return uint64(b[7]^b[15]) | uint64(b[6]^b[14])<<8 | uint64(b[5]^b[13])<<16 | uint64(b[4]^b[12])<<24 |
		uint64(b[3]^b[11])<<32 | uint64(b[2]^b[10])<<40 | uint64(b[1]^b[9])<<48 | uint64(b[0]^b[8])<<56
}
