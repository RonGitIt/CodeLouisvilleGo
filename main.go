package main

import (
	"fmt"
	"log"
	"net/http"
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
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("That was a post to handleImageUpload"))
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error getting file from request: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()
		w.Write([]byte(fmt.Sprintf("File info: %v", header.Size)))
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
