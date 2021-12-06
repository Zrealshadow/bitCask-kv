package store

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
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

// Test Merge :
// Register all filenames, and check exist file
// Keep anyfile exist all 50 duplicated k-v

func TestMergeFileManager(t *testing.T) {
	setupMetaTest()
	defer teardownMetaTest()
	meta, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Meta file Error : Err: %v", err)
	}

	fm, err := NewFileManager(meta)

	if err != nil {
		t.Fatalf("Create FileManager Error: Err: %v", err)
	}

	var klist []int
	for i := 1; i < 5000; i++ {
		klist = append(klist, i*100)
	}
	aph, err := fm.GetAF()
	activeFid := aph.fid
	existFileNameList := make([]string, 0)
	existFileNameList = append(existFileNameList, aph.fp)
	for _, key := range klist {
		for ts := 0; ts < 50; ts++ {
			k := fmt.Sprintf("%d", key)
			v := fmt.Sprintf("%d", key+ts)
			err := fm.Write([]byte(k), []byte(v))
			if err != nil {
				t.Fatalf("Write key %s Value %s Failed", k, v)
			}

			if fm.GetActiveFid() != activeFid {
				// update activefid
				aph, _ := fm.GetAF()
				activeFid = aph.fid
				existFileNameList = append(existFileNameList, aph.fp)
			}
		}
	}
	time.Sleep(1 * time.Minute)

	files, err := ioutil.ReadDir(datadir)

	if len(files) == len(existFileNameList) {
		t.Fatalf("No Merge Operations are triggered")
	}

	// check k - v
	for _, key := range klist {
		k := fmt.Sprintf("%d", key)
		v := fmt.Sprintf("%d", key+49)

		v_, _ := fm.Read([]byte(k))
		gotv_ := string(v_)

		if v != gotv_ {
			t.Fatalf("key %s Want V %s Got V %s\n", k, v, gotv_)
		}
	}
	fm.Close()

}

func TestConcurrentReadAndWriteFileManager(t *testing.T) {
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

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(workId int, f *FileManager) {
			defer wg.Done()
			klist := make([]string, 50)
			for i := 0; i < 50; i++ {
				klist[i] = fmt.Sprintf("%d", i) + "_work:" + fmt.Sprintf("%d", workId)
			}

			for _, key := range klist {
				f.Write([]byte(key), []byte(key))
			}

			// check read
			for _, key := range klist {
				v, _ := f.Read([]byte(key))
				value := string(v)
				if value != key {
					panic(fmt.Sprintf("Key %s Want %s Got %s", key, key, value))
				}
			}
		}(i, fm)
	}
	wg.Wait()
	fm.Close()
}

func TestReconnectToFileManager(t *testing.T) {
	SetupFileManager()
	defer teardownMetaTest()

	meta, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Meta file Error : Err: %v", err)
	}

	fm, err := NewFileManager(meta)
	if err != nil {
		t.Fatalf("Create FileManager Error: Err: %v", err)
	}

	// we insert many key-value
	klist := make([]string, 0)
	for i := 0; i < 100000; i++ {
		klist = append(klist, fmt.Sprint(i))
		fm.Write([]byte(klist[i]), []byte(klist[i]))
	}
	// Close the FM
	fm.Close()

	// now we re-declare it
	meta_, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Retry Create Meta file Error : Err: %v", err)
	}

	fm_, err := NewFileManager(meta_)
	if err != nil {
		t.Fatalf("Retry Create FileManager Error: Err: %v", err)
	}

	// check last inserted record
	for _, key := range klist {
		v, _ := fm_.Read([]byte(key))
		value := string(v)
		if value != key {
			t.Fatalf("ReStart FM but failed to get correct k-v Key %s Want %s Got %s", key, key, value)
		}
	}
}

func SetupFileManager() {
	os.MkdirAll(datadir, 0755)
}

func TearDownFileManager() {
	os.RemoveAll(datadir)
}
