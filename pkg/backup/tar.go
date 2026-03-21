package backup

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/fs"
	"time"
)

// TarGzFS is an ihfs.WriteFileFS and ihfs.MkdirAllFS backed by a tar.gz archive.
// Call Close to flush and close both the tar and gzip writers.
type TarGzFS struct {
	gz  *gzip.Writer
	tar *tar.Writer
}

// NewTarGzFS creates a new TarGzFS writing a tar.gz archive to w.
func NewTarGzFS(w io.Writer) *TarGzFS {
	gz := gzip.NewWriter(w)
	return &TarGzFS{gz: gz, tar: tar.NewWriter(gz)}
}

func (f *TarGzFS) Open(_ string) (fs.File, error) {
	return nil, fs.ErrInvalid
}

func (f *TarGzFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	hdr := &tar.Header{
		Name:    name,
		Mode:    int64(perm),
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}
	if err := f.tar.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := f.tar.Write(data)
	return err
}

// MkdirAll is a no-op as directories are implicit in tar archives.
func (f *TarGzFS) MkdirAll(_ string, _ fs.FileMode) error {
	return nil
}

// Close flushes and closes the tar and gzip writers.
// The underlying io.Writer is not closed.
func (f *TarGzFS) Close() error {
	if err := f.tar.Close(); err != nil {
		return err
	}
	return f.gz.Close()
}
