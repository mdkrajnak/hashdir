package filehash

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// FileHash stores the filename and its corresponding hash
type FileHash struct {
	Name string
	Hash string
}

// ComputeHash computes the hash of a file
func ComputeHash(filePath string, useSHA512 bool) (string, error) {
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

// GenerateFileHashes generates the file hashes for all files in a directory
func GenerateFileHashes(dirPath string, useSHA512 bool) ([]FileHash, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fileHashes []FileHash
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath := filepath.Join(dirPath, file.Name())
		hash, err := ComputeHash(fullPath, useSHA512)
		if err != nil {
			return nil, err
		}

		fileHashes = append(fileHashes, FileHash{
			Name: file.Name(),
			Hash: hash,
		})
	}

	sort.Slice(fileHashes, func(i, j int) bool {
		return fileHashes[i].Name < fileHashes[j].Name
	})

	return fileHashes, nil
}
