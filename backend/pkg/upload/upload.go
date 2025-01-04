package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// 允许上传的图片后缀
var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// SaveAvatar 保存头像文件
func SaveAvatar(file *multipart.FileHeader, userID int64) (string, error) {
	// 检查文件后缀
	ext := strings.ToLower(path.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("不支持的文件类型")
	}

	// 创建上传目录
	uploadDir := "static/uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
	filepath := path.Join(uploadDir, filename)

	// 打开文件以进行写入
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 将文件内容复制到目标文件
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/" + filepath, nil
}
