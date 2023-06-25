package entityfileuploader

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	fileUpload := FileUpload{
		UploadDir:   "uploads",
		MaxFileSize: 10,
		FileTypes:   []string{"png", "jpeg", "txt"},
	}

	_, err := fileUpload.Init("cats")

	if err != nil {
		t.Errorf("Init should not return an error: %v", err)
	}
	// UploadDir
	expectedUploadDir, _ := os.Getwd()
	expectedUploadDir += "/uploads"
	if fileUpload.UploadDir != "uploads" {
		t.Errorf("Expected  fileUpload.UploadDir to equal '%v' but got '%v'", "uploads", fileUpload.UploadDir)
	}
	if fileUpload.UploadDir != "uploads" {
		t.Errorf("Expected UploadDir to equal '%v' but got '%v'", "uploads", fileUpload.UploadDir)
	}
	// MaxFileSize
	if fileUpload.MaxFileSize != 10 {
		t.Errorf("MaxFileSize should equal 10 not %v", fileUpload.MaxFileSize)
	}
	// FileTypes
	expectedFileTypes := []string{"png", "jpeg"}
	if expectedFileTypes[0] != fileUpload.FileTypes[0] {
		t.Errorf("FileType[0] should equal '%v' but got '%v'", expectedFileTypes[0], fileUpload.FileTypes[0])
	}
	if expectedFileTypes[1] != fileUpload.FileTypes[1] {
		t.Errorf("FileType[1] should equal '%v' but got '%v'", expectedFileTypes[1], fileUpload.FileTypes[1])
	}
	tearDown()
}
