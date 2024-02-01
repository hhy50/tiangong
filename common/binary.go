package common

import "encoding/binary"

// 基于大端存储法, 高位补0

func Uint16(b []byte) uint16 {
	if len(b) == 0 {
		return 0
	}
	if len(b) == 1 {
		return uint16(b[0])
	}
	return binary.BigEndian.Uint16(b)
}

func Uint32(b []byte) uint32 {
	if len(b) == 0 {
		return 0
	}
	if len(b) == 1 {
		return uint32(b[0])
	}
	if len(b) == 2 {
		return uint32(binary.BigEndian.Uint16(b))
	}
	if len(b) == 3 {
		return uint32(binary.BigEndian.Uint16(b[1:])) + uint32(b[0])<<16
	}
	return binary.BigEndian.Uint32(b)
}

func Uint64(b []byte) uint64 {
	length := len(b)
	switch length {
	case 0, 1, 2, 3, 4:
		return uint64(Uint32(b))
	default:
		return uint64(Uint32(b[length-4:length])) + uint64(Uint32(b[0:length-4]))<<32
	}
}

func Uint64ToBytes(i uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, i)
	return bytes
}

func Uint32ToBytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}

func Uint16ToBytes(i uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, i)
	return bytes
}
