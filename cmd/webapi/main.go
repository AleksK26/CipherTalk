package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AleksK26/wasatext/service/auth"
	"github.com/AleksK26/wasatext/service/conversations"
	"github.com/AleksK26/wasatext/service/db"
	"github.com/AleksK26/wasatext/service/groups"
	"github.com/AleksK26/wasatext/service/messages"
	"github.com/AleksK26/wasatext/service/users"

	"github.com/gin-gonic/gin"
)

func main() {
	loadConfiguration()
	db.ConnectDatabase()

	router := gin.Default()
	router.Use(enableCORS())
	router.Use(auth.Middleware())

	registerAPIRoutes(router)
	registerWebUI(router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	gracefulShutdown(server)
}

func registerAPIRoutes(router *gin.Engine) {
	router.POST("/session", users.LoginHandler)
	router.PATCH("/users/me", users.SetMyUserName)
	router.GET("/conversations", conversations.GetConversationsHandler)
	router.GET("/conversations/:id", conversations.GetConversationHandler)
	router.POST("/messages", messages.SendMessageHandler)
	router.POST("/messages/:id/forward", messages.ForwardMessageHandler)
	router.POST("/messages/:id/comments", messages.CommentMessageHandler)
	router.DELETE("/messages/:id/comments", messages.UncommentMessageHandler)
	router.DELETE("/messages/:id", messages.DeleteMessageHandler)
	router.POST("/groups/:id/members", groups.AddToGroupHandler)
	router.DELETE("/groups/:id/members", groups.LeaveGroupHandler)
	router.PATCH("/groups/:id/name", groups.SetGroupNameHandler)
	router.PATCH("/users/me/photo", users.SetMyPhotoHandler)
	router.PATCH("/groups/:id/photo", groups.SetGroupPhotoHandler)
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
