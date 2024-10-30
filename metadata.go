// metadata.go
package main

type FileMetadata struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"` // File size in bytes
	FileHash string `json:"file_hash"` // A unique hash representing the file
}
