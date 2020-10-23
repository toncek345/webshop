package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ Storage = (*Disk)(nil)

type Disk struct {
	DirPath string
}

func (d *Disk) Put(key string, data []byte) (*File, error) {
	if err := ioutil.WriteFile(filepath.Join(d.DirPath, key), data, 0777); err != nil {
		return nil, err
	}

	return &File{
		Key:  key,
		Data: data,
	}, nil
}

func (d *Disk) Get(key string) (*File, error) {
	data, err := ioutil.ReadFile(filepath.Join(d.DirPath, key))
	if err != nil {
		return nil, err
	}

	return &File{
		Key:  key,
		Data: data,
	}, nil
}

func (d *Disk) Delete(key string) error {
	if err := os.Remove(filepath.Join(d.DirPath, key)); err != nil {
		return err
	}

	return nil
}
