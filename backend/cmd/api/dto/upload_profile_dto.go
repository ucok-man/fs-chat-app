package dto

import (
	"mime/multipart"
)

type ReqUploadProfileDto struct {
	File     *multipart.File
	Metadata *multipart.FileHeader
}
