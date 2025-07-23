package services

type ImageStruct struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	ETag string `json:"etag"`
	URL  string `json:"url"`
}

type RdsImageStruct struct {
	Filename    string `json:"filename"`
	Description string `json:"description"`
	S3Key       string `json:"s3_key"`
}
