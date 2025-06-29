package entityfileuploader

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type FileManager struct {
	*FileUpload
	TableName string
}

// Set the full entity directory path if it doesn't already exist.
// The full entity directory path is constructed with the following directories if we have a table / entity called `users`:
// upload path + entity name + plus the primary id
// e.g. `./assets/uploads/users/1`
func (f *FileManager) setEntityDirPath(id string) (string, error) {
	var err error
	fullDirPath := f.GetEntityDirPath(id)
	err = os.MkdirAll(fullDirPath, os.ModePerm)
	return fullDirPath, err
}

// GetEntityDirPath Gets the full entity file path only (does not include the filename)
func (f *FileManager) GetEntityDirPath(id string) string {
	fullDirPath := fmt.Sprintf("%s/%s/%s", f.fullDirPath, f.TableName, id)
	return fullDirPath
}

// GetEntityURL GetEntityUrl gets the full url
func (f *FileManager) GetEntityURL(fileName string, id string) string {
	url := fmt.Sprintf("%s/%s/%s/%s/%s", f.URL, f.UploadDir, f.TableName, id, fileName)
	return url
}

// GetEntityFilePath Get the full path which is constructed of the following:
// upload path + entity name + plus the primary id + the filename
// e.g. `./assets/uploads/users/1/cats.jpg`
func (f *FileManager) GetEntityFilePath(fullEntityDirPath string, fileName string) string {
	return fmt.Sprintf("%s/%s", fullEntityDirPath, fileName)
}

// Upload Uploads & stores the file on your server. If there is an error then the return value of string
// will be an empty string. If there are no errors to return the string will be the full path of
// with filename.
func (f *FileManager) Upload(w http.ResponseWriter, r *http.Request, id string, formName string) (string, error) {
	var err error
	r.ParseMultipartForm(f.MaxFileSize << 20) // 50mb
	file, handler, err := r.FormFile(formName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileName := handler.Filename
	// Create file
	// Default behavior is to overwrite file if filename is the same as incomming.
	fullPath, err := f.setEntityDirPath(id)
	if err != nil {
		return "", err
	}
	fullEntityFilePath := f.GetEntityFilePath(fullPath, fileName)
	dest, err := os.Create(fullEntityFilePath)
	defer dest.Close()
	if err != nil {
		return "", err
	}
	// Copy the uploaded file to the created file on the file system
	if _, err := io.Copy(dest, file); err != nil {
		return "", err
	}
	return fullEntityFilePath, err
}

// Get Gets the full file path including the filename
func (f *FileManager) Get(fileName string, id string) string {
	entityFilePath := f.GetEntityURL(fileName, id)
	return entityFilePath
}

// Update Updates only the filename not the file (Use Upload to update the file)
func (f *FileManager) Update(fileName string, id string, newFileName string) error {
	var err error
	entityDirPath := f.GetEntityDirPath(id)
	entityFilePath := f.GetEntityFilePath(entityDirPath, fileName)
	entityNewFilePath := f.GetEntityFilePath(entityDirPath, newFileName)
	err = os.Rename(entityFilePath, entityNewFilePath)
	return err
}

// Delete Deletes the file from the entity file path
func (f *FileManager) Delete(fileName string, id string) error {
	entityDirPath := f.GetEntityDirPath(id)
	entityFilePath := f.GetEntityFilePath(entityDirPath, fileName)
	err := os.Remove(entityFilePath)
	return err
}

// DeleteEntityByID Deletes a single entity by a UUID
func (f *FileManager) DeleteEntityByID(id string) error {
	entityDirPath := f.GetEntityDirPath(id)
	err := os.RemoveAll(entityDirPath)
	return err
}

// ReceiveMultiPartFormDataAndSaveToDir Handles a file uploaded via multipart form request over http
// The field argument is the form field's name
func (f *FileManager) ReceiveMultiPartFormDataAndSaveToDir(r *http.Request, field string, id string) error {
	r.ParseMultipartForm(50 << 20) // 50mb
	file, header, err := r.FormFile(field)

	fullPath, err := f.setEntityDirPath(id)
	fullEntityFilePath := f.GetEntityFilePath(fullPath, header.Filename)

	if err != nil {
		return err
	}
	defer file.Close()

	distFile, err := os.Create(fullEntityFilePath)
	defer distFile.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(distFile, file); err != nil {
		return err
	}
	return nil
}
