package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

// copyFile copies the file src to dest.
// If a file dest already already exists:
// 		if the existing file has identical contents to this one do nothing and return nil (TODO)
// 		if existing file is not identical then copy to a new filename with a suffix appended
//
func copyFile(src, dest string) (string, error) {

	if src == dest {
		return dest, nil // nothing to do
	}

	source, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer source.Close()

	destDir, _ := path.Split(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return dest, err
	}

	newDest := getDestination(dest)
	destination, err := os.Create(newDest)
	if err != nil {
		return newDest, err
	}
	defer destination.Close()

	// we now have a candidate location to copy the file to
	if _, err = io.Copy(destination, source); err != nil {
		return newDest, err
	}
	return newDest, destination.Sync()
}

// getDestination returns the location to copy the file to
//
func getDestination(dest string) string {

	// TODO: We should check if current version is the same and skip it if so
	dir, fName := path.Split(dest)
	ext := path.Ext(fName)
	base := fName[0 : len(fName)-len(ext)]
	version := 1
	newDest := dest
	for {
		d, err := os.Open(newDest)
		if err != nil {
			// the file failed to open - we'll assume because it doesn't exist!
			break
		}
		// the file already exists - generate a new destination and try again
		_ = d.Close()
		newDest = path.Join(dir, fmt.Sprintf("%s_%d%s", base, version, ext))
		version++
	}

	return newDest
}
