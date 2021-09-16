package entityfileuploader

import (
	"fmt"
	"os"
)

type FileUpload struct {
	UploadDir   string
	fullDirPath string
	MaxFileSize uint8
	FileTypes   []string
	URL         string
}

// Init is a factory function that returns a pointer to FileManager.
// tableName is the name of the entity associated with a database table name.
func (f *FileUpload) Init(tableName string) (*FileManager, error) {
	err := f.setUploadDir()
	f.setFileTypes()
	fm := FileManager{
		FileUpload: f,
		TableName:  tableName,
	}
	return &fm, err
}

func (f *FileUpload) setFileTypes() {
	if f.FileTypes == nil {
		f.FileTypes = make([]string, 4, 20)
		f.FileTypes[0] = "png"
		f.FileTypes[1] = "jpg"
		f.FileTypes[2] = "jpeg"
		f.FileTypes[3] = "mpeg"
	}
}

func (f *FileUpload) setUploadDir() error {
	var err error
	apiRootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if f.UploadDir == "" {
		// Set f.UploadDir using default
		f.UploadDir = "uploads"
		f.fullDirPath = fmt.Sprintf("%s/%s", apiRootDir, f.UploadDir)
		fmt.Printf("Warning: FileUpload not set so using default: %s\n", f.UploadDir)
	} else {
		// Set f.fullDirPath folder using UploadDir
		f.fullDirPath = fmt.Sprintf("%s/%s", apiRootDir, f.UploadDir)
	}
	// Sets upload dir if doesnt already exist
	err = os.MkdirAll(f.fullDirPath, os.ModePerm)
	return err
}
