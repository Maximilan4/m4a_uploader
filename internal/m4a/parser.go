package m4a

import (
    "fmt"
    "github.com/frolovo22/tag"
    "github.com/sirupsen/logrus"
)

type AudioFile struct {
    SearchTitle string
    Isrc        string
    Path        string
    Amid        int64
}

func ParseFiles(pathChan chan string) chan *AudioFile {
    output := make(chan *AudioFile)

    go func() {
        defer close(output)
        var err error
        for path := range pathChan {
            err = parseFile(path, output)
            if err != nil {
                logrus.WithError(err).Errorf("unable to parse %s", path)
            }
        }
    }()

    return output
}

func parseFile(path string, output chan *AudioFile) error {
    var err error
    var artist, title string
    var meta tag.Metadata

    meta, err = tag.ReadFile(path)
    if err != nil {
        return err
    }

    artist, err = meta.GetArtist()
    if err != nil {
        return err
    }

    title, err = meta.GetTitle()
    if err != nil {
        return err
    }

    output <- &AudioFile{
        SearchTitle: fmt.Sprintf("%s - %s", artist, title),
        Isrc:        "",
        Path:        path,
    }

    return nil
}
