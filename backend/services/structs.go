package services

type ImageStruct struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	ETag string `json:"etag"`
	URL  string `json:"url"`
}
