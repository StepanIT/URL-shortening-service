package storage

// var UrlMap = make(map[string]string)

var ServerAddress = "localhost:8080"

// интерфейс сохраняет id и url, получает оригинальный url по id
type Repositories interface {
	Save(id string, url string) error
	Get(id string) (string, error)
}
