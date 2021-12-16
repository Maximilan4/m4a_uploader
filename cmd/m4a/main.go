package main

import (
    "flag"
    "github.com/sirupsen/logrus"
    "m4a_manager/internal/apple"
    "m4a_manager/internal/input"
    "m4a_manager/internal/m4a"
    "m4a_manager/internal/source"
    "m4a_manager/internal/upload"
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

    m4aPaths := input.ScanForM4aPaths(*scanDir)
    audioFiles := m4a.ParseFiles(m4aPaths)
    matchedFiles := apple.MatchAudioFiles(audioFiles, dataset)
    upload.M4a(matchedFiles, awsCfg.Bucket)
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
