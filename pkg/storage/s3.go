package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sigit14ap/go-commerce/internal/config"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"
)

type UploadInput struct {
	File string
}

type StorageProvider interface {
	Upload(typeFile string, file *multipart.FileHeader) string
}

type Provider struct {
	cfg *config.Config
}

func NewStorageProvider(cfg *config.Config) *Provider {
	return &Provider{
		cfg: cfg,
	}
}

func (p *Provider) Upload(typeFile string, file *multipart.FileHeader) string {

	s3Config := &aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	}

	s3Session, _ := session.NewSession(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	f, err := file.Open()
	if err != nil {
		panic(err)
	}

	extension := strings.Split(file.Filename, ".")

	name := generateName(20)
	ext := extension[1]

	filename := name + "." + ext

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(typeFile + "/" + filename),
		Body:   f,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("file uploaded")

	return result.Location
}

func generateName(length int) string {
	var output strings.Builder

	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
