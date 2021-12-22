package source

import (
    "encoding/csv"
    "github.com/sirupsen/logrus"
    "io"
    "sort"
    "strconv"
)

type AppleSource struct {
    Items *TracksInfo
}

func (as *AppleSource) Search(amid int64) *AppleTrackInfo {
    start := 0
    end := as.Items.Len() - 1
    var (
        mid      int
        midValue int64
    )
    for start <= end {
        mid = (start + end) / 2
        midValue = as.Items.data[mid].Amid
        if midValue == amid {
            return &as.Items.data[mid]
        } else if midValue < amid {
            start = mid + 1
        } else if midValue > amid {
            end = mid - 1
        }
    }

    return nil
}

func (as *AppleSource) LoadFromCsv(reader *csv.Reader) error {
    var (
        amid int64
        err  error
        line []string
    )
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

        as.Items.Push(AppleTrackInfo{
            Amid: amid,
            Isrc: line[1],
        })
    }

    sort.Sort(as.Items)
    return nil
}

func NewAppleSource() *AppleSource {
    info := &TracksInfo{
        data: make([]AppleTrackInfo, 0, 32*1024),
    }

    return &AppleSource{info}
}
