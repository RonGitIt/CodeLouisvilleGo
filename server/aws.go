package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
)

type AWS struct {
	Session *session.Session
	Bucket string
}

func NewAws(bucket string) (*AWS, error) {
	awsS3 := &AWS {
		Bucket: bucket,
	}
	err := awsS3.setupSession()
	if err != nil {
		return nil, err
	}
	return awsS3, nil
}

func (a *AWS) setupSession() error {
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

func (a *AWS) AlreadyExists(objectName string) (bool, error) {
	s3Service := s3.New(a.Session)
	resp, err := s3Service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(a.Bucket),
	})
	if err != nil {
		log.Printf("S3 error: Could not set up S3 connection: %s", err)
		return true, err
	}

	// Look through the items in the bucket and see if an object with
	// the same name already exists.
	for _, item := range resp.Contents {
		if *item.Key == objectName {
			// Found it!
			return true, nil
		}
	}

	// Not found
	return false, nil
}


