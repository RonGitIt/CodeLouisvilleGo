package server

import (
	"log"
)

type Server struct {
	awsS3 *AWS
}

func NewServer(bucket string, passphrase string) *Server {
	s3, err := NewAws(bucket)
	if err != nil {
		log.Fatalf("Error setting up S3: %s", err)
	}

	return &Server{
		awsS3: s3,
	}
}

func (s *Server) TesthelperDuplicateCheck(objectName string) (bool, error) {
	return s.awsS3.AlreadyExists(objectName)
}

func(s *Server) TesthelperDeleteFile(objectName string) (bool, error) {
	return s.awsS3.DeleteFile(objectName)
}