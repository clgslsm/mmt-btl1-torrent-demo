package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	// Check if port number is provided
	if len(os.Args) < 2 {
		log.Fatal("Port number is required. Usage: go run client.go <port>")
	}
	port := os.Args[1]

	// Load .torrent metadata
	metainfoData, _ := ioutil.ReadFile("a.jpg.torrent")
	var metainfo Metainfo
	json.Unmarshal(metainfoData, &metainfo)

	// // Register pieces with the tracker
	// for i := range metainfo.PieceHashes {
	// 	piece := map[string]interface{}{
	// 		"file_name": metainfo.FileName,
	// 		"piece_id":  i,
	// 		"peer_addr": "http://localhost:5000",
	// 	}
	// 	pieceData, _ := json.Marshal(piece)
	// 	http.Post(fmt.Sprintf("%s/register_piece", metainfo.TrackerURL), "application/json", bytes.NewReader(pieceData))
	// }

	// Serve pieces at /{file_name}_piece_{piece_id}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the piece file name from the URL path
		pieceFileName := r.URL.Path[1:] // Remove the leading '/'
		data, err := ioutil.ReadFile(pieceFileName)
		if err != nil {
			http.Error(w, "Piece not found", http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	fmt.Printf("Client A serving pieces on port %s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
