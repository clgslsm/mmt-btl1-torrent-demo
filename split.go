// metainfo_generator.go
package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Metainfo struct {
	TrackerURL  string   `json:"tracker_url"`
	FileName    string   `json:"file_name"`
	FileSize    int64    `json:"file_size"`
	PieceLength int      `json:"piece_length"`
	PieceHashes []string `json:"piece_hashes"`
}

func main() {
	filePath := "a.jpg"      // Original file to be split and shared
	pieceLength := 52 * 1024 // 52KB piece size

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	var pieceHashes []string

	buffer := make([]byte, pieceLength)
	pieceIndex := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		// Save each piece as a separate file
		pieceFileName := fmt.Sprintf("%s_piece_%d", filePath, pieceIndex)
		pieceFile, err := os.Create(pieceFileName)
		if err != nil {
			log.Fatalf("Failed to create piece file: %v", err)
		}
		pieceFile.Write(buffer[:bytesRead])
		pieceFile.Close()

		// Calculate and store SHA-1 hash for the piece
		hash := sha1.Sum(buffer[:bytesRead])
		pieceHashes = append(pieceHashes, fmt.Sprintf("%x", hash))

		pieceIndex++
	}

	// Create the .torrent file with all metadata
	metainfo := Metainfo{
		TrackerURL:  "http://localhost:8000",
		FileName:    filePath,
		FileSize:    fileSize,
		PieceLength: pieceLength,
		PieceHashes: pieceHashes,
	}

	torrentFileName := filePath + ".torrent"
	metainfoFile, err := os.Create(torrentFileName)
	if err != nil {
		log.Fatalf("Failed to create metainfo file: %v", err)
	}
	defer metainfoFile.Close()

	encoder := json.NewEncoder(metainfoFile)
	if err := encoder.Encode(metainfo); err != nil {
		log.Fatalf("Failed to encode metainfo: %v", err)
	}

	fmt.Printf("Metainfo (.torrent) file created as %s, with %d pieces.\n", torrentFileName, len(pieceHashes))
}
