package apple

import (
    "encoding/json"
    "fmt"
    "github.com/sirupsen/logrus"
    "m4a_manager/internal/m4a"
    "m4a_manager/internal/source"
    "net/http"
    "net/url"
    "strconv"
)

type Search struct {
    Results struct {
        Songs struct {
            Data []struct {
                Id         string `json:"id"`
                Attributes struct {
                    Isrc string `json:"isrc"`
                } `json:"attributes"`
            } `json:"data"`
        } `json:"songs"`
    } `json:"results"`
}

func MatchAudioFiles(files chan *m4a.AudioFile, dataset *source.AppleSource) chan *m4a.AudioFile {
    matched := make(chan *m4a.AudioFile)
    go func(searchedTracks chan *m4a.AudioFile) {
        defer close(searchedTracks)
        for file := range files {
            founded, err := match(file, dataset)
            if err != nil {
                logrus.WithError(err).Warningf("Error while match track %s", file.SearchTitle)
                continue
            }
            logrus.Infof("founded track %s", file.SearchTitle)
            searchedTracks <- founded
        }

    }(matched)

    return matched
}

func match(file *m4a.AudioFile, dataset *source.AppleSource) (*m4a.AudioFile, error) {
    requestUrl, _ := url.Parse("https://amp-api.music.apple.com/v1/catalog/ru/search")
    query := url.Values{}
    query.Add("term", file.SearchTitle)
    query.Add("l", "ru")
    query.Add("platform", "web")
    query.Add("limit", "25")
    query.Add("types", "songs")
    query.Add("fields[songs]", "id,isrc")
    requestUrl.RawQuery = query.Encode()

    request, err := http.NewRequest("get", requestUrl.String(), nil)
    if err != nil {
        return nil, err
    }

    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IldlYlBsYXlLaWQifQ.eyJpc3MiOiJBTVBXZWJQbGF5IiwiaWF0IjoxNjM5MDg1NjY1LCJleHAiOjE2NTQ2Mzc2NjV9.hjX-hCq2xVAgjOyiYvnlT6vbhBWl-RZETjVRZ6fGiVHaPW_yjsHv_jOJs57mrt-7uNa8kODd1Eo8dc179YkYoQ"))

    response, err := http.DefaultClient.Do(request)
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()
    var parsedResult Search
    decoder := json.NewDecoder(response.Body)
    err = decoder.Decode(&parsedResult)
    if err != nil {
        return nil, err
    }

    foundedCount := len(parsedResult.Results.Songs.Data)
    if foundedCount == 0 {
        return nil, fmt.Errorf("unable to find any tracks in apple music by title %s", file.SearchTitle)
    }

    var id int64
    for _, song := range parsedResult.Results.Songs.Data {
        id, err = strconv.ParseInt(song.Id, 10, 64)
        if err != nil {
            logrus.WithError(err).Warningf("wrong id for track %s", file.SearchTitle)
            continue
        }

        result := dataset.Search(id)
        if result == nil {
            continue
        }

        if result.Isrc != song.Attributes.Isrc {
            logrus.WithError(err).Warningf("isrc mismatch for track %s", file.SearchTitle)
            continue
        }

        file.Amid = result.Amid
        file.Isrc = result.Isrc

        return file, nil
    }

    return nil, fmt.Errorf("unable to find any tracks in apple music by title %s", file.SearchTitle)
}
