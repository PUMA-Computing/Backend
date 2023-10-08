package services

import (
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type FilesService struct {
}

func NewFilesService() *FilesService {
	return &FilesService{}
}

func (ns *FilesService) GetFileNameWithoutExtension(fileName string) string {
	return path.Base(fileName[:len(fileName)-len(path.Ext(fileName))])
}

func (ns *FilesService) ConvertFileName(fileName string) string {
	fileExtension := filepath.Ext(fileName)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	return timestamp + "-" + ns.GetFileNameWithoutExtension(fileName) + fileExtension
}

func (ns *FilesService) ConvertToPublicDirectory(fileName string) string {
	return filepath.Join("public/", filepath.Base(fileName))
}
