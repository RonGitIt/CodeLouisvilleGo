package main

import (
	"awsuploader/server"
	"fmt"
	"log"
	"net/http"
)

const (
	BUCKETNAME = "go-inventory-images"
	ID         = "7JTD2xL8fQE56RrJ5H4kwUpa0+PqEYZk1PsmOh7WBhf5zD0o2b0idJPuN1Nof2Rc"
	SECRET     = "9ZfHw94LP6P4jnXBMRCUhKFpj+5Z82x3vOajVaecsZ4PTFXcM1o5XGxonLAS4dT+GJajSwUv1zGw82LXdMR6IqiyTio="
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
