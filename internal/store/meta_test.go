package store

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var datadir string = "../tmp"

func TestBasicMeta(t *testing.T) {
	setupMetaTest()
	defer teardownMetaTest()

	m, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Empty Meta error in dir %s, err: %v", datadir, err)
	}
	tt := time.Now().UnixNano()
	m.SetActiveFid(tt)
	for i := 0; i < 10; i++ {
		m.AddOldFid(int64(i))
	}
	m.Save()

	mm, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Empty Meta error in dir %s, err: %v", datadir, err)
	}

	if mm.ActiveFid != tt {
		t.Fatalf("Save Error in ActiveFid")
	}

	for idx, i := range mm.GetFids() {
		if int64(idx) != i {
			t.Fatalf("Save Error in OlderFid")
		}
	}
}

func TestEmptyMeta(t *testing.T) {
	setupMetaTest()
	defer teardownMetaTest()
	m, err := NewMeta(datadir)
	if err != nil {
		t.Fatalf("Create Empty Meta error in dir %s, err: %v", datadir, err)
	}
	// exist
	fp := filepath.Join(datadir, defaultMetaName)
	_, err = os.Stat(fp)
	if os.IsNotExist(err) {
		t.Fatalf("Create Meta but no meta files are created %s", m.Dirpath)
	}

	lastActiveFid := m.ActiveFid
	// create another Meta
	m, err = NewMeta(datadir)
	fmt.Printf("%+v\n", m)
	if lastActiveFid != m.ActiveFid {
		t.Fatalf("can not load exist meta file")
	}

}

func setupMetaTest() {
	os.MkdirAll(datadir, 0755)
}

func teardownMetaTest() {
	// os.RemoveAll(datadir)
}
