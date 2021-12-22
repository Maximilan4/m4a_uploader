package source

import (
    "encoding/csv"
    "github.com/sirupsen/logrus"
    "io"
    "m4a_manager/internal/m4a"
    "sort"
    "strconv"
)

type UploadedM4aSource struct {
    UploadedByAmids AudioFilesByAmid
    UploadedByPath  AudioFilesByPath
}

func (as *UploadedM4aSource) SearchByAmid(amid int64) *m4a.AudioFile {
    start := 0
    end := as.UploadedByAmids.Len() - 1
    var (
        mid      int
        midValue int64
    )

    for start <= end {
        mid = (start + end) / 2
        midValue = as.UploadedByAmids[mid].Amid
        if midValue == amid {
            return &as.UploadedByAmids[mid]
        } else if midValue < amid {
            start = mid + 1
        } else if midValue > amid {
            end = mid - 1
        }
    }

    return nil
}

func (as *UploadedM4aSource) SearchByPath(path string) *m4a.AudioFile {
    start := 0
    end := as.UploadedByPath.Len() - 1
    var (
        mid      int
        midValue string
    )
    for start <= end {
        mid = (start + end) / 2
        midValue = as.UploadedByPath[mid].Path
        if midValue == path {
            return &as.UploadedByAmids[mid]
        } else if midValue < path {
            start = mid + 1
        } else if midValue > path {
            end = mid - 1
        }
    }

    return nil
}

func (as *UploadedM4aSource) Push(file m4a.AudioFile) {
    as.UploadedByPath = append(as.UploadedByPath, file)
    as.UploadedByAmids = append(as.UploadedByAmids, file)
    sort.Sort(as.UploadedByAmids)
    sort.Sort(as.UploadedByPath)
}

func (as *UploadedM4aSource) LoadFromCsv(reader *csv.Reader) error {
    var (
        amid int64
        err  error
        line []string
    )
    var file m4a.AudioFile
    for {
        line, err = reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        amid, err = strconv.ParseInt(line[0], 10, 64)
        if err != nil {
            logrus.Error(err)
            continue
        }

        file = m4a.AudioFile{
            Isrc: line[1],
            Path: line[2],
            Amid: amid,
        }
        as.UploadedByPath = append(as.UploadedByPath, file)
        as.UploadedByAmids = append(as.UploadedByAmids, file)
    }

    sort.Sort(as.UploadedByAmids)
    sort.Sort(as.UploadedByPath)
    return nil
}

func NewUploadedM4aSource() *UploadedM4aSource {
    return &UploadedM4aSource{
        UploadedByAmids: make(AudioFilesByAmid, 0, 32*1024),
        UploadedByPath:  make(AudioFilesByPath, 0, 32*1024),
    }
}
