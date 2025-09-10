package storage

import (
	"os"
	"testing"
)

func TestPostgresStorage_SaveAndGet(t *testing.T) {
	// Читаем DSN из переменной окружения
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		t.Skip("DATABASE_DSN not set — skipping Postgres test")
	}

	// Создаём хранилище
	store, err := NewPostgresStorage(dsn)
	if err != nil {
		t.Fatalf("failed to create PostgresStorage: %v", err)
	}

	userID := "test-user"
	shortURL := "abc123"
	originalURL := "https://example.com"

	// Сохраняем
	if err := store.Save(shortURL, originalURL, userID); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Получаем по shortURL
	got, err := store.Get(shortURL)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}
	if got != originalURL {
		t.Errorf("expected %s, got %s", originalURL, got)
	}

	// Получаем все ссылки пользователя
	all, err := store.GetAllByUser(userID)
	if err != nil {
		t.Fatalf("failed to get all by user: %v", err)
	}

	if url, ok := all[shortURL]; !ok || url != originalURL {
		t.Errorf("expected in map: %s -> %s, got: %v", shortURL, originalURL, all)
	}
}
