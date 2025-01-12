package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	hashgen "hashdir/hashgen"
)

// Command line flags
var (
	leftDirPath  string
	rightDirPath string
	useSHA512    bool
)

func init() {
	// Define both short and long forms of flags
	flag.StringVar(&leftDirPath, "left", ".", "Left directory path to scan")
	flag.StringVar(&leftDirPath, "l", ".", "Left directory path to scan (shorthand)")

	flag.StringVar(&rightDirPath, "right", ".", "Right directory path to scan")
	flag.StringVar(&rightDirPath, "r", ".", "Right directory path to scan (shorthand)")

	flag.BoolVar(&useSHA512, "sha512", false, "Use SHA512 instead of SHA256")
	flag.BoolVar(&useSHA512, "s", false, "Use SHA512 instead of SHA256 (shorthand)")

	// Customize flag usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fmt.Fprintf(os.Stderr, "  -l, --left string\n\tLeft directory path to scan (default \".\")\n")
		fmt.Fprintf(os.Stderr, "  -r, --right string\n\tRight directory path to scan (default \".\")\n")
		fmt.Fprintf(os.Stderr, "  -s, --sha512\n\tUse SHA512 instead of SHA256\n")
	}
}

func main() {
	// Process args to handle --flag syntax
	processArgs()
	flag.Parse()

	// Generate file hashes for both directories
	leftFileHashes, err := hashgen.GenerateFileHashes(leftDirPath, useSHA512)
	if err != nil {
		fmt.Printf("Error generating file hashes for left directory: %v\n", err)
		os.Exit(1)
	}

	rightFileHashes, err := hashgen.GenerateFileHashes(rightDirPath, useSHA512)
	if err != nil {
		fmt.Printf("Error generating file hashes for right directory: %v\n", err)
		os.Exit(1)
	}

	// Compare file hashes
	identicalFiles, differentFiles, leftOnlyFiles, rightOnlyFiles := compareFileHashes(leftFileHashes, rightFileHashes)

	// Print results
	fmt.Printf("\nComparison Results:\n")
	fmt.Printf("Identical Files: %d\n", identicalFiles)
	fmt.Printf("Different Files: %d\n", differentFiles)
	fmt.Printf("Files only in left directory: %d\n", leftOnlyFiles)
	fmt.Printf("Files only in right directory: %d\n", rightOnlyFiles)
}

// processArgs converts --flag style args to -flag style
func processArgs() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			os.Args[i] = "-" + strings.TrimPrefix(arg, "--")
		}
	}
}

// compareFileHashes compares two slices of FileHash and returns detailed comparison results
func compareFileHashes(left, right []hashgen.FileHash) (int, int, int, int) {
	leftMap := make(map[string]string)
	for _, fh := range left {
		leftMap[fh.Name] = fh.Hash
	}

	rightMap := make(map[string]string)
	for _, fh := range right {
		rightMap[fh.Name] = fh.Hash
	}

	identicalFiles := 0
	differentFiles := 0
	leftOnlyFiles := 0
	rightOnlyFiles := 0

	// Check files in the left directory
	for name, leftHash := range leftMap {
		if rightHash, exists := rightMap[name]; exists {
			if leftHash == rightHash {
				identicalFiles++
			} else {
				differentFiles++
			}
		} else {
			leftOnlyFiles++
		}
	}

	// Check files in the right directory that are not in the left directory
	for name := range rightMap {
		if _, exists := leftMap[name]; !exists {
			rightOnlyFiles++
		}
	}

	return identicalFiles, differentFiles, leftOnlyFiles, rightOnlyFiles
}
