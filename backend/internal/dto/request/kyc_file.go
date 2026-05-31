package request

type PresignKycFileReq struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type PresignKycFileAccessReq struct {
	ObjectKey string `json:"object_key" binding:"required"`
}
