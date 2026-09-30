package requests

type UploadFileRequest struct {
	File string `json:"file" validate:"required,notblank" example:"https://example.com/pc.png"`
}
