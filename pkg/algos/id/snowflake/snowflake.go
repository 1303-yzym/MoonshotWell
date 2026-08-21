package snowflake

import (
	"sync/atomic"
	"time"
)

const (
	signBits       = 1  // head
	timestampBits  = 41 // 毫秒
	machineIDBits  = 10 // 机器位数 2^10 1024
	sequenceBits   = 12 // 序列位数 毫秒 (2^12)*1000 = 4,096,000
	maxMachineID   = -1 ^ (-1 << machineIDBits)
	maxSequenceNum = -1 ^ (-1 << sequenceBits)
)

type Snowflake struct {
	timestamp   atomic.Int64
	machineID   int64
	sequenceNum atomic.Int64
}

func New(machineID int64) *Snowflake {
	if machineID < 0 || machineID > maxMachineID {
		panic("Invalid machine ID")
	}

	return &Snowflake{
		machineID: machineID,
	}
}

func (s *Snowflake) GenerateID() int64 {
	currentTimestamp := time.Now().UnixNano() / 1e6

	lastTimestamp := s.timestamp.Load()
	if currentTimestamp < lastTimestamp {
		panic("Clock moved backwards, refusing to generate ID")
	}

	sequence := s.sequenceNum.Load()

	if currentTimestamp == lastTimestamp {
		sequence = (sequence + 1) & maxSequenceNum
		if sequence == 0 {
			currentTimestamp = s.waitNextMillis(lastTimestamp)
		}
	} else {
		sequence = 0
	}

	s.timestamp.Store(currentTimestamp)
	s.sequenceNum.Store(sequence)

	id := (currentTimestamp << (machineIDBits + sequenceBits)) |
		(s.machineID << sequenceBits) |
		sequence

	return id
}

func (s *Snowflake) waitNextMillis(lastTimestamp int64) int64 {
	currentTimestamp := time.Now().UnixNano() / 1e6
	for currentTimestamp <= lastTimestamp {
		currentTimestamp = time.Now().UnixNano() / 1e6
	}

	return currentTimestamp
}
