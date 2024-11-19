package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/api"
	"github.com/AleksK26/WASA_AleksK_2024-25/service/database"

	"github.com/gin-gonic/gin"
)

func main() {
	loadConfiguration() // Call loadConfiguration
	database.ConnectDatabase()

	router := gin.Default()
	router.Use(enableCORS()) // Call enableCORS
	router.Use(api.Middleware())

	registerAPIRoutes(router)
	registerWebUI(router) // Call registerWebUI

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	gracefulShutdown(server)
}

func registerAPIRoutes(router *gin.Engine) {
	router.POST("/session", api.LoginHandler)
	router.PATCH("/users/me", api.SetMyUserName)
	router.GET("/conversations", api.GetConversationsHandler)
	router.GET("/conversations/:id", api.GetConversationHandler)
	router.POST("/messages", api.SendMessageHandler)
	router.POST("/messages/:id/forward", api.ForwardMessageHandler)
	router.POST("/messages/:id/comments", api.CommentMessageHandler)
	router.DELETE("/messages/:id/comments", api.UncommentMessageHandler)
	router.DELETE("/messages/:id", api.DeleteMessageHandler)
	router.POST("/groups/:id/members", api.AddToGroupHandler)
	router.DELETE("/groups/:id/members", api.LeaveGroupHandler)
	router.PATCH("/groups/:id/name", api.SetGroupNameHandler)
	router.PATCH("/users/me/photo", api.SetMyPhotoHandler)
	router.PATCH("/groups/:id/photo", api.SetGroupPhotoHandler)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shut down: %s\n", err)
	}

	log.Println("Server shut down successfully")
}
