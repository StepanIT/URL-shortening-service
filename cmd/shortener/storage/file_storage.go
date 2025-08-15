package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// structure for storing data in a JSON file
type FileStorage struct {
	filePath string
	data     map[string]string
}

// constructor to create a new FileStorage instance
func NewFileStorage(filePath string) (*FileStorage, error) {
	fs := &FileStorage{
		filePath: filePath,
		data:     make(map[string]string),
	}

	// trying to load data from a file
	if err := fs.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return fs, nil
}

// Save saves the URL under the given ID to a file
func (fs *FileStorage) Save(id string, url string) error {
	fs.data[id] = url
	return fs.save()
}

// returns URL by ID or error
func (fs *FileStorage) Get(id string) (string, error) {
	url, ok := fs.data[id]
	if !ok {
		return "", errors.New("url not found")
	}
	return url, nil
}

// save writes all pairs of id and url
func (fs *FileStorage) save() error {
	file, err := os.Create(fs.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for id, url := range fs.data {
		// each entry as a separate JSON object in a row
		line := fmt.Sprintf(`{"id":"%s","url":"%s"}`+"\n", id, url)
		if _, err := file.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

// load reads data from a file and fills the map
func (fs *FileStorage) load() error {
	file, err := os.Open(fs.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fs.data = make(map[string]string)

	var id, url string
	decoder := json.NewDecoder(file)
	for {
		var entry map[string]string
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		id = entry["id"]
		url = entry["url"]
		fs.data[id] = url
	}

	return nil
}
