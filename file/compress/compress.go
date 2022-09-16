package compress

import (
	"bytes"
	"io"

	"github.com/alexmullins/zip"
)

type FileItem struct {
	FileName string
	Reader   io.Reader
}

func ZipBuffer(password string, fileItems ...FileItem) (*bytes.Buffer, error) {
	out := new(bytes.Buffer)
	writer := zip.NewWriter(out)
	defer writer.Close()
	for _, item := range fileItems {
		var (
			err error
			w   io.Writer
		)
		if len(password) > 0 {
			w, err = writer.Encrypt(item.FileName, password)
		} else {
			w, err = writer.Create(item.FileName)
		}
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(w, item.Reader); err != nil {
			return nil, err
		}
		writer.Flush()
	}
	return out, nil
}
