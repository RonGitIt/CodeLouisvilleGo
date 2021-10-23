package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello world"))
	} )
	http.HandleFunc("/upload", handleImageUpload)
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleImageUpload(w http.ResponseWriter, r *http.Request){
	defer func() {fmt.Println("Exiting handleImageUpload")}()

	outputDir := "/home/jon/tmp"
	switch r.Method {
	case http.MethodPost:
		r.ParseMultipartForm(32 << 20)

		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error getting file from request: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		io.Copy(&buf, file)

		outputFilename := filepath.Join(outputDir, header.Filename)
		outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0666)
		defer outputFile.Close()
		io.Copy(outputFile, &buf)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("File size: %v\n", header.Size)))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
