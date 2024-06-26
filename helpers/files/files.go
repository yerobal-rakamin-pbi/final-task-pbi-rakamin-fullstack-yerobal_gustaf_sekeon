package files

import (
	"mime/multipart"
	"os"
	"strings"
)

func IsExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

type File struct {
	Content multipart.File
	Meta    *multipart.FileHeader
}

func Init(content multipart.File, meta *multipart.FileHeader, key string) (*File, error) {
	file := &File{
		Content: content,
		Meta:    meta,
	}

	return file, nil
}

func (f *File) SetFileName(newName string) {
	fileName := strings.Split(f.Meta.Filename, ".")
	fileName[0] = newName
	f.Meta.Filename = strings.Join(fileName, ".")
}

func (f *File) IsImage() bool {
	return strings.Contains(f.Meta.Header.Get("Content-Type"), "image")
}

func GetFileNameFromURL(url string) string {
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]

	return fileName
}
