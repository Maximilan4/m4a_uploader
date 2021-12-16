package input

import (
    "encoding/csv"
    "os"
)

func OpenSourceFile(path string) (*csv.Reader, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }

    return csv.NewReader(file), nil
}
