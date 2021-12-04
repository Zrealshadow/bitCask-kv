package store

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type FileRole int

const (
	ACTIVE FileRole = iota + 1
	OLD
)

type appendFile struct {
	fp     string
	offset int64
	role   FileRole
	fo     *os.File
	fid    int64
}

func NewAppendFile(dir string, role FileRole, fid int64) (*appendFile, error) {
	if role != OLD && role != ACTIVE {
		return nil, fmt.Errorf("File Role %d is not found", role)
	}

	af := &appendFile{
		fp:     filepath.Join(dir, strconv.FormatInt(fid, 10)),
		offset: 0,
		role:   role,
		fid:    fid,
	}

	var err error

	if role == ACTIVE {

		af.fo, err = os.OpenFile(af.fp, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		af.offset, err = af.Size()

		if err != nil {
			return nil, err
		}

		return af, nil
	}

	af.fo, err = os.OpenFile(af.fp, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return af, nil
}

func (af *appendFile) Size() (int64, error) {
	fi, err := af.fo.Stat()
	if err != nil {
		return -1, err
	}
	return fi.Size(), nil
}

// return offset of this record
func (af *appendFile) Write(b []byte) (int64, error) {

	if af.role == OLD {
		return -1, fmt.Errorf("write operation are not supported in old file %s", af.fp)
	}

	off := af.offset
	n, err := af.fo.Write(b)
	if err != nil {
		return -1, err
	}
	if n != len(b) {
		af.fo.Seek(off, 0)
		return -1, fmt.Errorf("write %d bytes, actually write %d bytes", len(b), n)
	}

	af.offset += int64(n)
	return off, nil
}

func (af *appendFile) Read(offset int64, b []byte) (int, error) {
	return af.fo.ReadAt(b, offset)
}

func (af *appendFile) SetOlder() {
	af.role = OLD
}

func (af *appendFile) Close() {
	af.Sync()
	if af.fo != nil {
		af.fo.Close()
	}
}

func (af *appendFile) IsClosed() bool {

	return af.fo == nil
}

func (af *appendFile) GetRole() FileRole {
	return af.role
}

func (af *appendFile) Sync() {

	if af.fo != nil {
		err := af.fo.Sync()
		if err != nil {
			log.Fatalf("sync %v", err)
		}
	}
}
