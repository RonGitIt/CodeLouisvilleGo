package server

import (
	"fmt"
)

type Server struct {
	awsS3 *AWS
}

func NewServer(config AwsConfig) (*Server, error) {
	s3, err := NewAws(config)
	if err != nil {
		return nil, fmt.Errorf("Error setting up S3: %w", err)
	}

	return &Server{ awsS3: s3 }, nil
}

func (s *Server) TesthelperDuplicateCheck(objectName string) (bool, error) {
	return s.awsS3.AlreadyExists(objectName)
}

func(s *Server) TesthelperDeleteFile(objectName string) (bool, error) {
	return s.awsS3.DeleteFile(objectName)
}