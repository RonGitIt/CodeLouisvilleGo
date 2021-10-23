package server

import (
	"log"
)

type Server struct {
	awsS3 *AWS
}

func NewServer(bucket string) *Server {
	s3, err := NewAws(bucket)
	if err != nil {
		log.Fatalf("Error setting up S3: %s", err)
	}

	return &Server{
		awsS3: s3,
	}
}