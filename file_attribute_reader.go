package main

import (
	"fmt"
	"os"
	"time"
)

type FileAttributeReader struct{}

var _ MediaTimestampReader = FileAttributeReader{}

func (FileAttributeReader) Source() string {
	return "ctime"
}

func (FileAttributeReader) ReadTimestamp(dir string, file os.FileInfo) (bool, time.Time, error) {
	// Return the file last update time
	// Ideally we would use file creation time, but this seems to be the best we can do
	// Could possible do something OS dependent to get a better answer, though at best it will
	// be image upload time rather than image taken time.

	if ft := file.ModTime(); !ft.IsZero() {
		return true, ft, nil
	}
	return false, time.Time{}, fmt.Errorf("no file modification time available")
}
