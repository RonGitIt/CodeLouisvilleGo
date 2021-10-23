package main

import (
	"awsuploader/server"
	"log"
	"net/http"
)

const (
	bucketName = "go-inventory-images"
)

func main() {
	server := server.NewServer("go-inventory-images")
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello world"))
	} )
	http.HandleFunc("/upload", server.HandleImageUpload)
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal(err)
	}
}