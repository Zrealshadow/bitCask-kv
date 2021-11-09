package store

import (
	"fmt"
	"testing"
)

func TestBasicRecord(t *testing.T) {
	fmt.Printf("\nBasic Test of KV Record:\n\n")
	var kvs = []struct {
		key   string
		value string
	}{
		{"1", "Go"},
		{"2", "Python"},
		{"3", "C++"},
		{"4", "Java"},
		{"5", "C"},
		{"6", "Javascript"},
	}

	for _, kv := range kvs {
		r, err := NewRecord([]byte(kv.key), []byte(kv.value))
		if err != nil {
			t.Error(err.Error())
		}

		if r.KeySize != uint8(len(kv.key)) || r.ValueSize != uint32(len(kv.value)) {
			t.Errorf("Size is not same in Record k,v: %d %d\t origin %d %d", r.KeySize, r.ValueSize, len(kv.key), len(kv.value))
		}

		b := encodeRecord(r)

		dr, _ := decodeRecord(b)
		drh, _ := decodeRecordHeader(b)

		if dr.Crc != r.Crc {
			t.Errorf("Checksum ruined in encoding - decoding process, origin : %d, Got :%d", r.Crc, dr.Crc)
		}
		CheckStr(t, string(dr.Key), kv.key, "Key")
		CheckStr(t, string(dr.Value), kv.value, "Value")
		CheckInt32(t, uint32(dr.KeySize), uint32(r.KeySize), "KeySize")
		CheckInt32(t, uint32(dr.ValueSize), r.ValueSize, "ValueSize")

		if drh.Crc != r.Crc {
			t.Errorf("Checksum ruined in encoding - decoding head process, origin : %d, Got :%d", r.Crc, drh.Crc)
		}

		CheckInt32(t, uint32(drh.KeySize), uint32(r.KeySize), "RecordHead KeySize")
		CheckInt32(t, uint32(drh.ValueSize), r.ValueSize, "RecordHeadValueSize")
	}
}

func CheckStr(t *testing.T, get string, want string, desc string) {
	if get != want {
		t.Errorf("Get %s Want %s in Field %s Compare", get, want, desc)
	}
}

func CheckInt32(t *testing.T, get uint32, want uint32, desc string) {
	if get != want {
		t.Errorf("Get %d Want %d in Field %s Compare", get, want, desc)
	}
}
