package media

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/logger"
)

type Media struct {
	client *cloudinary.Cloudinary
}

func New(url string) (*Media, error) {
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	cld.Logger.SetLevel(logger.ERROR)
	cld.Config.URL.Secure = true

	return &Media{
		client: cld,
	}, nil
}

type uploadopt func(p *uploader.UploadParams)

func UploadWithFolder(folder string) uploadopt {
	return func(p *uploader.UploadParams) {
		p.AssetFolder = folder
		p.UseAssetFolderAsPublicIDPrefix = api.Bool(true)
	}
}

func UploadWithReplaceable(id string) uploadopt {
	return func(p *uploader.UploadParams) {
		p.PublicID = id
		p.UniqueFilename = api.Bool(false)
		p.Overwrite = api.Bool(true)
	}
}

func (m *Media) Upload(file io.Reader, opts ...uploadopt) (string, error) {
	param := uploader.UploadParams{}
	for _, optfn := range opts {
		optfn(&param)
	}

	resp, err := m.client.Upload.Upload(context.Background(), file, param)
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}
