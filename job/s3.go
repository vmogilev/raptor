package job

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// BootAws starts aws session
func BootAws() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

// ToStringP - "cat" s3 key to string and if not found panic
func ToStringP(ctx context.Context, downloader *s3manager.Downloader, log *log.Logger, bucket string, item string) string {
	buff := &aws.WriteAtBuffer{}

	numBytes, err := downloader.DownloadWithContext(ctx, buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		log.Panicf("ERROR: Unable to download item %q, %v\n", item, err)
	}

	return strings.TrimSpace(string(buff.Bytes()[:numBytes]))
}

// ToString - "cat" an s3 object to string and if not found return "" and false
func ToString(ctx context.Context, downloader *s3manager.Downloader, bucket string, item string) (string, bool) {
	buff := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "NoSuchKey" {
				return "", false
			}
		}
		log.Fatalf("Unable to download item %q, %v\n", item, err)
	}

	return strings.TrimSpace(string(buff.Bytes()[:numBytes])), true
}

// CatPath - lists S3 path and returns a string map of keys to their values
func CatPath(ctx context.Context, svc *s3.S3, bucket, path string) (map[string]string, error) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	}

	resp, err := svc.ListObjectsWithContext(ctx, params)
	if err != nil {
		return nil, err
	}

	ret := map[string]string{}
	downloader := s3manager.NewDownloaderWithClient(svc)
	for _, key := range resp.Contents {
		if v, ok := ToString(ctx, downloader, bucket, *key.Key); ok {
			k := filepath.Base(*key.Key)
			ret[k] = v
		}
	}

	return ret, nil
}
