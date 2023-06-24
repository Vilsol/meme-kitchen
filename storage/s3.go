package storage

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	BaseURL  string `json:"base_url"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
}

type S3 struct {
	BaseURL   string
	S3Client  *s3.S3
	S3Session *session.Session
	Config    Config
}

func ConnecToS3(config Config) *S3 {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(config.Key, config.Secret, ""),
		Endpoint:         aws.String(config.Endpoint),
		Region:           aws.String(config.Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)

	if err != nil {
		log.Err(err).Msg("failed to create S3 session")
		return nil
	}

	s3Client := s3.New(newSession)

	return &S3{
		BaseURL:   config.BaseURL,
		S3Client:  s3Client,
		S3Session: newSession,
		Config:    config,
	}
}

func (s3o *S3) Get(key string) (io.ReadCloser, int64, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	object, err := s3o.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3o.Config.Bucket),
		Key:    aws.String(cleanedKey),
	})

	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get object")
	}

	return object.Body, *object.ContentLength, nil
}

func (s3o *S3) Put(ctx context.Context, key string, body io.Reader) (string, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	uploader := s3manager.NewUploader(s3o.S3Session)

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Key:    aws.String(cleanedKey),
	})

	if err != nil {
		return cleanedKey, errors.Wrap(err, "failed to upload file")
	}

	return key, nil
}

func (s3o *S3) Rename(from string, to string) error {
	cleanedKey := strings.TrimPrefix(to, "/")

	_, err := s3o.S3Client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(viper.GetString("storage.bucket")),
		CopySource: aws.String(viper.GetString("storage.bucket") + from),
		Key:        aws.String(cleanedKey),
	})

	return errors.Wrap(err, "failed to copy object")
}

func (s3o *S3) Delete(key string) error {
	cleanedKey := strings.TrimPrefix(key, "/")

	for i := 0; i < 10; i++ {
		versions, err := s3o.S3Client.ListObjectVersions(&s3.ListObjectVersionsInput{
			Bucket:    aws.String(viper.GetString("storage.bucket")),
			KeyMarker: aws.String(cleanedKey),
			Prefix:    aws.String(cleanedKey),
		})

		if err != nil {
			return errors.Wrap(err, "failed to list object versions")
		}

		objects := make([]*s3.ObjectIdentifier, len(versions.Versions)+len(versions.DeleteMarkers))

		for i, version := range versions.Versions {
			objects[i] = &s3.ObjectIdentifier{
				Key:       version.Key,
				VersionId: version.VersionId,
			}
		}

		for i, marker := range versions.DeleteMarkers {
			objects[i+len(versions.Versions)] = &s3.ObjectIdentifier{
				Key:       marker.Key,
				VersionId: marker.VersionId,
			}
		}

		if len(objects) == 0 {
			return nil
		}

		_, err = s3o.S3Client.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(viper.GetString("storage.bucket")),
			Delete: &s3.Delete{
				Objects: objects,
			},
		})

		if err != nil {
			return errors.Wrap(err, "failed to delete objects")
		}
	}

	return nil
}

type ObjectMeta struct {
	ContentLength *int64
	ContentType   *string
}

func (s3o *S3) Meta(key string) (*ObjectMeta, error) {
	cleanedKey := strings.TrimPrefix(key, "/")

	data, err := s3o.S3Client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(viper.GetString("storage.bucket")),
		Key:    aws.String(cleanedKey),
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get object meta")
	}

	return &ObjectMeta{
		ContentLength: data.ContentLength,
		ContentType:   data.ContentType,
	}, nil
}
