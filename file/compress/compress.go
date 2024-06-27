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
		header := &zip.FileHeader{
			Name:   item.FileName,
			Flags:  1 << 11, // 指定 utf-8 编码
			Method: zip.Deflate,
		}
		if len(password) > 0 {
			header.SetPassword(password)
		}
		w, err := writer.CreateHeader(header)
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
