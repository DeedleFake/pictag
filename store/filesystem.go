package store

import (
	"io/fs"
	"path/filepath"
)

type filesystem struct {
	store *Store
}

func (f filesystem) Open(name string) (fs.File, error) {
	// TODO: Allow opening a fake root directory that can list every image?

	if !validName(name) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}

	file, err := f.store.root.Open(filepath.Join(name[:2], name))
	if err != nil {
		return nil, err
	}
	return file, nil
}
