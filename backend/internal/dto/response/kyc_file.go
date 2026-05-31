package response

type PresignKycFileResp struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresIn int    `json:"expires_in"`
}

type PresignKycFileAccessResp struct {
	AccessURL string `json:"access_url"`
	ExpiresIn int    `json:"expires_in"`
}
