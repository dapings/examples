package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix not supported")
	}
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit")
	}

	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create the file save dir")
		}
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions")
	}

	dst := filepath.Join(uploadSavePath, fileName)
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}
	accessUrl := filepath.Join(global.AppSetting.UploadServerUrl, fileName)
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
