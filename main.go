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

	config := server.AwsConfig{
		Bucket: bucketName,
		Password: passphrase,
		Id: "7JTD2xL8fQE56RrJ5H4kwUpa0+PqEYZk1PsmOh7WBhf5zD0o2b0idJPuN1Nof2Rc",
		Secret: "9ZfHw94LP6P4jnXBMRCUhKFpj+5Z82x3vOajVaecsZ4PTFXcM1o5XGxonLAS4dT+GJajSwUv1zGw82LXdMR6IqiyTio=",
	}
	server := server.NewServer(config)
	http.HandleFunc("/upload", server.HandleImageUpload)
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal(err)
	}
}