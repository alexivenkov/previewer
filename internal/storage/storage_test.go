package storage

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestStorageSave(t *testing.T) {
	basePath = "testdata"

	f, err := os.Open("testdata/gopher.jpg")
	require.NoError(t, err)
	defer f.Close()

	b, err := io.ReadAll(f)
	require.NoError(t, err)

	storage := NewImageStorage()
	path, err := storage.Save(b, "test.jpg")
	require.NoError(t, err)
	require.Equal(t, "testdata/test.jpg", path)

	tf, err := os.Open(path)
	require.NoError(t, err)
	defer func() {
		tf.Close()
		os.Remove(path)
	}()

	tb, err := io.ReadAll(tf)
	require.NoError(t, err)

	require.Equal(t, b, tb)
}

func TestStorageGet(t *testing.T) {
	basePath = "testdata"

	storage := NewImageStorage()
	b, err := storage.Get("gopher.jpg")
	require.NoError(t, err)

	f, err := os.Open("testdata/gopher.jpg")
	require.NoError(t, err)
	defer f.Close()

	tb, err := io.ReadAll(f)
	require.NoError(t, err)

	require.Equal(t, tb, b)
}

func TestStorageFailGet(t *testing.T) {
	basePath = "testdata"

	storage := NewImageStorage()
	_, err := storage.Get("123.jpg")
	require.ErrorIs(t, err, os.ErrNotExist)
}
