package udt

import (
	"bytes"
	"log"
	"testing"
)

func TestBitmap(t *testing.T) {
	bit := NewBitmap([16]byte{16})
	log.Println(bit)
}

func TestBitmapScanAndValue8(t *testing.T) {
	cases := []struct {
		name   string
		input  []byte
		expect [8]byte
	}{
		{"all zero", make([]byte, 8), [8]byte{}},
		{"all one", bytes.Repeat([]byte{0xFF}, 8), [8]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"mixed", []byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}, [8]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}},
	}
	for _, c := range cases {
		bm := NewBitmap([8]byte{})
		if err := bm.Scan(c.input); err != nil {
			t.Errorf("Scan failed for %s: %v", c.name, err)
			continue
		}
		if bm.bit != c.expect {
			t.Errorf("Scan result mismatch for %s: got %v, want %v", c.name, bm.bit, c.expect)
		}
		val, err := bm.Value()
		if err != nil {
			t.Errorf("Value failed for %s: %v", c.name, err)
			continue
		}
		bs, ok := val.([]byte)
		if !ok {
			t.Errorf("Value type assertion failed for %s", c.name)
			continue
		}
		if !bytes.Equal(bs[:], c.input) {
			t.Errorf("Value result mismatch for %s: got %v, want %v", c.name, bs, c.input)
		}
	}
}

func BenchmarkBitmapScanValue(b *testing.B) {
	type bmCase struct {
		name  string
		size  int
		newbm func() any
		scan  func(any, []byte) error
		value func(any) ([]byte, error)
	}

	table := []bmCase{
		{
			name:  "8byte",
			size:  8,
			newbm: func() any { return NewBitmap([8]byte{}) },
			scan:  func(bm any, data []byte) error { return bm.(*Bitmap[[8]byte]).Scan(data) },
			value: func(bm any) ([]byte, error) {
				v, err := bm.(*Bitmap[[8]byte]).Value()
				if err != nil {
					return nil, err
				}
				return v.([]byte), nil
			},
		},
		{
			name:  "255byte",
			size:  255,
			newbm: func() any { return NewBitmap([255]byte{}) },
			scan:  func(bm any, data []byte) error { return bm.(*Bitmap[[255]byte]).Scan(data) },
			value: func(bm any) ([]byte, error) {
				v, err := bm.(*Bitmap[[255]byte]).Value()
				if err != nil {
					return nil, err
				}
				return v.([]byte), nil
			},
		},
	}

	for _, tc := range table {
		b.Run(tc.name+"/Scan", func(b *testing.B) {
			data := bytes.Repeat([]byte{0xAA}, tc.size)
			b.ResetTimer()
			for b.Loop() {
				bm := tc.newbm()
				_ = tc.scan(bm, data)
			}
		})
		b.Run(tc.name+"/Value", func(b *testing.B) {
			bm := tc.newbm()
			data := bytes.Repeat([]byte{0xAA}, tc.size)
			_ = tc.scan(bm, data)
			b.ResetTimer()
			for b.Loop() {
				_, _ = tc.value(bm)
			}
		})
	}
}
