// TODO: add a "strategy" where you can either do "permutations" or "sequential"

// Package cmd ...
package cmd

import (
	"fmt"
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
	strategy  string // "stitch" strategy
	verbose   bool   // verbose output

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

	// create the final output directory
	createOutputDir(output)

	// the final list of all files to build
	buildlist := [][]string{}

	// loop through "args" to build "buildlist"; if an arg is a file it's added to
	// buildlist as []string{path} if it's a directory it's looped through adding
	// file to buildlist as []string{path}
	for _, arg := range args {

		// get file info for the current "arg" (should be either a file or dir)
		fi := getFileInfo(arg)

		// determine if "arg" is a file or directory
		switch mode := fi.Mode(); {

		// if a dir, get all files and range through looking for files only to add to
		// buildlist
		case mode.IsDir():

			// get all the files in the directory
			files := readDir(arg)

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
		tmpf := createFile(tmpfile)
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
			contents := readFile(path)

			// write contents to tmpfile
			writeFile(tmpf, []byte(fmt.Sprintf("%s\n", contents)))
		}

		final := fmt.Sprintf("%s/%s%s", filepath.Dir(tmpfile), strings.TrimPrefix(finalname, "-"), extension)

		// convert tmpfile to the final file
		renameFile(tmpfile, final)

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
