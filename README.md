[![Go](https://github.com/joegasewicz/entity-file-uploader/actions/workflows/go.yml/badge.svg)](https://github.com/joegasewicz/entity-file-uploader/actions/workflows/go.yml)
[![GitHub license](https://img.shields.io/github/license/joegasewicz/entity-file-uploader)](https://github.com/joegasewicz/entity-file-uploader/blob/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/joegasewicz/entity-file-uploader)](https://github.com/joegasewicz/entity-file-uploader/issues)
# Entity File Uploader
Handles file uploads &amp; organises files based on your database entities.
Read the [docs](https://pkg.go.dev/github.com/joegasewicz/entity-file-uploader).

### Install
```bash
go get -u github.com/joegasewicz/entity-file-uploader
```
## Examples

### Create a new File Manager Entity
```go
import (
	entityfileuploader "github.com/joegasewicz/entity-file-uploader"
)

var FileUpload = entityfileuploader.FileUpload{
	UploadDir:   "uploads",
	MaxFileSize: 10,
	FileTypes:   []string{"png", "jpeg"},
	URL: "http://localhost:8080",
}
```

### Create a FileManager
A FileManager is specific to your database table name
```go
catUpload, err := FileUpload.Init("cats")
```

### Upload a file 
We can now use our FileManager to save a file to `/uploads/users/1/cats/catpic.png`
```go
// Cat.ID is the primary key value & Cat is a Gorm / ORM struct
// and we are uploading a file with a name of `catpic.jpg`
func Post(w http.ResponseWriter, r *http.Request) {
    // GetFileName util function takes the form name of your file upload as the 2nd arg
    avatarFileName, _ :=  entityfileuploader.GetFileName(r, "catImage")
    // Example uses Gorm >>
    Cat := models.Cat{
        Avatar: avatarFileName,
    }
    result := DB.Create(&Cat)
    // Gorm <<
	// Upload method takes the UUID of your saved & created entity
    fileName, err := CatUpload.Upload(w, r, Cat.ID, Cat.Avatar)
    if err == nil {
    // Handle error
    }
    fmt.Printf("Saved new file to: %s\n", fileName)	// /uploads/users/1/cats/catpic.png
}
```

### Fetch the file's filepath
Gets the full file path including the filename
```go
fileName := CatUpload.Get(Cat.Avatar, Cat.ID)
fmt.Println(fileName) // http://localhost:8080/uploads/cats/1/catpic.png
```

### Update a file 
Updates only the filename not the file (Use `Upload` to update the file)
```go
err := CatUpload.Update(Cat.Avatar, Cat.ID, "tomcat.png")
```

### Delete a file
// Deletes the file from the entity file path
```go
err := CatUpload.Delete(Cat.Avatar, Cat.ID)
```

### Handle file uploads over http for multipart formdata
```go
err := fileManager.ReceiveMultiPartFormDataAndSaveToDir(r, "logo", fileModel.ID)
```