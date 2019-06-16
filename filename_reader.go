package main

import (
	"os"
	"regexp"
	"strings"
	"time"
)

type FileNameReader struct{}

var _ MediaTimestampReader = FileAttributeReader{}

func (FileNameReader) Source() string {
	return "name"
}

func (FileNameReader) ReadTimestamp(dir string, file os.FileInfo) (bool, time.Time, error) {
	// Parse the name of the file to see if it indicates the timestamp
	// Works for DropBox and some Phones
	// Will not work for iPhone or a camera
	// Other phones probably use different name formats which could be added here
	var dateFormats = []struct {
		r string // regExp to extract date to parse
		f string // go format string to parse result
	}{
		{`^\d{4}-\d\d-\d\d \d\d\.\d\d\.\d\d`, "2006-01-02 15.04.05"}, // Dropbox format
		{`^vid_\d{8}_\d{6}`, "vid_20060102_150405"},                  // Pixel format
		{`^img_\d{8}_\d{6}`, "img_20060102_150405"},                  // Pixel format
		{`^pano_\d{8}_\d{6}`, "pano_20060102_150405"},                // Pixel format
	}

	// extract the potential date time part of a file name using a regEx and try to parse it
	name := strings.ToLower(file.Name())
	for _, f := range dateFormats {
		re := regexp.MustCompile(f.r)
		n := re.FindString(name)
		if n != "" {
			t, err := time.Parse(f.f, n)
			if err == nil {
				return true, t, nil
			}
		}
	}

	return false, time.Time{}, nil
}
