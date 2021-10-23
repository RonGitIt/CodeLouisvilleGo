package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

type AWS struct {
	Session *session.Session
	Bucket string
}

func NewAWS(bucket string) (*AWS, error){
	aws := AWS{
		Bucket: bucket,
	}
	err := aws.SetupSession()
	if err != nil {
		return nil, err
	}
	return &aws, nil
}

func (a *AWS) SetupSession() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN"),
	})
	if err != nil {
		log.Printf("Error setting up session: %v\n", err)
		return err
	}
	a.Session = sess
	return nil
}

func (a *AWS) UploadFile(filename string, file io.Reader) error {
	uploader := s3manager.NewUploader(a.Session)
	uploader.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(a.Bucket),
		Key: aws.String(filename),
		Body: file,
	})
	return nil
}


