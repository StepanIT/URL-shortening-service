package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	db *sql.DB
}

func NewPingHandler(db *sql.DB) *PingHandler {
	return &PingHandler{db: db}
}

func (h *PingHandler) Ping(c *gin.Context) {
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			log.Println("Database ping failed:", err)
			c.String(http.StatusInternalServerError, "Database connection failed")
			return
		}
	}
	c.String(http.StatusOK, "OK")
}
