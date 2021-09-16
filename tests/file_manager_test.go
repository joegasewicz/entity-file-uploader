package tests

import (
	entityfileuploader "github.com/joegasewicz/entity-file-uploader"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var TEST_ROOT_PATH, _ = os.Getwd()
var TEST_DIR_PATH = TEST_ROOT_PATH + "/uploads/cats/1"
var TEST_FILE_PATH = TEST_ROOT_PATH + "/uploads/cats/1/cat.txt"

func setUp() {
	// Copy the local catpic.png to the entity destination directory
	os.MkdirAll(TEST_DIR_PATH, os.ModePerm)
	file, err := os.Open("cat.txt")
	defer file.Close()
	if err != nil {
		println(err)
	}
	fileDest, err := os.Create(TEST_FILE_PATH)
	defer fileDest.Close()
	if err != nil {
		println(err)
	}
	io.Copy(fileDest, file)
}

func tearDown() {
	// Remove uploads and child dirs etc.
	filePath, _ := os.Getwd()
	os.RemoveAll(filePath + "/uploads")
}

func TestUpload(t *testing.T) {
	fileUpload := entityfileuploader.FileUpload{
		UploadDir:   "uploads",
		MaxFileSize: 10,
		FileTypes:   []string{"png", "jpeg", "txt"},
	}

	catUpload, _ := fileUpload.Init("cats")

	catHandler := func(w http.ResponseWriter, r *http.Request) {
		filePath, _ := catUpload.Upload(w, r, 1, "catPic")
		if filePath != TEST_FILE_PATH {
			t.Errorf("Expected '%v' but got '%v'", TEST_FILE_PATH, filePath)
		}
	}

	// Add an image form file to the request
	pr, writer := ImageTestHelper(t)
	req := httptest.NewRequest("POST", "http://localhost:8000/", pr)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	catHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode to be 200 but got '%v'", resp.StatusCode)
	}

	// Test the file is there
	file, err := os.Open(TEST_FILE_PATH)

	if err != nil {
		t.Error(err)
	}

	if file.Name() != TEST_FILE_PATH {
		t.Errorf("Expected 'cat.txt' but got '%v'", file.Name())
	}

	tearDown()
}

func TestGet(t *testing.T) {
	setUp()
	fileUpload := entityfileuploader.FileUpload{
		UploadDir:   "uploads",
		MaxFileSize: 10,
		FileTypes:   []string{"png", "jpeg", "txt"},
		URL:         "http://localhost:8080",
	}

	catUpload, _ := fileUpload.Init("cats")

	catHandler := func(w http.ResponseWriter, r *http.Request) {
		filePath := catUpload.Get("cat.txt", 1)
		expected := "http://localhost:8080/uploads/cats/1/cat.txt"
		if filePath != expected {
			t.Errorf("Expected '%v' but got '%v'", expected, filePath)
		}
	}
	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()
	catHandler(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("Expected StatusCode to be 200 but got '%v'", resp.StatusCode)
	}
	tearDown()
}

func TestUpdate(t *testing.T) {
	setUp()
	fileUpload := entityfileuploader.FileUpload{
		UploadDir:   "uploads",
		MaxFileSize: 10,
		FileTypes:   []string{"png", "jpeg"},
	}

	catUpload, _ := fileUpload.Init("cats")
	err := catUpload.Update("cat.txt", 1, "tomcat.png")
	if err != nil {
		t.Error(err)
	}
	file, _ := os.Open(TEST_DIR_PATH + "/tomcat.png")
	defer file.Close()
	expected := TEST_DIR_PATH + "/tomcat.png"
	if file.Name() != expected {
		t.Errorf("Expected '%v' but got '%v'", expected, file.Name())
	}
	tearDown()
}

func TestDelete(t *testing.T) {
	setUp()
	fileUpload := entityfileuploader.FileUpload{
		UploadDir:   "uploads",
		MaxFileSize: 10,
		FileTypes:   []string{"png", "jpeg"},
	}

	catUpload, _ := fileUpload.Init("cats")
	err := catUpload.Delete("cat.txt", 1)
	if err != nil {
		t.Error(err)
	}
	_, err = os.Stat(TEST_FILE_PATH)
	if os.IsNotExist(err) != true {
		t.Error("Expected file to not be in path")
	}
	tearDown()
}
