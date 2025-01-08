package upload

import (
	"TalkSphere/pkg/oss"
	"TalkSphere/setting"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

var (
	allowedExts = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	MaxFileSize = 10 << 20 // 10MB
)

// SaveImageToOSS 保存图片
func SaveImageToOSS(file *multipart.FileHeader, prefix string, userID int64) (string, error) {
	// 检查文件大小
	if file.Size > int64(MaxFileSize) {
		return "", fmt.Errorf("文件大小超过限制")
	}

	// 检查文件后缀
	ext := strings.ToLower(path.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("不支持的文件类型")
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "upload-*"+ext)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name()) // 使用完后删除临时文件
	defer tmpFile.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 将上传的文件内容复制到临时文件
	if _, err = io.Copy(tmpFile, src); err != nil {
		return "", err
	}

	// 确保文件指针回到开始位置
	if _, err = tmpFile.Seek(0, 0); err != nil {
		return "", err
	}

	// 生成对象键名
	objectKey := fmt.Sprintf("%s/%d_%d%s", prefix, userID, time.Now().UnixNano(), ext)

	// 上传到 COS
	_, err = oss.Client.Object.Put(context.Background(), objectKey, tmpFile, nil)
	if err != nil {
		return "", err
	}

	// 生成访问URL
	url := fmt.Sprintf("https://%s.cos.%s.myqcloud.com/%s",
		setting.Conf.OSSConfig.BucketName,
		setting.Conf.OSSConfig.Region,
		objectKey)

	return url, nil
}
