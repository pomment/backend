package dao

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

func ReadJSON(p string) (content string, err error) {
	absPath := filepath.Join(BasePath, p)

	jsonFile, err := os.Open(absPath)
	defer jsonFile.Close()

	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func WriteJSON(p string, content string) (err error) {
	absPath := filepath.Join(BasePath, p)

	var jsonFile *os.File
	defer jsonFile.Close()

	jsonFile, err = os.OpenFile(absPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	n, _ := jsonFile.Seek(0, io.SeekEnd)
	_, err = jsonFile.WriteAt([]byte(content), n)

	return nil
}

func GetThreadPath(id string) string {
	return path.Join("threads", fmt.Sprintf("%s.json", id))
}

func GetThreadMetaPath(id string) string {
	return path.Join("threads", fmt.Sprintf("%s.meta.json", id))
}
