package store

import (
	"fmt"
	"os"
	"testing"
)

var kvs = []*struct {
	record string
	offset int64
}{
	{"1", 0},
	{"2", 0},
	{"3", 0},
	{"4", 0},
	{"5", 0},
	{"6", 0},
	{"Skip-list", 0},
	{"B tree", 0},
	{"B+ tree", 0},
	{"NUS Garbage Collection University", 0},
	{"Fuck NUS HR", 0},
}

var dirpath string = "./"

func TestBasicAppendFile(t *testing.T) {
	fmt.Printf("\nBasic Test of Appendfile:\n\n")
	af, err := NewAppendFile(dirpath, ACTIVE, int64(1))
	_, err = os.Stat(af.fp)

	if os.IsNotExist(err) {
		t.Fatalf("appendfile %s created failed in NewAppendiFile Function", af.fp)
	}

	for _, kv := range kvs {
		offset, err := af.Write([]byte(kv.record))
		if err != nil {
			t.Fatal(err.Error())
		}
		kv.offset = offset
	}

	for idx, kv := range kvs {
		b := make([]byte, len(kv.record))
		n, err := af.Read(kv.offset, b)

		if err != nil {
			t.Fatal(err.Error())
		}

		if n != len(kv.record) {
			t.Fatalf("[id:%d record %s] Error Read bytes Want Read %d bytes , but got %d bytes", idx, kv.record, len(kv.record), n)
		}

		r := string(b)

		if r != kv.record {
			t.Fatalf("[id:%d example] Error Read want %s got %s", idx, kv.record, r)
		}
	}

	defer func() {
		// delete created file
		af.Close()
		os.Remove(af.fp)
	}()
}
