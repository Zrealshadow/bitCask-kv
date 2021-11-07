package bitcask

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"math"
)

type Record struct {
	Key       []byte
	Value     []byte
	KeySize   uint8
	ValueSize uint32
	Crc       uint32
}

func NewRecord(key []byte, value []byte) (*Record, error) {
	kl := len(key)
	vl := len(value)

	if kl > math.MaxUint8 || vl > math.MaxUint32 {
		return nil, errors.New("key or value length exceeds maximum")
	}

	r := &Record{
		KeySize:   uint8(kl),
		ValueSize: uint32(vl),
		Key:       key,
		Value:     value,
	}
	r.Crc = 0
	r.Crc = crc32.ChecksumIEEE(encodeRecord(r))
	return r, nil
}

func encodeRecord(r *Record) []byte {
	b := make([]byte, uint(r.KeySize)+uint(r.ValueSize)+9)
	binary.BigEndian.PutUint32(b[0:4], r.Crc)       // 32 / 8 = 4
	b[4] = r.KeySize                                // 8/8 = 1
	binary.BigEndian.PutUint32(b[5:9], r.ValueSize) // 32 / 8 =4
	s, e := uint(9), 9+uint(r.KeySize)
	copy(b[s:e], r.Key)
	s, e = 9+uint(r.KeySize), 9+uint(r.KeySize)+uint(r.ValueSize) // uint default uint32
	copy(b[s:e], r.Value)
	return b
}

func decodeRecord(b []byte) (*Record, error) {
	if len(b) < 9 {
		return nil, errors.New("packet header exception")
	}
	r := &Record{}
	// decode crc
	crc := binary.BigEndian.Uint32(b[:4])
	r.KeySize = uint8(b[4])
	r.ValueSize = binary.BigEndian.Uint32(b[5:9])
	// fmt.Printf("KeySize %d, ValueSize %d\n", r.KeySize, r.ValueSize)
	s, e := uint(9), 9+uint(r.KeySize)
	r.Key = b[s:e]
	s, e = 9+uint(r.KeySize), 9+uint(r.KeySize)+uint(r.ValueSize)
	r.Value = b[s:e]
	r.Crc = crc32.ChecksumIEEE(encodeRecord(r))
	if crc != r.Crc {
		return nil, errors.New("packet crc32 exception")
	}
	return r, nil
}

func decodeRecordHeader(b []byte) (*Record, error) {
	if len(b) < 9 {
		return nil, errors.New("packet header exception")
	}

	r := &Record{
		Crc:       binary.BigEndian.Uint32(b[:4]),
		KeySize:   b[4],
		ValueSize: binary.BigEndian.Uint32(b[5:9]),
	}
	return r, nil
}
