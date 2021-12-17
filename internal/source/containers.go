package source

import "m4a_manager/internal/m4a"

type AppleTrackInfo struct {
    Amid int64
    Isrc string
}

type AudioFilesByAmid []*m4a.AudioFile

func (a AudioFilesByAmid) Len() int {
    return len(a)
}

func (a AudioFilesByAmid) Less(i, j int) bool {
    return a[i].Amid < a[j].Amid
}

func (a AudioFilesByAmid) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

type AudioFilesByPath []*m4a.AudioFile

func (a AudioFilesByPath) Len() int {
    return len(a)
}

func (a AudioFilesByPath) Less(i, j int) bool {
    return a[i].Path < a[j].Path
}

func (a AudioFilesByPath) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

type TracksInfo struct {
    data []AppleTrackInfo
}

func (t *TracksInfo) Push(info AppleTrackInfo) {
    t.data = append(t.data, info)
}

func (t *TracksInfo) Len() int {
    return len(t.data)
}

func (t *TracksInfo) Less(i, j int) bool {
    return t.data[i].Amid < t.data[j].Amid
}

func (t *TracksInfo) Swap(i, j int) {
    t.data[i], t.data[j] = t.data[j], t.data[i]
}
