package main

import (
	"os"
	"time"
)

type FileNameReader struct{}

var _ MediaTimestampReader = FileAttributeReader{}

func (FileNameReader) Source() string {
	return "name"
}

func (FileNameReader) ReadTimestamp(dir string, file os.FileInfo) (bool, time.Time, error) {
	// Parse the name of the file to see it indicates the timestamp
	// Works for DropBox uploads, or Google Phones (and maybe others)
	// Will not work for iPhone or camera
	return false, time.Time{}, nil
}
