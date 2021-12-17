package main

import (
    "errors"
    "flag"
    "github.com/sirupsen/logrus"
    "m4a_manager/internal/apple"
    "m4a_manager/internal/input"
    "m4a_manager/internal/m4a"
    "m4a_manager/internal/source"
    "m4a_manager/internal/state"
    "m4a_manager/internal/upload"
    "os"
)

var sourceFile *string
var scanDir *string
var awsConfig *string

func init() {
    sourceFile = flag.String("s", "source.csv", "-s <path_to_source_file>")
    scanDir = flag.String("d", "", "-d <path_to_dir_with_m4a>")
    awsConfig = flag.String("a", "aws_config.json", "-a <path_to_aws_config>")
}

func main() {
    flag.Parse()
    if *sourceFile == "" || *scanDir == "" || *awsConfig == "" {
        flag.Usage()
        return
    }

    awsCfg, err := upload.ParseConfig(*awsConfig)
    if err != nil {
        logrus.Fatal(err)
    }

    err = upload.Init(awsCfg)
    if err != nil {
        logrus.Fatal(err)
    }

    dataset, err := initDataSource(*sourceFile)
    if err != nil {
        logrus.Fatal(err)
    }

    uploadedDataset, err := initUploadedDataSource("uploaded.csv")
    if err != nil {
        logrus.Fatal(err)
    }

    m4aPaths := input.ScanForM4aPaths(*scanDir, uploadedDataset)
    audioFiles := m4a.ParseFiles(m4aPaths)
    matchedFiles := apple.MatchAudioFiles(audioFiles, dataset)
    uploadedFiles := upload.M4a(matchedFiles, awsCfg.Bucket, uploadedDataset)
    err = state.SaveUploaded(uploadedFiles, "uploaded.csv")
    if err != nil {
        logrus.Fatal(err)
    }
}

func initDataSource(sourceFile string) (*source.AppleSource, error) {
    reader, err := input.OpenSourceFile(sourceFile)
    if err != nil {
        logrus.Fatal(err)
    }

    dataSource := source.NewAppleSource()
    err = dataSource.LoadFromCsv(reader)
    if err != nil {
        logrus.Fatal(err)
    }

    return dataSource, nil
}

func initUploadedDataSource(sourceFile string) (*source.UploadedM4aSource, error) {
    dataSource := source.NewUploadedM4aSource()
    if _, err := os.Stat(sourceFile); errors.Is(err, os.ErrNotExist) {
        return dataSource, nil
    } else if err != nil {
        return nil, err
    }

    reader, err := input.OpenSourceFile(sourceFile)
    if err != nil {
        logrus.Fatal(err)
    }

    err = dataSource.LoadFromCsv(reader)
    if err != nil {
        logrus.Fatal(err)
    }

    return dataSource, nil
}
