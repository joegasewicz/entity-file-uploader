package tests

import (
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"testing"
)

// ImageTestHelper Adds an image form file to the handler
func ImageTestHelper(t *testing.T) (*io.PipeReader, *multipart.Writer) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("catPic", "cat.txt")
		if err != nil {
			t.Error(err)
		}
		catFile, _ := os.Open("./cat.txt")
		defer catFile.Close()
		cat, _ := png.Decode(catFile)
		err = png.Encode(part, cat)
		if err != nil {
			t.Error(err)
		}
	}()
	return pr, writer
}
