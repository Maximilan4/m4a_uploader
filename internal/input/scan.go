package input

import (
    "github.com/sirupsen/logrus"
    "io/fs"
    "io/ioutil"
    "m4a_manager/internal/m4a"
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
    var existing *m4a.AudioFile
    var dirItem fs.FileInfo
    for _, dirItem = range elements {
        name = dirItem.Name()
        filePath = path.Join(dir, dirItem.Name())
        if name == "." || name == ".." {
            continue
        }

        if dirItem.IsDir() {
            scan(filePath, outputChan, uploaded)
            continue
        }

        if !strings.Contains(name, ".m4a") {
            continue
        }

        existing = uploaded.SearchByPath(filePath)
        if existing == nil {
            outputChan <- filePath
        } else {
            logrus.Infof("file %s already uploaded", filePath)
        }
    }
}
