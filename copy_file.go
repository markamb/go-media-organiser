package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// copyFile copies the file src to dest.
// If a file dest already already exists:
// 		if the existing file has identical contents to this one do nothing (returns false, "", nil)
// 		if existing file is not identical then copy to a new filename with a suffix appended (returns true, <new filename>, nil)
//
func copyFile(src, dest string) (bool, string, error) {

	if src == dest {
		return false, dest, nil // nothing to do
	}

	destDir, _ := path.Split(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return false, dest, err
	}

	shouldCopy, newDest, err := getDestination(src, dest)
	if err != nil {
		return false, "", err
	}

	if !shouldCopy {
		return false, "", nil
	}

	destination, err := os.Create(newDest)
	if err != nil {
		return false, newDest, err
	}
	defer destination.Close()

	source, err := os.Open(src)
	if err != nil {
		return false, "", err
	}
	defer source.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return false, newDest, err
	}
	return true, newDest, destination.Sync()
}

// getDestination returns the location to copy the file 'source' to
// Returns true if the file should be copied along with the destination
func getDestination(source, dest string) (bool, string, error) {

	// TODO - we could do a lot better here - comparinmg files a chunk at a time is a bit tricky
	// as reading a given chunk size is not guaranteed, causing chunks to get out of sync between
	// 2 buffers. This is doable, but not worth the extra complications (and therefore risk) given
	// we expect any files we deal with to fit into memory relatively easily.
	// A simple optimisation would be to just compare the 2 file sizes up front however.

	dir, fName := path.Split(dest)
	ext := path.Ext(fName)
	base := fName[0 : len(fName)-len(ext)]
	version := 1
	newDest := dest
	var sBuff []byte
	var sBuffRead bool
	for {
		if d, err := os.Open(newDest); err != nil {
			break // the file failed to open - we'll assume because it doesn't exist!
		} else {
			_ = d.Close()
		}

		// the file already exists - check if it has the same contents as the one we want to copy
		if !sBuffRead {
			var err error
			sBuff, err = ioutil.ReadFile(source)
			if err != nil {
				return false, "", err
			}
			sBuffRead = true
		}

		dBuff, err := ioutil.ReadFile(dest)
		if err != nil {
			return false, "", err
		}

		if bytes.Equal(sBuff, dBuff) {
			return false, "", nil // skip copy - the file already exists
		}

		newDest = path.Join(dir, fmt.Sprintf("%s_%d%s", base, version, ext))
		version++
	}

	return true, newDest, nil
}
