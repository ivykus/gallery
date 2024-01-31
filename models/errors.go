package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email is already in use")
	ErrNotFound   = errors.New("models: no resource could be found with the provided information")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("file error: %s", fe.Issue)
}

func CheckContentType(r io.ReadSeeker, allowedTypes []string) error {
	buf := make([]byte, 512)
	_, err := r.Read(buf)
	if err != nil {
		return fmt.Errorf("check content type: %w", err)
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("check content type: %w", err)
	}

	contentType := http.DetectContentType(buf)
	for _, t := range allowedTypes {
		if contentType == t {
			return nil
		}
	}
	return FileError{fmt.Sprintf("content type %s not allowed", contentType)}
}

func CheckExtension(filename string, allowedExtensions []string) error {
	ext := filepath.Ext(filename)
	for _, e := range allowedExtensions {
		if ext == e {
			return nil
		}
	}
	return FileError{fmt.Sprintf("extension %s not allowed", ext)}
}
