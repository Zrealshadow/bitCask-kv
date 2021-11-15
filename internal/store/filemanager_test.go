package store

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestBasicReadAndWriteFileManager(t *testing.T) {
	SetupFileManager()
	defer TearDownFileManager()

	meta, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Meta file Error : Err: %v", err)
	}

	fm, err := NewFileManager(meta)
	if err != nil {
		t.Fatalf("Create FileManager Error: Err: %v", err)
	}
	m := make(map[string]string)
	var klist []int
	var spaceValue []int
	for i := 1; i < 21; i++ {
		klist = append(klist, i*100)
		spaceValue = append(spaceValue, i*5)
	}

	for ts := 0; ts < 100000; ts++ {
		v := klist[rand.Intn(20)]
		space := spaceValue[rand.Intn(20)]
		key := fmt.Sprint(v)
		value := fmt.Sprintf("%d_%d", v, v+space)
		err := fm.Write([]byte(key), []byte(value))
		if err != nil {
			t.Fatalf("Write Fatal Key: %s Value: %s", key, value)
		}
		m[key] = value
	}

	// check read
	for k, want := range m {
		g, err := fm.Read([]byte(k))
		if err != nil {
			t.Fatalf("Read Key: %s Fatal", k)
		}
		get := string(g)
		if want != string(get) {
			t.Fatalf("KV Key: %s Got :%s Want :%s", k, get, want)
		}
	}
}

func TestConcurrentReadAndWriteFileManager(t *testing.T) {
	SetupFileManager()
	defer TearDownFileManager()

	meta, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Meta file Error : Err: %v", err)
	}

	_, err := NewFileManager(meta)
	if err != nil {
		t.Fatalf("Create FileManager Error: Err: %v", err)
	}

}
func SetupFileManager() {
	os.MkdirAll(datadir, 0755)
}

func TearDownFileManager() {
	// os.RemoveAll(datadir)
}
