package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func registerWebUI(router *gin.Engine) {
	webUIPath := filepath.Join("..", "webui", "dist")
	if _, err := os.Stat(webUIPath); os.IsNotExist(err) {
		log.Println("Web UI not found, skipping registration")
		return
	}

	router.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(webUIPath, "index.html"))
	})

	router.StaticFS("/static", http.Dir(webUIPath))
}
