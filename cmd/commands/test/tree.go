package test

import (
	"crypto/rand"
	"fmt"
	"io"
	mathrand "math/rand"
	"os"
	"path/filepath"
)

// CreateTestTree recursively creates a tree of directories and files under 'root',
// with random file sizes, until the total size reaches maxSizeMB (in megabytes).
func CreateTestTree(root string, maxSizeMB int) error {
	const (
		minFileSize = 1024        // 1KB
		maxFileSize = 1024 * 1024 // 1MB
	)
	maxSize := int64(maxSizeMB) * 1024 * 1024
	var totalSize int64

	var create func(path string) error
	create = func(path string) error {
		if mathrand.Float64() < 0.7 || totalSize+minFileSize > maxSize {
			remaining := maxSize - totalSize
			if remaining < minFileSize {
				return nil
			}
			fileSize := int64(mathrand.Intn(maxFileSize-minFileSize+1) + minFileSize)
			if fileSize > remaining {
				fileSize = remaining
			}
			fileName := fmt.Sprintf("file_%d.dat", mathrand.Intn(1e6))
			filePath := filepath.Join(path, fileName)
			f, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := io.CopyN(f, rand.Reader, fileSize); err != nil {
				return err
			}
			totalSize += fileSize
		} else {
			dirName := fmt.Sprintf("dir_%d", mathrand.Intn(1e6))
			dirPath := filepath.Join(path, dirName)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return err
			}
			for i := 0; i < mathrand.Intn(4)+1; i++ {
				if totalSize >= maxSize {
					break
				}
				if err := create(dirPath); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	for totalSize < maxSize {
		if err := create(root); err != nil {
			return err
		}
	}
	return nil
}
