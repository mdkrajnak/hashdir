package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileHash stores the filename and its corresponding hash
type FileHash struct {
	Name string
	Hash string
}

// Command line flags
var (
	dirPath   string
	useSHA512 bool
)

func init() {
	// Define both short and long forms of flags
	flag.StringVar(&dirPath, "dir", ".", "Directory path to scan")
	flag.StringVar(&dirPath, "d", ".", "Directory path to scan (shorthand)")

	flag.BoolVar(&useSHA512, "sha512", false, "Use SHA512 instead of SHA256")
	flag.BoolVar(&useSHA512, "s", false, "Use SHA512 instead of SHA256 (shorthand)")

	// Customize flag usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		fmt.Fprintf(os.Stderr, "  -d, --dir string\n\tDirectory path to scan (default \".\")\n")
		fmt.Fprintf(os.Stderr, "  -s, --sha512\n\tUse SHA512 instead of SHA256\n")
	}
}

func main() {
	// Process args to handle --flag syntax
	processArgs()
	flag.Parse()

	// Get list of files in directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	// Create slice to store file hashes
	var fileHashes []FileHash

	// Process each file
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirPath, file.Name())
		hash, err := computeHash(fullPath, useSHA512)
		if err != nil {
			fmt.Printf("Error computing hash for %s: %v\n", file.Name(), err)
			continue
		}

		fileHashes = append(fileHashes, FileHash{
			Name: file.Name(),
			Hash: hash,
		})
	}

	// Sort fileHashes by filename
	sort.Slice(fileHashes, func(i, j int) bool {
		return fileHashes[i].Name < fileHashes[j].Name
	})

	// Print results
	hashType := "SHA256"
	if useSHA512 {
		hashType = "SHA512"
	}
	fmt.Printf("\nFile Hashes (%s):\n", hashType)
	fmt.Println("----------------------------------------")
	for _, fh := range fileHashes {
		fmt.Printf("%s: %s\n", fh.Name, fh.Hash)
	}
}

// processArgs converts --flag style args to -flag style
func processArgs() {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--") {
			os.Args[i] = "-" + strings.TrimPrefix(arg, "--")
		}
	}
}

func computeHash(filePath string, useSHA512 bool) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var hash string
	if useSHA512 {
		hasher := sha512.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		hash = hex.EncodeToString(hasher.Sum(nil))
	} else {
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		hash = hex.EncodeToString(hasher.Sum(nil))
	}

	return hash, nil
}
