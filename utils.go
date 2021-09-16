package entityfileuploader

import "net/http"

func GetFileName(r *http.Request, fileName string) (string, error) {
	var err error
	_, handler, err := r.FormFile(fileName)
	if err != nil {
		return "", err
	}
	return handler.Filename, err
}
