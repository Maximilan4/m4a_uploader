package source

type AppleTrackInfo struct {
    Amid int64
    Isrc string
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
