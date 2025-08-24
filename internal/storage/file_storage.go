package storage

import (
	"encoding/json"
	"errors"
	"log"
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
	log.Printf("Сохранение в файл %s: %s -> %s", fs.filePath, id, url)
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
	// convert map to JSON
	data, err := json.Marshal(fs.data)
	if err != nil {
		return err
	}

	// write the JSON to the file
	return os.WriteFile(fs.filePath, data, 0644)
}

// load reads data from a file and fills the map
func (fs *FileStorage) load() error {
	// Read the entire file
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err // This will include os.ErrNotExist
	}

	// If the file is empty, do nothing
	if len(data) == 0 {
		return nil
	}

	// Unmarshal JSON into the map
	return json.Unmarshal(data, &fs.data)
}
