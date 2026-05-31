package dto

// FileUploadResp is the data payload from POST /rest/v2.0/upload.
// See: https://opendocs.kun.global/docs/api/file-upload
type FileUploadResp struct {
	URL string `json:"url"`
}
