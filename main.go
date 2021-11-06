package main

import (
	"awsuploader/server"
	"fmt"
	"log"
	"net/http"
)

const (
	bucketName = "go-inventory-images"
)

func main() {
	var passphrase string
	fmt.Println("Enter password to unlock AWS services: ")
	fmt.Scanln(&passphrase)

	server := server.NewServer("go-inventory-images", passphrase)
	http.HandleFunc("/upload", server.HandleImageUpload)
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal(err)
	}
}