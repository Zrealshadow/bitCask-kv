package store

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileManager struct {
	meta *Meta
	// activeFile *appendFile
	// dcnt       map[int64]int         // record fid -> duplicated item
	// index      map[string]*Item      // record key -> item {position in file}
	// afmap      map[int64]*appendFile // record fid -> open file handler

	dcnt  *sync.Map // record fid -> duplicated item
	index *sync.Map // record key -> item {position in file}
	afmap *sync.Map // record fid -> open file handler

	mu sync.Mutex
	// this mux is to protect active file is not change when concurrent write
	// File Write is concurrent safe, but we have to protect the file target is same
	done chan struct{}
}

func NewFileManager(m *Meta) (*FileManager, error) {
	mf := &FileManager{meta: m, dcnt: &sync.Map{}, index: &sync.Map{}, afmap: &sync.Map{}, done: make(chan struct{})}
	err := mf.Load()
	if err != nil {
		return nil, err
	}
	go mf.MergeMoniter()
	return mf, nil
}

func (fm *FileManager) Load() error {
	t := time.Now()
	log.Printf("[%v] Start to Load Meta data of FileManager\n", t)
	err := fm.loadMeta()
	if err != nil {
		return err
	}
	d := time.Since(t)
	log.Printf("[%v] Finish to load Meta data of FileManager, Spend %v \n", time.Now(), d)

	t = time.Now()
	if fm.meta.IndexValid {
		log.Printf("[%v] Start to Load Index data of FileManager", t)
		err := fm.loadIndex()
		if err != nil {
			return err
		}
		d := time.Since(t)
		log.Printf("[%v] Finish to Load Index data of FileManager, Spend %v", time.Now(), d)
	} else {
		log.Printf("[%v] Start to Scan all data to generate index map\n", t)
		/*
			pass all file according to time order and load in index map
		*/
		fids := fm.meta.GetFids()
		for _, fid := range fids {
			af, ok := fm.afmap.Load(fid)
			if ok {
				err := fm.loadAppendFile(af.(*appendFile))
				if err != nil {
					log.Printf("Failed to Load index in fid : %d \n: err: %v", fm.meta.ActiveFid, err)
				}
			}
		}
		af, _ := fm.GetAF()
		err := fm.loadAppendFile(af)
		if err != nil {
			log.Printf("Failed to Load index in fid : %d \n: err: %v", fm.meta.ActiveFid, err)
		}
		d := time.Since(t)
		log.Printf("[%v] Finish to Load Index data of FileManager, Spend %v\n", time.Now(), d)
	}
	return nil
}

func (fm *FileManager) loadMeta() error {

	if fm.meta == nil {
		return fmt.Errorf("meta info is nil when create file manager")
	}
	// load old files
	fids := fm.meta.GetFids()
	datadir := filepath.Join(fm.meta.Dirpath, defaultFileDirName)
	for _, fid := range fids {
		handler, err := NewAppendFile(datadir, OLD, fid)
		if err != nil {
			return err
		}
		fm.afmap.Store(handler.fid, handler)
	}

	// load active files

	handler, err := NewAppendFile(datadir, ACTIVE, fm.meta.ActiveFid)
	if err != nil {
		return err
	}

	fm.afmap.Store(handler.fid, handler)
	return nil
}

func (fm *FileManager) loadAllFile() error {
	fids := fm.meta.GetFids()
	// sorted fids according to time
	for _, fid := range fids {
		af, ok := fm.afmap.Load(fid)
		if !ok {
			return fmt.Errorf("fid %d is not Found in afmap", fid)
		}
		err := fm.loadAppendFile(af.(*appendFile))
		if err != nil {
			return err
		}
	}
	return nil
}

func (fm *FileManager) loadAppendFile(af *appendFile) error {
	b := make([]byte, 9)
	off := int64(0)
	var err error
	n := 0
	for {
		n, err = af.Read(off, b)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if n < 9 {
			return nil
		}

		kv, err := decodeRecordHeader(b)
		if err != nil {
			return err
		}
		s := int(kv.ValueSize) + int(kv.KeySize) + 9
		d := make([]byte, s)
		n, err = af.Read(off, d)

		if n < s {
			return nil
		}

		kv, err = decodeRecord(d)

		fm.index.Store(string(kv.Key), NewItem(af.fid, off, uint32(s)))
		// fm.index[string(kv.Key)] = NewItem(af.fid, off, uint32(s))

		off += int64(s)
	}
}

func (fm *FileManager) loadIndex() error {
	fn := filepath.Join(fm.meta.Dirpath, defaultIndexName)
	f, err := os.OpenFile(fn, os.O_RDONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()
	b := make([]byte, 2)
	off := int64(0)
	n := 0
	for {
		n, err = f.ReadAt(b, off)

		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if n < 2 {
			break
		}

		s := binary.BigEndian.Uint16(b)
		d := make([]byte, s)
		n, err = f.ReadAt(d, off)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		it := &Item{}
		key := decodeIndex(d, it)
		fm.index.Store(string(key), it)
		// fm.index[string(key)] = it
		off += int64(s)
	}
	return nil
}

func (fm *FileManager) saveIndex() error {
	fn := filepath.Join(fm.meta.Dirpath, defaultIndexName)
	f, err := os.OpenFile(fn, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// for key, value := range fm.index {
	// 	b := encodeIndex(key, value)
	// 	_, err := f.Write(b)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	fm.index.Range(
		func(key, value interface{}) bool {
			k := key.(string)
			it := value.(*Item)
			b := encodeIndex(k, it)
			_, err := f.Write(b)
			if err != nil {
				// return false
				log.Printf("index save key = %s, item = %v, err = %v\n", k, it, err)
			}
			return true
		})
	return nil
}

func (fm *FileManager) GetOlderFile() ([]*appendFile, error) {
	fhs := make([]*appendFile, 0)
	for _, fid := range fm.meta.GetFids() {
		af, ok := fm.afmap.Load(fid)
		if !ok {
			return nil, fmt.Errorf("fid %d is not exist in afmap", fid)
		}
		fhs = append(fhs, af.(*appendFile))
	}
	return fhs, nil
}

func (fm *FileManager) GetAF() (*appendFile, error) {
	fid := fm.meta.ActiveFid
	af, ok := fm.afmap.Load(fid)
	if !ok {
		return nil, fmt.Errorf("fid %d is not exist in afmap", fid)
	}
	return af.(*appendFile), nil
}

func (fm *FileManager) Write(k []byte, v []byte) error {
	// t := time.Now().Unix()
	kv, err := NewRecord(k, v)
	if err != nil {
		return err
	}

	b := encodeRecord(kv)

	fm.mu.Lock()
	defer fm.mu.Unlock()

	af, err := fm.GetAF()
	if err != nil {
		return err
	}

	if af.IsClosed() {
		return fmt.Errorf("active fi is timeout")
	}
	// before write, we should judge is it necessary to create a new file
	activefilesize, err := af.Size()
	if err != nil {
		return err
	}
	if activefilesize+int64(len(b)) > MaxAppendFileSize {
		// create new active file
		af.SetOlder()
		d := filepath.Join(fm.meta.Dirpath, defaultFileDirName)
		new_af, err := NewAppendFile(d, ACTIVE, time.Now().UnixNano())
		if err != nil {
			return err
		}

		// update map
		fm.afmap.Store(new_af.fid, new_af)
		// update meta
		fm.meta.OlderFids = append(fm.meta.OlderFids, af.fid)
		fm.meta.ActiveFid = new_af.fid
		fm.meta.Save()
		af = new_af
	}

	off, err := af.Write(b)
	if err != nil {
		return fmt.Errorf("write(%s, %s) err:%v", k, v, err)
	}

	if itptr, ok := fm.index.Load(string(kv.Key)); ok {
		item := itptr.(*Item)
		if c, ok := fm.dcnt.Load(item.fid); !ok {
			fm.dcnt.Store(item.fid, uint32(0))
		} else {
			fm.dcnt.Store(item.fid, c.(uint32)+1)
		}
	}
	// fm.index[string(k)] = &Item{
	// 	fid:    af.fid,
	// 	offset: off,
	// 	size:   uint32(len(b)),
	// }
	fm.index.Store(string(k), &Item{
		fid:    af.fid,
		offset: off,
		size:   uint32(len(b)),
	})

	return nil
}

func (fm *FileManager) Read(k []byte) ([]byte, error) {
	itptr, ok := fm.index.Load(string(k))
	if !ok {
		return nil, fmt.Errorf("Key %s is not found", string(k))
	}
	item := itptr.(*Item)

	afptr, ok := fm.afmap.Load(item.fid)
	if !ok {
		return nil, fmt.Errorf("fid %d file is not found in afmap", item.fid)
	}

	af := afptr.(*appendFile)

	if af.IsClosed() {
		return nil, fmt.Errorf("afmap fid is timeout")
	}

	b := make([]byte, item.size)
	n, err := af.Read(item.offset, b)
	if err != nil {
		return nil, err
	}

	if n != int(item.size) {
		return nil, fmt.Errorf("read (%v), want %d bytes but get %d bytes", *item, item.size, n)
	}

	kv, err := decodeRecord(b)
	if err != nil {
		return nil, err
	}

	if string(kv.Value) == DELETEFLAG {
		fm.index.Delete(string(kv.Key))
		return nil, fmt.Errorf("Key %s is not found", string(k))
	}

	return kv.Value, nil
}

func (fm *FileManager) Merge(af *appendFile) error {

	b := make([]byte, 9)

	off := int64(0)
	var err error

	n := 0
	for {
		n, err = af.Read(off, b)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if n < 9 {
			break
		}

		kv, err := decodeRecordHeader(b)
		if err != nil {
			return err
		}

		s := int(kv.KeySize) + int(kv.ValueSize) + 9
		d := make([]byte, s)
		n, err = af.Read(off, d)
		if err != nil {
			return err
		}

		if n < s {
			return nil
		}
		kv, err = decodeRecord(d)
		if err != nil {
			return err
		}

		itptr, ok := fm.index.Load(string(kv.Key))
		it := itptr.(*Item)
		if ok {
			if af.fid == it.fid && off == it.offset && it.size == uint32(s) {
				// this value is the updated value
				if string(kv.Value) != DELETEFLAG {
					// write into active file
					err = fm.Write(kv.Key, kv.Value)
					if err != nil {
						return err
					}
				} else {
					// it is been flaged deleted
					fm.index.Delete(string(kv.Key))
				}
			}
		}
		off = off + int64(s)
	}
	return nil
}

func (fm *FileManager) MergeMoniter() {
	// every 1 minute , we will check weather is there any file have to be merged
	var d time.Duration

	if !Debug {
		d = time.Minute
	} else {
		d = time.Second * 5
	}
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {

		case <-fm.done:
			break
		case <-ticker.C:
			fm.dcnt.Range(func(fid, value interface{}) bool {
				cnt := value.(uint32)
				if cnt <= MaxDuplicatedIten {
					return true
				}

				// find appendfile
				af, ok := fm.afmap.Load(fid)
				if !ok {
					return true
				}

				if af.(*appendFile).GetRole() == ACTIVE {
					// no need to merge active file
					return true
				}

				err := fm.Merge(af.(*appendFile))
				if err != nil {
					log.Printf("fid %d merge failure , err = %v\n", fid, err)
				} else {
					// remove this file and update meta file
					fm.Remove(af.(*appendFile))
					log.Printf("fid = %d merge success \n", fid)
				}
				return true
			})
		}
	}
}

func (fm *FileManager) IndexMoniter() {
	lastActiveFid := fm.meta.ActiveFid
	for range time.Tick(10 * time.Second) {
		if fm.meta.ActiveFid == lastActiveFid {
			continue
		} else {
			// lastActive file  != 0 && ActiveFile update
			// save index
			t := time.Now()
			fm.saveIndex()
			d := time.Since(t)
			log.Printf("indexSave finished | spend time %v", d)
		}
	}
}

func (fm *FileManager) Remove(af *appendFile) {
	// close file
	af.Close()
	// update afmap dnt
	fm.afmap.Delete(af.fid)
	fm.dcnt.Delete(af.fid)
	// update MetaData
	fm.meta.DeleteFids(af.fid)
	fm.meta.Save()
	// rm this file in os filesystem
	os.Remove(af.fp)
}

func (fm *FileManager) Close() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.done <- struct{}{}
	close(fm.done)

	// close all appendfile handler
	fm.afmap.Range(func(key, value interface{}) bool {
		af := value.(*appendFile)
		af.Close()
		return true
	})
}

//--------------- Some helper Function for Unit Test ------------/

func (fm *FileManager) GetActiveFid() int64 {
	return fm.meta.ActiveFid
}
