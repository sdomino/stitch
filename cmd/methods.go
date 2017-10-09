package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// createOutputDir will attempt to make "output" dir if it doesn't exist
func createOutputDir(path string) {
	if err := os.MkdirAll(path, 0777); err != nil {
		fatal(err.Error())
	}
}

// getAbsolutePath
func getAbsolutePath(path string) string {
	abspath, err := filepath.Abs(path)
	if err != nil {
		fatal(err.Error())
	}
	return abspath
}

// getFileInfo
func getFileInfo(path string) os.FileInfo {
	fi, err := os.Stat(path)
	if err != nil {
		fatal(err.Error())
	}
	return fi
}

// readDir
func readDir(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fatal(err.Error())
	}
	return files
}

// createFile
func createFile(path string) *os.File {
	tmpf, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fatal(err.Error())
	}
	return tmpf
}

// readFile
func readFile(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fatal(err.Error())
	}
	return contents
}

// writeFile
func writeFile(file *os.File, contents []byte) {
	if _, err := file.Write(contents); err != nil {
		fatal(err.Error())
	}
}

// renameFile
func renameFile(file, name string) {
	if err := os.Rename(file, name); err != nil {
		fatal(err.Error())
	}
}

// getPermutations will get all permutations of a slice of string slices. Adapted
// from https://stackoverflow.com/a/43973743
func getPermutation(things [][]string) [][]string {
	returnVal := [][]string{}

	// return if we got an empty slice
	if len(things) < 1 {
		return nil
	}

	// if there is only one slice element to be permuted
	if len(things) == 1 {
		// each sub-element is a permutation
		for i := range things[0] {
			returnVal = append(returnVal, []string{things[0][i]})
		}
		return returnVal
	}

	// build from right to left
	t := getPermutation(things[1:])
	// append permutations of the last elements to this element
	for i := range things[0] {
		for x := range t {
			returnVal = append(returnVal, append([]string{things[0][i]}, t[x]...))
		}
	}

	return returnVal
}
