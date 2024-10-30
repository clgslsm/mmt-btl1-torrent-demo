// server.go
package main

import (
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

type PieceInfo struct {
	Client  string `json:"client"`
	Address string `json:"address"`
}

var registeredMetainfos = map[string]Metainfo{}

type FilePieces map[string]map[string][]PieceInfo

func getPieceAddress(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")

	if fileName == "" {
		http.Error(w, "Missing fileName", http.StatusBadRequest)
		return
	}

	file, err := os.Open("file.json")
	if err != nil {
		http.Error(w, "Could not open JSON file", http.StatusInternalServerError)
		return
	}

	defer file.Close()
	var filePieces FilePieces
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &filePieces)

	if pieces, ok := filePieces[fileName]; ok {
		addresses := make(map[string]string)
		for pieceID, pieceInfo := range pieces {
			if len(pieceInfo) > 0 {
				// Collect the first available address for each piece
				addresses[pieceID] = pieceInfo[0].Address
			}
		}
		if len(addresses) > 0 {
			json.NewEncoder(w).Encode(addresses)
			return
		}
	}
	http.Error(w, "Pieces not found", http.StatusNotFound)
}

func registerMetainfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var metainfo Metainfo
		if err := json.NewDecoder(r.Body).Decode(&metainfo); err != nil {
			http.Error(w, "Invalid data", http.StatusBadRequest)
			return
		}

		// Use file hash (or other unique identifier) as the key
		fileKey := metainfo.FileName
		registeredMetainfos[fileKey] = metainfo

		fmt.Printf("Registered metainfo: %+v\n", metainfo)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Metainfo registered successfully"))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/register_metainfo", registerMetainfoHandler)
	http.HandleFunc("/getPieceAddress", getPieceAddress)
	fmt.Println("Tracker listening on port 8000...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
