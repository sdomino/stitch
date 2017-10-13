// Package cmd ...
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (

	// vars
	version = "0.0.2" //

	// flags
	debug     bool   // debug output
	extension string // final file extention (.md)
	output    string // output/directory/path
	// strategy  string // "stitch" strategy
	verbose bool // verbose output

	// StitchCmd -
	StitchCmd = &cobra.Command{
		Use:   "stitch [file/path file/path]",
		Short: "",
		Long:  ``,

		//
		// print version or help, or continue, depending on flag settings
		PreRun: prerun,

		//
		RunE: stitch,
	}
)

//
func init() {
	StitchCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug Output")
	StitchCmd.Flags().StringVarP(&extension, "extension", "e", ".md", "Output file extension, with dot (.)")
	StitchCmd.Flags().StringVarP(&output, "output", "o", "./", "Output directory")
	// StitchCmd.Flags().StringVarP(&strategy, "strategy", "s", "permute", "Stitch strategy [permute, sequence]")
	StitchCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose Output")
}

//
func prerun(cmd *cobra.Command, args []string) {

	// if no args are given print help
	if len(args) == 0 {
		cmd.HelpFunc()(cmd, args)
		os.Exit(0)
	}

	// print current version of stitch
	if args[0] == "version" {
		fmt.Printf("stitch %s\n", version)
		os.Exit(0)
	}
}

//
func stitch(cmd *cobra.Command, args []string) error {

	printv(fmt.Sprintf("Found: '%s'", args))

	// create the final output directory; will attempt to make "output" dir if it
	// doesn't exist
	if err := os.MkdirAll(output, 0777); err != nil {
		fatal(err.Error())
	}

	// loop through "args" to build "buildlist"; if an arg is a file it's added to
	// buildlist as []string{path} if it's a directory it's looped through adding
	// file to buildlist as []string{path}
	buildlist := [][]string{}
	for _, arg := range args {

		// get file info for the current "arg" (should be either a file or dir)
		fi, err := os.Stat(arg)
		if err != nil {
			fatal(err.Error())
		}

		// determine if "arg" is a file or directory
		switch mode := fi.Mode(); {

		// if a dir, get all files and range through looking for files only to add to
		// buildlist
		case mode.IsDir():

			// get all the files in the directory
			files, err := ioutil.ReadDir(arg)
			if err != nil {
				fatal(err.Error())
			}

			// range through the files looking only for files to add to the buildlist
			filelist := []string{}
			for _, file := range files {
				if file.Mode().IsRegular() {
					filepath := getAbsolutePath(fmt.Sprintf("%s/%s", arg, file.Name()))
					printv(fmt.Sprintf("Adding: '%s'", filepath))
					filelist = append(filelist, filepath)
				}
			}

			// add files to buildlist
			buildlist = append(buildlist, filelist)

		// if a file, add it to buildlist as []string{path}
		case mode.IsRegular():
			filepath := getAbsolutePath(arg)
			printv(fmt.Sprintf("Adding: '%s'", filepath))
			buildlist = append(buildlist, []string{filepath})
		}
	}

	// the "stitch" file. A tmp file where all contents are appended
	tmpfile := fmt.Sprintf("%s/stitch.tmp", output)

	// range over buildlist to get all permutations of files
	for _, list := range getPermutation(buildlist) {

		// create tmpfile file at "output" to write contents to
		tmpf, err := os.OpenFile(tmpfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fatal(err.Error())
		}
		defer tmpf.Close()

		// the name of the final output file; format: name-of-each-file.ext
		finalname := ""

		// range over the list of files and start "stitching"
		for _, path := range list {

			// get the base name of the file and split on "." to get [name extention]
			filename := strings.Split(filepath.Base(path), ".")

			// build finalname using filenames
			finalname += fmt.Sprintf("-%s", filename[0])

			// read contents of current file
			contents, err := ioutil.ReadFile(path)
			if err != nil {
				fatal(err.Error())
			}

			// write contents to tmpfile
			if _, err := tmpf.Write([]byte(fmt.Sprintf("%s\n", contents))); err != nil {
				fatal(err.Error())
			}
		}

		final := fmt.Sprintf("%s/%s%s", filepath.Dir(tmpfile), strings.TrimPrefix(finalname, "-"), extension)

		// convert tmpfile to the final file
		if err := os.Rename(tmpfile, final); err != nil {
			fatal(err.Error())
		}

		printv(fmt.Sprintf("Complete: '%s'", final))
	}

	return nil
}

//
func printv(output string) {
	if verbose {
		fmt.Printf("%s\n", output)
	}
}

//
func printd(output string) {
	if debug {
		fmt.Printf("DEBUG: %s\n", output)
	}
}

//
func fatal(msg string) {
	fmt.Printf("ERROR: %s\n", msg)
	os.Exit(1)
}

// getAbsolutePath
func getAbsolutePath(path string) string {
	abspath, err := filepath.Abs(path)
	if err != nil {
		fatal(err.Error())
	}
	return abspath
}

// getPermutation will get all permutations of a slice of string slices. Adapted
// from https://stackoverflow.com/a/43973743
func getPermutation(files [][]string) [][]string {
	permutation := [][]string{}

	// return if we got an empty slice
	if len(files) == 0 {
		return nil
	}

	// if there is only one slice element to be permuted
	if len(files) == 1 {
		// each sub-element is a permutation
		for i := range files[0] {
			permutation = append(permutation, []string{files[0][i]})
		}
		return permutation
	}

	// build from right to left
	t := getPermutation(files[1:])
	// append permutations of the last elements to this element
	for i := range files[0] {
		for x := range t {
			permutation = append(permutation, append([]string{files[0][i]}, t[x]...))
		}
	}

	return permutation
}
