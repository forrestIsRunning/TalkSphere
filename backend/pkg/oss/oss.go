package oss

import (
	"fmt"
	"github.com/TalkSphere/backend/setting"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var Client *cos.Client

func Init(cfg *setting.OSSConfig) error {
	// 初始化 COS 客户端
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", cfg.BucketName, cfg.Region))
	if err != nil {
		return err
	}

	b := &cos.BaseURL{BucketURL: u}
	Client = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})

	return nil
}

// GetObjectURL 获取对象的访问URL
func GetObjectURL(objectKey string) string {
	return fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		setting.Conf.OSSConfig.BucketName,
		setting.Conf.OSSConfig.Region,
		objectKey)
}
