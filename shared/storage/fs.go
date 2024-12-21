package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sst "github.com/picosh/pobj/storage"
)

type StorageFS struct {
	*sst.StorageFS
}

func NewStorageFS(dir string) (*StorageFS, error) {
	st, err := sst.NewStorageFS(dir)
	if err != nil {
		return nil, err
	}
	return &StorageFS{st}, nil
}

func (s *StorageFS) ServeObject(bucket sst.Bucket, fpath string, opts *ImgProcessOpts) (io.ReadCloser, *sst.ObjectInfo, error) {
	var rc io.ReadCloser
	var info *sst.ObjectInfo
	var err error
	mimeType := GetMimeType(fpath)
	if !strings.HasPrefix(mimeType, "image/") || opts == nil || os.Getenv("IMGPROXY_URL") == "" {
		rc, info, err = s.GetObject(bucket, fpath)
		// StorageFS never returns a content-type.
		info.Metadata.Set("content-type", mimeType)
	} else {
		filePath := filepath.Join(bucket.Name, fpath)
		dataURL := fmt.Sprintf("s3://%s", filePath)
		rc, info, err = HandleProxy(dataURL, opts)
	}
	if err != nil {
		return nil, nil, err
	}
	return rc, info, err
}
