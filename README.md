![sticth](./stitch.png)

Go command line tool for "stitching" files together. Provide any type of file and stitch will do its best to stitch them together.

If given a folder, stitch will only take the root level files of that folder into consideration when stitching (ie. it doesn't recursively stitch the contents of nested folders together).

> Stitch doesn't make any assumptions about the format or content of your files. It's up to you to be aware of this when stitching files together and determining the final output file type.

### Installing

You can either clone this repo and build from source with `go build` or you can download a release.

### Usage

```bash
Usage:
  stitch [file/path file/path] [flags]

Flags:
  -d, --debug              Debug Output
  -e, --extension string   Output file extension, with dot (.) (default ".md")
  -h, --help               help for stitch
  -o, --output string      Output directory (default "./")
  -s, --strategy string    Stitch strategy [permute, sequence] (default "permute") [NOT IMPLEMENTED]
  -v, --verbose            Verbose Output
```

> NOTE: "strategy" is not yet implemented. Stitch does a "permute" by default, meaning it will stitch together all possible permutations of files in a directory with single files (see verbose output example below)

### Examples

```
# Stitch two files together
$ » stitch intro.md outro.md

# Stitch a file and the contents of a directory together
$ » stitch intro.md content/

# Stitch a file, the contents of a directory, and another file together
$ » stitch intro.md content/ outro.md

# Stitch a file, a specific file in a directory, and another file together
$ » stitch intro.md content/body.md outro.md

# Change the file extension
$ » stitch intro.md outro.md -e .txt

# Variable file extensions
$ » stitch text.txt markdown.md styles.css index.html -e .html

# Change the output file location
$ » stitch intro.md outro.md -o /path/to/output/dir

# Verbose Output
$ » stitch colors/ eggs-and-ham.md -v
Found args: '[colors/ eggs-and-ham.md]'
Adding: '/Users/sdomino/Desktop/example/colors/blue.md'
Adding: '/Users/sdomino/Desktop/example/colors/green.md'
Adding: '/Users/sdomino/Desktop/example/colors/red.md'
Adding: '/Users/sdomino/Desktop/example/eggs-and-ham.md'
Complete: './blue-eggs-and-ham.md'
Complete: './green-eggs-and-ham.md'
Complete: './red-eggs-and-ham.md'
```

## TODO

- [] Add more tests
- [] Ability to specify final output file name
- [] Variable injection
