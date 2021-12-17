package input

import (
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "m4a_manager/internal/source"
    "path"
    "strings"
)

func ScanForM4aPaths(dir string, uploaded *source.UploadedM4aSource) chan string {
    outputChan := make(chan string)
    go func() {
        defer close(outputChan)
        scan(dir, outputChan, uploaded)
    }()
    return outputChan
}

func scan(dir string, outputChan chan string, uploaded *source.UploadedM4aSource) {
    elements, err := ioutil.ReadDir(dir)
    if err != nil {
        logrus.Error(err)
        return
    }

    var name string
    var filePath string
    for _, f := range elements {
        name = f.Name()
        filePath = path.Join(dir, f.Name())
        if name == "." || name == ".." {
            continue
        }

        if f.IsDir() {
            scan(filePath, outputChan, uploaded)
            continue
        }

        if !strings.Contains(name, ".m4a") {
            continue
        }

        existing := uploaded.SearchByPath(filePath)
        if existing == nil {
            outputChan <- filePath
        } else {
            logrus.Infof("file %s already uploaded", filePath)
        }
    }
}
