package entityfileuploader

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
		part, err := writer.CreateFormFile("catPic", "cat.png")
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		catFile, err := os.Open("./test_data/cat.png")
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
		defer catFile.Close()
		cat, err := png.Decode(catFile)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}

		err = png.Encode(part, cat)
		if err != nil {
			_ = pw.CloseWithError(err)
			return
		}
	}()
	return pr, writer
}
