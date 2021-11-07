package main

import (
	"awsuploader/server"
	"fmt"
	"log"
	"net/http"
)

const (
	BUCKETNAME = "go-inventory-images"
	ID         = "OZOXcDQ7SEuTUqrKSRM3AJFHIBudZGSChkdGiTmDlIowjjj2uMk7IMfEeLf6+pV+"
	SECRET     = "ASlR6S3dVtqJslXc4ruSQBt6ndA3DAJVrXDEsbj2x0dOM3k5oo7JDdfTOE0uBlkgHosK2RYojuRDueu41aGKBgKAojY="
)

func main() {
	var passphrase string
	fmt.Print("Enter password to unlock AWS services: ")
	fmt.Scanln(&passphrase)

	// Set up AWS connection details
	config := server.AwsConfig{
		Bucket:   BUCKETNAME,
		Password: passphrase,
		Id:       ID,
		Secret:   SECRET,
	}
	server, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Error while setting up server: %v", err)
	}
	log.Println("password accepted... starting up the server")

	// Set up routes and turn fire it up
	http.HandleFunc("/upload", server.HandleImageUpload)
	http.HandleFunc("/get/", server.HandleGetImage)
	if err = http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal(err)
	}
}
