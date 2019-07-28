package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type Config struct {
	// SourceDirectory is the directory to be scanned for media files to be moved
	SourceDirectory string
	// DestinationDirectory specifies the root directory for media files to be moved to
	// Files are arranged by date under this location. Example:
	//   <rootdir>/2019/2019-02 February
	//   <rootdir>/2019/2019-03 March
	DestinationDirectory string
	// ArchiveDirectory is an optional location to move file into after copying and renaming
	// The original filename is not changed in the move
	ArchiveDirectory string
}

type MediaInfo struct {
	// Our best guess at the timestamp this media was taken
	Time time.Time
	// Source of the timestamp extracted
	TimeSource string
}

type MediaTimestampReader interface {
	// ReadTimestamp will read the file timestamp for a given media file
	// The definition of this will depend on the implementation
	// Returns:
	// bool 		- was a date returned (if false, not supported format or error occurred)
	// time.Time	- the date (of Zero date if none)
	// error		- an error if time should have been available but failed
	ReadTimestamp(dir string, file os.FileInfo) (bool, time.Time, error)
	// Source returns a string describing the source of the timestamp
	Source() string
}

// The following readers are all tried in the given order to extract a timestamp for
// a given media file
var timestampReaders = []MediaTimestampReader{
	ExifReader{},
	FileNameReader{},
	FileAttributeReader{},
}

func main() {
	var config = defaultConfig
	ProcessDirectory(config.SourceDirectory, config.DestinationDirectory, config.ArchiveDirectory)
}

func ProcessDirectory(srcDir string, baseDestDir string, archiveDir string) {

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		panic(err)
	}
	archiveCreated := false
	for _, file := range files {
		var info MediaInfo
		var err error

		if skipFile(file.Name()) {
			continue
		}

		if !isVideo(file.Name()) && !isImage(file.Name()) {
			fmt.Printf("SKIPPING file %s: unsupported file type\n", path.Join(srcDir, file.Name()))
			continue
		}

		if info.Time, info.TimeSource, err = getTimestamp(srcDir, file); err != nil {
			fmt.Printf("SKIPPING file %s: cannot extract date: %v\n", path.Join(srcDir, file.Name()), err)
			continue
		}

		sourcePath := path.Join(srcDir, file.Name())
		destDir := newDestinationDir(baseDestDir, info) // full directory to move  file to
		destFileName := newFileName(file.Name(), info)
		dest := path.Join(destDir, destFileName)
		fmt.Printf("[%s] %s\t\t => %s (%v)...", info.TimeSource, sourcePath, dest, info.Time)
		actualDest, err := copyFile(sourcePath, dest)
		if err != nil {
			// skip this file and leave it where it is
			fmt.Printf("FAILED copy to %s (%v)\n", actualDest, err)
			continue
		}
		if archiveDir != "" {
			// if an archive directory is supplied we move the original file into it
			if !archiveCreated {
				archiveCreated = true
				archiveDir = path.Join(archiveDir, time.Now().Format("2006-01-02_150405"))
				if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil {
					panic(err)
				}
			}
			if err := os.Rename(sourcePath, path.Join(archiveDir, file.Name())); err != nil {
				panic(fmt.Sprintf(": Failed to move %s to %s (%v)", sourcePath, archiveDir, err))
			}
		}
		fmt.Println(" DONE")
	}
}

func getTimestamp(srcDir string, file os.FileInfo) (time.Time, string, error) {
	var errorStr string
	// Iterate over the timeStamp readers in order until one is successful
	for _, r := range timestampReaders {
		found, tm, err := r.ReadTimestamp(srcDir, file)
		if found {
			return tm, r.Source(), nil
		}
		if err != nil {
			errorStr = errorStr + ": " + err.Error()
		}
	}
	return time.Time{}, "", fmt.Errorf(errorStr)
}

// newFileLocation returns the directory the file should be copied to
func newDestinationDir(baseDir string, fm MediaInfo) string {
	dest := baseDir
	if !fm.Time.IsZero() {
		dest = path.Join(dest, fmt.Sprintf("%d", fm.Time.Year()))
		dest = path.Join(dest, fmt.Sprintf("%d-%.2d %s", fm.Time.Year(), int(fm.Time.Month()), fm.Time.Month().String()))
	}
	return dest
}

func newFileName(currentName string, fm MediaInfo) string {
	if fm.Time.IsZero() {
		panic(fmt.Sprintf("failed to generate new filename for %s as no file time available", currentName))
	}
	prefix := getPrefix(currentName)
	suffix := getSuffix(currentName)
	return fmt.Sprintf("%s%s%s", prefix, fm.Time.Format("2006_01_02_150405"), suffix)
}

func getSuffix(fname string) string {
	suffix := path.Ext(fname)
	return suffix
}

func skipFile(name string) bool {
	if len(name) == 0 {
		// unlikely, but...
		return true
	}
	if strings.HasPrefix(name, ".") {
		return true
	}
	if strings.ToLower(path.Ext(name)) == ".ini" {
		return true
	}
	return false
}

func isImage(name string) bool {
	var imageExt = map[string]struct{}{
		"jpg":  struct{}{},
		"jpeg": struct{}{},
		"tif":  struct{}{},
		"gif":  struct{}{},
	}

	ext := strings.ToLower(path.Ext(name))
	if len(ext) == 0 {
		return false
	}

	if _, ok := imageExt[ext[1:]]; ok {
		// example: ".jpg"
		return true
	}
	return false
}

func isVideo(name string) bool {
	var videoExt = map[string]struct{}{
		"mov": struct{}{},
		"mpg": struct{}{},
		"mp4": struct{}{},
	}

	ext := strings.ToLower(path.Ext(name))
	if len(ext) == 0 {
		return false
	}

	if _, ok := videoExt[ext[1:]]; ok {
		// example: ".mov"
		return true
	}
	return false
}

func getPrefix(name string) string {
	if isImage(name) {
		return "img_"
	}
	if isVideo(name) {
		return "vid_"
	}
	return ""
}
