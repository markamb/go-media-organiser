package main

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"os"
	"path"
	"strings"
	"time"
)

type ExifReader struct{}

var _ MediaTimestampReader = ExifReader{}

func (ExifReader) Source() string {
	return "exif"
}

func (ExifReader) ReadTimestamp(dir string, file os.FileInfo) (bool, time.Time, error) {
	fpath := path.Join(dir, file.Name())
	switch {
	case strings.HasSuffix(strings.ToLower(fpath), "jpeg"):
	case strings.HasSuffix(strings.ToLower(fpath), "jpg"):
	default:
		return false, time.Time{}, nil // Not an error - just not supported
	}
	f, err := os.Open(fpath)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("%s : failed to open file: %v\n", fpath, err)
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err != nil {
		// should be supported - failed to read file metadata
		return false, time.Time{}, fmt.Errorf("%s : failed to extract exif date: %v\n", fpath, err)
	}
	tm, err := x.DateTime()
	if err != nil {
		return false, time.Time{}, fmt.Errorf("%s : Failed to extract exif date: %v\n", fpath, err)
	}
	return true, tm, nil
}
