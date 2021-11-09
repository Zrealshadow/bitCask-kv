package bitcask

import "encoding/binary"

type Item struct {
	fid    int64
	offset int64
	size   uint32
}

func NewItem(fid int64, offset int64, size uint32) *Item {
	return &Item{fid: fid, offset: offset, size: size}
}

/*
	first element save bytes length of all Index in 2 byte -> int16
	item size : fid(8) + offset(8) + size(4) == 20
	key size : kl
	thus all length is 20 + 2 + kl saved in first 2 byte
*/
func encodeIndex(key string, it *Item) []byte {
	kl := len(key)
	s := kl + 22
	b := make([]byte, s)
	binary.BigEndian.PutUint16(b[:2], uint16(s))
	binary.BigEndian.PutUint64(b[2:10], uint64(it.fid))
	binary.BigEndian.PutUint64(b[10:18], uint64(it.offset))
	binary.BigEndian.PutUint32(b[18:22], it.size)
	copy(b[22:s], []byte(key))
	return b
}

func decodeIndex(b []byte, it *Item) []byte {
	s := binary.BigEndian.Uint16(b[:2])
	kl := s - 22
	it.fid = int64(binary.BigEndian.Uint64(b[2:10]))
	it.offset = int64(binary.BigEndian.Uint64(b[10:18]))
	it.size = binary.BigEndian.Uint32(b[18:22])
	k := make([]byte, kl)
	copy(k, b[22:s])
	return k
}
