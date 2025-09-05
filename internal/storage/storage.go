package storage

// интерфейс сохраняет id и url, получает оригинальный url по id
type URLShortenerRepositories interface {
	Save(id string, url string, userID string) error
	Get(id string) (string, error)
	GetAllByUser(userID string) (map[string]string, error)
}
