package comic

///跨包使用的话 属性名页要大写
type Comic struct {
	Title   string `json:"comicName"`
	ComicID string `json:"comicID"`
	Cover   string `json:"cover"`
}

//漫画剧集
type Volume struct {
	Tittle string `json:"title"`
	VolID  string `json:"volID"`
}

type VolumeDetail struct {
	Images string `json:"image"`
}
