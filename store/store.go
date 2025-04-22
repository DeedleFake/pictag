// Package store implements a simple hash-addressed storage system for images.
package store

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
)

var ErrInvalidName = errors.New("invalid name")

// Store is a hash-addressed on-disk storage system for images. It
// converts images to a standard format and writes them to disk,
// deduplicating them by the hash of the resulting data. It can then
// retreive those images by identifying them by the resulting hash.
type Store struct {
	root *os.Root
}

// Open opens a store rooted at path. The directory at path must
// already exist. No checks are done for data in an existing
// directory, so the store will just write data right alongside
// anything unrelated that is already there.
func Open(path string) (*Store, error) {
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}

	return &Store{
		root: root,
	}, nil
}

// Close closes the store.
func (s *Store) Close() error {
	return s.root.Close()
}

func (s *Store) store(name string, data []byte) error {
	prefix := name[:2]
	err := s.root.Mkdir(prefix, 0644)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	w, err := s.root.Create(filepath.Join(prefix, name))
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write(data)
	return err
}

// Store stores the given image into s, returning its identifying hash.
func (s *Store) Store(img image.Image) (string, error) {
	h := sha256.New()
	var buf bytes.Buffer // TODO: Move this to the disk in case it's huge?
	w := io.MultiWriter(h, &buf)
	err := png.Encode(w, img)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf("%x", h.Sum(nil))
	return name, s.store(name, buf.Bytes())
}

// Load returns the image identified by name, if it exists. If name is
// invalid, an error equal to ErrInvalidName is returned.
func (s *Store) Load(name string) (image.Image, error) {
	if len(name) != 64 {
		return nil, fmt.Errorf("%w: %q", ErrInvalidName, name)
	}

	file, err := s.root.Open(filepath.Join(name[:2], name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}
