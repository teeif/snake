package snake

import (
	"archive/tar"
	"bytes"
	"io/fs"
	"path/filepath"

	"github.com/dsnet/compress/bzip2"
)

type Tarlib struct {
	Buffer   *bytes.Buffer
	FS       *tar.Writer
	Gzip     *bzip2.Writer
	FileName string
}

func Tar(tarfile string) *Tarlib {
	t := new(Tarlib)
	t.Buffer = new(bytes.Buffer)
	t.Gzip, _ = bzip2.NewWriter(t.Buffer, &bzip2.WriterConfig{
		Level: 9,
	})
	t.FS = tar.NewWriter(t.Gzip)
	t.FileName = tarfile
	return t
}

func (t *Tarlib) Add(path string, stat fs.FileInfo, body []byte) bool {
	if !String(path).Find(".DS_Store", true) && !String(path).Find("__MACOSX", true) && !String(path).Find(".gitignore", true) && !String(path).Find(".index", true) {
		header, _ := tar.FileInfoHeader(stat, path)
		header.Name = filepath.ToSlash(path)
		header.Size = int64(len(body))
		if err := t.FS.WriteHeader(header); err == nil {
			_, err := t.FS.Write(body)
			return err == nil
		}
	}
	return false
}

func (t *Tarlib) Close() error {
	t.FS.Close()
	t.Gzip.Close()
	_, err := FS(t.FileName).ByteWriter(t.Buffer.Bytes())
	return err
}
