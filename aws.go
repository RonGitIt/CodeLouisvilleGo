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

func SetupAws(bucket string) (error){
	awsS3 = &AWS{
		Bucket: bucket,
	}
	err := awsS3.SetupSession()
	if err != nil {
		return err
	}
	return nil
}

func (a *AWS) SetupSession() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials("AKIA56L2FSSKHAQU6XET", "0km7B4XbLWak8yBXXBdTi9vRSjB4V6wqKCL7hGGQ", ""),
	})
	if err != nil {
		log.Printf("Error setting up session: %v\n", err)
		return err
	}
	a.Session = sess
	return nil
}

func (a *AWS) UploadFile(filename string, file io.Reader) (string, error) {
	uploader := s3manager.NewUploader(a.Session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(a.Bucket),
		Key: aws.String(filename),
		Body: file,
	})
	if err != nil {
		log.Printf("Error during S3 upload: %s", err)
		return "", err
	}
	return result.Location, nil
}


