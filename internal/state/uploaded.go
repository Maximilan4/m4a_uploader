package state

import (
    "encoding/csv"
    "github.com/sirupsen/logrus"
    "m4a_manager/internal/m4a"
    "os"
    "strconv"
)

func SaveUploaded(files chan *m4a.AudioFile, outputFile string) error {
    output, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    defer output.Close()
    writer := csv.NewWriter(output)
    for file := range files {
        err = writer.Write([]string{strconv.FormatInt(file.Amid, 10), file.Isrc, file.Path})
        if err != nil {
            logrus.WithError(err).Error("cant write output to uploaded state file")
        }
        writer.Flush()
    }

    return nil
}
