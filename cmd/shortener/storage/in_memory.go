package storage

import (
	"errors"
)

// структура в которой map для хранения ключа ID и значения URL
type InMemoryStorage struct {
	Data map[string]string
}

// конструктор для создания нового экземпляра InMemoryStorage, с новой map
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Data: make(map[string]string),
	}
}

// функция сохраняет оригинальный URL по заданному короткому ID
func (s *InMemoryStorage) Save(id, url string) error {
	s.Data[id] = url
	return nil
}

// функция извлекает оригинальный URL по короткому ID
func (s *InMemoryStorage) Get(id string) (string, error) {
	url, ok := s.Data[id]
	if !ok {
		return "", errors.New("not found")
	}
	return url, nil
}
