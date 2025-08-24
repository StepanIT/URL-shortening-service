package storage

// интерфейс сохраняет id и url, получает оригинальный url по id
type URLShortenerRepositories interface {
	Save(id string, url string) error
	Get(id string) (string, error)
}
