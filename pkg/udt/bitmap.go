package udt

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type FixedBinary interface {
	~[8]byte | ~[16]byte | ~[32]byte | ~[64]byte | ~[128]byte | ~[255]byte
}

type Bitmap[B FixedBinary] struct {
	bit B
}

func NewBitmap[B FixedBinary](b B) *Bitmap[B] {
	return &Bitmap[B]{bit: b}
}

func (b *Bitmap[B]) Set(i uint, value bool) {
	byteIdx := i / 8
	bitIdx := i % 8

	switch a := any(&b.bit).(type) {
	case *[8]byte:
		if byteIdx < 8 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	case *[16]byte:
		if byteIdx < 16 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	case *[32]byte:
		if byteIdx < 32 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	case *[64]byte:
		if byteIdx < 64 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	case *[128]byte:
		if byteIdx < 128 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	case *[255]byte:
		if byteIdx < 255 {
			if value {
				a[byteIdx] |= 1 << bitIdx
			} else {
				a[byteIdx] &^= 1 << bitIdx
			}
		}
	}
}

func (b *Bitmap[B]) Get(i uint) bool {
	byteIdx := i / 8
	bitIdx := i % 8

	switch a := any(b.bit).(type) {
	case [8]byte:
		if byteIdx < 8 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	case [16]byte:
		if byteIdx < 16 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	case [32]byte:
		if byteIdx < 32 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	case [64]byte:
		if byteIdx < 64 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	case [128]byte:
		if byteIdx < 128 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	case [255]byte:
		if byteIdx < 255 {
			return (a[byteIdx] & (1 << bitIdx)) != 0
		}
	}

	return false
}

func (b *Bitmap[B]) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	case nil:
		var zero B

		b.bit = zero

		return nil
	default:
		return errors.New("type assertion .([]byte) or (string) failed")
	}

	copyToBitArray := func(arr []byte) error {
		if len(data) != len(arr) {
			return fmt.Errorf("length mismatch: got %d, want %d", len(data), len(arr))
		}

		copy(arr, data)

		return nil
	}

	switch arr := any(&b.bit).(type) {
	case *[8]byte:
		return copyToBitArray(arr[:])
	case *[16]byte:
		return copyToBitArray(arr[:])
	case *[32]byte:
		return copyToBitArray(arr[:])
	case *[64]byte:
		return copyToBitArray(arr[:])
	case *[128]byte:
		return copyToBitArray(arr[:])
	case *[255]byte:
		return copyToBitArray(arr[:])
	default:
		return errors.New("unsupported BitArray type")
	}
}

func (b Bitmap[B]) Value() (driver.Value, error) {
	switch arr := any(&b.bit).(type) {
	case *[8]byte:
		return arr[:], nil
	case *[16]byte:
		return arr[:], nil
	case *[32]byte:
		return arr[:], nil
	case *[64]byte:
		return arr[:], nil
	case *[128]byte:
		return arr[:], nil
	case *[255]byte:
		return arr[:], nil
	default:
		return nil, errors.New("unsupported BitArray type")
	}
}

func (b *Bitmap[B]) String() string {
	var n int

	switch any(b.bit).(type) {
	case [8]byte:
		n = 8 * 8
	case [16]byte:
		n = 16 * 8
	case [32]byte:
		n = 32 * 8
	case [64]byte:
		n = 64 * 8
	case [128]byte:
		n = 128 * 8
	case [255]byte:
		n = 255 * 8
	default:
		return ""
	}

	bits := make([]byte, n)
	for i := 0; i < n; i++ {
		if b.Get(uint(i)) {
			bits[i] = '1'
		} else {
			bits[i] = '0'
		}
	}

	return string(bits)
}
