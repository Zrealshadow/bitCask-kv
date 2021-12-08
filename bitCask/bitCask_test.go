package bitcask

import (
	"fmt"
	"os"
	"testing"
)

const MountDir string = "../storage/DB"

func TestBitCaskBasic(t *testing.T) {
	//cretae BitCaskEngine
	defer TearDownTmp()
	e, err := NewBitCaskEngine(MountDir)
	if err != nil {
		t.Fatalf("Create BitCask Engine Error")
	}
	err = e.NewBlock("Test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	// b, _ := e.Block("Test")
	var klist []string
	for i := 0; i < 10000; i++ {
		k := fmt.Sprintf("K:%d", i*100)
		klist = append(klist, k)
		e.Put("Test", []byte(klist[i]), []byte(klist[i]))
		// b.Put([]byte(klist[i]), []byte(klist[i]))
	}

	//check
	for i := 0; i < 10000; i++ {
		v, _ := e.Get("Test", []byte(klist[i]))
		// v, _ := b.Get([]byte(klist[i]))
		value := string(v)
		if value != klist[i] {
			t.Fatalf("Get Error")
		}
	}

}

func TearDownTmp() {
	os.RemoveAll(MountDir)
}
