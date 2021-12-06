package bitcask

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	store "github.com/bitCaskKV/bitCask/internal/store"
)

type BitCaskEngine struct {
	blockMap map[string]*BitCaskBlock
	MountDir string
}

func NewBitCaskEngine(MountDir string) *BitCaskEngine {
	_, err := os.Stat(MountDir)
	engine := &BitCaskEngine{
		blockMap: map[string]*BitCaskBlock{},
		MountDir: MountDir,
	}
	if os.IsNotExist(err) {
		return engine
	}

	// Load exist Mount DB
	fmlist, err := ioutil.ReadDir(MountDir)
	for _, fm := range fmlist {
		path := MountDir + string(os.PathSeparator) + fm.Name()
		engine.blockMap[fm.Name()] = NewBitCaskBlock(path)
	}
	return engine
}

func (e *BitCaskEngine) NewBlock(blockName string) {
	if _, ok := e.blockMap[blockName]; !ok {
		path := e.MountDir + string(os.PathSeparator) + blockName
		e.blockMap[blockName] = NewBitCaskBlock(path)
	}
}

func (e *BitCaskEngine) Get(BlockName string, Key []byte) ([]byte, error) {
	if _, ok := e.blockMap[BlockName]; !ok {
		return nil, errors.New(fmt.Sprintf("Database %s is not Exist ", BlockName))
	}
	return e.blockMap[BlockName].Get(Key)
}

func (e *BitCaskEngine) Put(BlockName string, Key []byte, Value []byte) error {
	if _, ok := e.blockMap[BlockName]; !ok {
		return errors.New(fmt.Sprintf("Database %s is not Exist ", BlockName))
	}
	return e.blockMap[BlockName].Put(Key, Value)
}

func (e *BitCaskEngine) Del(BlockName string, Key []byte) error {
	if _, ok := e.blockMap[BlockName]; !ok {
		return errors.New(fmt.Sprintf("Database %s is not Exist ", BlockName))
	}
	return e.blockMap[BlockName].Del(Key)
}

func (e *BitCaskEngine) Block(BlockName string) (*BitCaskBlock, error) {
	if _, ok := e.blockMap[BlockName]; !ok {
		return nil, errors.New("qqq")
	}
	return e.blockMap[BlockName], nil
}

type BitCaskBlock struct {
	FM   *store.FileManager
	path string
}

func NewBitCaskBlock(path string) *BitCaskBlock {
	b := &BitCaskBlock{path: path}
	meta, _ := store.NewMeta(path)
	fmp, _ := store.NewFileManager(meta)
	b.FM = fmp
	return b
}

func (b *BitCaskBlock) Name() string {
	l := strings.Split(b.path, string(os.PathSeparator))
	return l[len(l)-1]
}

func (b *BitCaskBlock) Get(Key []byte) ([]byte, error) {
	value, err := b.FM.Read(Key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BitCaskBlock) Put(Key []byte, Value []byte) error {
	return b.FM.Write(Key, Value)
}

func (b *BitCaskBlock) Del(Key []byte) error {
	return b.FM.Write(Key, []byte(store.DELETEFLAG))
}
