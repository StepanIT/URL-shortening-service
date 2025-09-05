package storage

import (
	"errors"
)

// структура в которой map для хранения ключа ID и значения URL
type InMemoryStorage struct {
	data     map[string]string
	userData map[string]string
}

// конструктор для создания нового экземпляра InMemoryStorage, с новой map
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data:     make(map[string]string),
		userData: make(map[string]string),
	}
}

// функция сохраняет оригинальный URL по заданному короткому ID
func (s *InMemoryStorage) Save(id, url, userID string) error {
	s.data[id] = url
	s.userData[id] = userID
	return nil
}

// функция извлекает оригинальный URL по короткому ID
func (s *InMemoryStorage) Get(id string) (string, error) {
	url, ok := s.data[id]
	if !ok {
		return "", errors.New("not found")
	}
	return url, nil
}

func (s *InMemoryStorage) GetAllByUser(userID string) (map[string]string, error) {
	result := make(map[string]string)

	// ищем все shortID, которые принадлежат пользователю
	for shortID, uID := range s.userData {
		if uID == userID {
			// если нашли, добавляем в результат пару shortID -> originalURL
			result[shortID] = s.data[shortID]
		}
	}

	return result, nil
}
