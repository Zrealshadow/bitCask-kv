package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Meta struct {
	Dirpath    string  `json:"dirpath"`
	OlderFids  []int64 `json:"olderFids"`
	ActiveFid  int64   `json:"activeFid"`
	IndexValid bool    `json:"indexValid"`
}

func NewMeta(dirpath string) (*Meta, error) {
	fp := filepath.Join(dirpath, defaultMetaName)
	ok, err := IsExist(fp)
	if err != nil {
		return nil, err
	}

	if ok {
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			return nil, err
		}
		m := &Meta{}
		err = json.Unmarshal(b, m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	m := &Meta{ActiveFid: time.Now().UnixNano(), OlderFids: make([]int64, 0), IndexValid: defaultIndexValid}
	err = m.Save()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Meta) Save() error {
	fp := filepath.Join(m.Dirpath, defaultMetaName)
	ok, err := IsExist(fp)
	if err != nil {
		return err
	}

	if !ok {
		err = os.MkdirAll(m.Dirpath, 0755)
		if err != nil {
			return err
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return ioutil.WriteFile(fp, b, 0644)
}

func (m *Meta) GetFids() []int64 {
	return m.OlderFids
}

func IsExist(fp string) (bool, error) {
	_, err := os.Stat(fp)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (m *Meta) DeleteFids(fid int64) {
	newOlderMid := make([]int64, len(m.OlderFids)-1)
	for _, id := range m.OlderFids {
		if id == fid {
			continue
		}
		newOlderMid = append(newOlderMid, id)
	}
	m.OlderFids = newOlderMid
}
