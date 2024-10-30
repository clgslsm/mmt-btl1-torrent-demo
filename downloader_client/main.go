// downloader_client.go
package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
	metainfoData, err := os.ReadFile("a.jpg.torrent")
	if err != nil {
		log.Fatalf("Failed to read metainfo file: %v", err)
	}

	var metainfo Metainfo
	if err := json.Unmarshal(metainfoData, &metainfo); err != nil {
		log.Fatalf("Failed to parse metainfo: %v", err)
	}

	// Create output file to reassemble pieces
	output, err := os.Create(metainfo.FileName + "_downloaded.jpg")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer output.Close()

	// Request piece addresses from the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:8000/getPieceAddress?fileName=%s", metainfo.FileName))
	if err != nil {
		log.Fatalf("Failed to get piece addresses: %v", err)
	}
	defer resp.Body.Close()

	var pieceAddresses map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&pieceAddresses); err != nil {
		log.Fatalf("Failed to decode piece addresses: %v", err)
	}

	// Loop through each piece and request it from the provided addresses
	for i, pieceHash := range metainfo.PieceHashes {
		pieceURL, ok := pieceAddresses[fmt.Sprintf("%d", i)]
		if !ok {
			log.Fatalf("No address found for piece %d", i)
		}

		// Log the piece URL
		fmt.Printf("Downloading piece %d from URL: %s\n", i, pieceURL)
		resp, err := http.Get(fmt.Sprintf("%s/%s_piece_%d", pieceURL, metainfo.FileName, i))
		if err != nil {
			log.Fatalf("Failed to download piece %d: %v", i, err)
		}
		defer resp.Body.Close()

		pieceContent, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read piece content: %v", err)
		}

		// Verify the piece hash using SHA-1
		if actualHash := fmt.Sprintf("%x", sha1.Sum(pieceContent)); actualHash != pieceHash {
			log.Fatalf("Hash mismatch for piece %d: expected %s, got %s", i, pieceHash, actualHash)
		}

		output.Write(pieceContent)
		fmt.Printf("Downloaded and wrote piece %d with hash %s\n", i, pieceHash)
	}

	fmt.Println("File download and reassembly complete!")
}
