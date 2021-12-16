package input

import (
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "path"
    "strings"
)

func ScanForM4aPaths(dir string) chan string {
    outputChan := make(chan string)
    go func() {
        defer close(outputChan)
        scan(dir, outputChan)
    }()
    return outputChan
}

func scan(dir string, outputChan chan string) {
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
            scan(filePath, outputChan)
            continue
        }

        if !strings.Contains(name, ".m4a") {
            continue
        }

        outputChan <- filePath
    }
}
