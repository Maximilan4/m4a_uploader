package upload

import (
    "context"
    "errors"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/sirupsen/logrus"
    "m4a_manager/internal/m4a"
    "os"
    "time"
)

var uploader *manager.Uploader

func Init(cfg *Config) error {
    awsConfig, err := config.LoadDefaultConfig(context.TODO(),
        config.WithEndpointResolverWithOptions(getResolverFunc(cfg)),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.SecretAccessKey, "")),
    )

    if err != nil {
        return err
    }

    client := s3.NewFromConfig(awsConfig)
    uploader = manager.NewUploader(client)
    return nil
}

func getResolverFunc(cfg *Config) aws.EndpointResolverWithOptionsFunc {
    return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        var endpoint aws.Endpoint
        if cfg.Region == "" || cfg.Host == "" {
            return endpoint, errors.New("specify aws_host env variable")
        }

        endpoint.PartitionID = "aws"
        endpoint.URL = cfg.Host
        endpoint.SigningRegion = cfg.Region
        endpoint.HostnameImmutable = false
        endpoint.Source = aws.EndpointSourceCustom

        return endpoint, nil
    }
}

func M4a(files chan *m4a.AudioFile, bucketName string) {
    var err error
    for file := range files {
        err = uploadM4aFile(file, bucketName)
        if err != nil {
            logrus.WithError(err).Warningf("Uploading file %s is failed", file.Path)
            continue
        }

        logrus.Infof("upload is complete for file %s", file.Path)
    }
}

func uploadM4aFile(file *m4a.AudioFile, bucketName string) error {
    ctx, done := context.WithTimeout(context.Background(), time.Second*15)
    defer done()
    key := fmt.Sprintf("%s/%s", file.Isrc, "track.m4a")

    reader, err := os.Open(file.Path)
    if err != nil {
        return err
    }

    _, err = uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket: &bucketName,
        Key:    &key,
        Body:   reader,
    })
    if err != nil {
        return err
    }

    return nil
}
