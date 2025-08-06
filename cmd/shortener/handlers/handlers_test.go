package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StepanIT/URL-shortening-service/cmd/shortener/storage"
	"github.com/gin-gonic/gin"
)

// creating a Gin context
func getTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	return c
}

// successful POST request
func TestPostHandler_Success(t *testing.T) {
	// create inMemory storage
	repo := storage.NewInMemoryStorage()
	h := &Handler{Repo: repo}

	// create test request
	w := httptest.NewRecorder()
	c := getTestGinContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://youtube.com"))

	// call handler
	h.PostHandler(c)

	// checks
	if w.Code != http.StatusCreated {
		t.Errorf("ожидался статус 201 Created, а получили %d", w.Code)
	}

}

// request body is empty
func TestPostHandler_EmptyBody(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	h := &Handler{Repo: repo}

	w := httptest.NewRecorder()
	c := getTestGinContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	h.PostHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("ожидался статус 400 Bad Request, получили %d", w.Code)
	}
}

// successful GET request
func TestGetHandler_Success(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	id := "abc123"
	url := "https://youtube.com"

	// save id and url
	repo.Save(id, url)
	h := &Handler{Repo: repo}

	w := httptest.NewRecorder()
	c := getTestGinContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/get/"+id, nil)
	c.Params = []gin.Param{{Key: "id", Value: id}}

	h.GetHandler(c)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("ожидался статус 307 Temporary Redirect, получили %d", w.Code)
	}

	// location check
	location := w.Header().Get("Location")
	if location != url {
		t.Errorf("ожидался редирект на %s, а получили %s", url, location)
	}
}

// if the ID does not exist in the database
func TestGetHandler_NotFound(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	h := &Handler{Repo: repo}

	w := httptest.NewRecorder()
	c := getTestGinContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/get/неизвестныйID", nil)

	h.GetHandler(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("ожидался статус 404 Not Found, получили %d", w.Code)
	}
}

// ошибка, если ID не передан
func TestGetHandler_EmptyID(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	h := &Handler{Repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/get/", nil)
	w := httptest.NewRecorder()

	h.GetHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("ожидался статус 400 Bad Request при пустом ID, получили %d", resp.StatusCode)
	}
}

// ошибка, если метод запроса не GET
func TestGetHandler_WrongMethod(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	h := &Handler{Repo: repo}

	req := httptest.NewRequest(http.MethodPost, "/get/someid", nil)
	w := httptest.NewRecorder()

	h.GetHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("ожидался статус 404 при методе POST, получили %d", resp.StatusCode)
	}
}
