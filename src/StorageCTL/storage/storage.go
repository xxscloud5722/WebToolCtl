package storage

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

type CloudStorage interface {
	Put(localPath, cloudPath string) (*string, error)
}

type TencentCosClient struct {
	client  *cos.Client
	baseURL *cos.BaseURL
}

func NewTencentCOS(cosURL, secretId, secretKey string) (CloudStorage, error) {
	u, err := url.Parse(cosURL)
	if err != nil {
		return nil, err
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	var storage = TencentCosClient{client: client, baseURL: b}
	return &storage, nil
}
