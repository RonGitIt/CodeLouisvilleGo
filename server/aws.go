package server

import (
	"errors"
	"fmt"
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
	Id string
	Secret string
}

type AwsConfig struct {
	Bucket string
	Password string
	Id string
	Secret string
}

func NewAws(config AwsConfig) (*AWS, error) {
	awsS3 := &AWS {
		Bucket: config.Bucket,
		Id: "",
		Secret: "",
	}

	//id, err := Decrypt("7JTD2xL8fQE56RrJ5H4kwUpa0+PqEYZk1PsmOh7WBhf5zD0o2b0idJPuN1Nof2Rc", passphrase)
	id, err := Decrypt(config.Id, config.Password)
	if err != nil {
		return awsS3, fmt.Errorf("error decrypting id during AWS struct creation: %w", err)
	}
	awsS3.Id = id

	//secret, err := Decrypt("9ZfHw94LP6P4jnXBMRCUhKFpj+5Z82x3vOajVaecsZ4PTFXcM1o5XGxonLAS4dT+GJajSwUv1zGw82LXdMR6IqiyTio=", passphrase)
	secret, err := Decrypt(config.Secret, config.Password)
	if err != nil {
		return awsS3, fmt.Errorf("error decrypting secret during AWS struct creation: %w", err)
	}
	awsS3.Secret = secret

	err = awsS3.setupSession()
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
	// Check if file already exists. Don't overwrite existing files
	if exists, err := a.AlreadyExists(filename); err != nil {
		return "", err
	} else if exists {
		return "", errors.New("Duplicate filename already exists in bucket")
	}
	uploader := s3manager.NewUploader(a.Session)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:  aws.String(a.Bucket),
		Key: aws.String(filename),
		Body: file,
	})
	if err != nil {
		return "", fmt.Errorf("error during S3 upload: %w", err)
	}
	return result.Location, nil
}

func (a *AWS) AlreadyExists(objectName string) (bool, error) {
	s3Service := s3.New(a.Session)
	resp, err := s3Service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(a.Bucket),
	})
	if err != nil {
		return true, fmt.Errorf("could not set up S3 connection: %w", err)
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

func (a *AWS) DeleteFile(objectName string) (bool, error) {
		s3Service := s3.New(a.Session)
		_, err := s3Service.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(a.Bucket),
			Key: aws.String(objectName),
		})
		if err != nil {
			log.Printf("Error deleting object (%s) from bucket (%s): %s", objectName, a.Bucket, err)
			return false, err
		}
		err = s3Service.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(a.Bucket),
			Key: aws.String(objectName),
		})
		if err != nil {
			log.Printf("Error waiting for object deleteion (%s:%s): %s", a.Bucket, objectName, err)
			return false, err
		}
		return true, nil
}

func (a *AWS) GetFile(objectName string) (*[]byte, int64, error) {
	downloader := s3manager.NewDownloader(a.Session)
	fileData := aws.NewWriteAtBuffer([]byte{})
	n, err := downloader.Download(fileData,
		&s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key: aws.String(objectName),
		})
	if err != nil {
		return nil, 0, fmt.Errorf("error while downloading in GetFile: %w", err)
	}
	fileBytes := fileData.Bytes()
	return &fileBytes, n, nil
}