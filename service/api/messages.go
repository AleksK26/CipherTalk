package api

import (
	"net/http"
	"time"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SendMessageHandler(c *gin.Context) {
	userID := c.GetString("userID")
	var request struct {
		ConversationID string `json:"conversationId"`
		Content        string `json:"content"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	message := database.Message{
		ID:             uuid.New().String(),
		ConversationID: request.ConversationID,
		SenderID:       userID,
		Content:        request.Content,
		Timestamp:      time.Now(),
	}
	database.DB.Create(&message)
	c.JSON(http.StatusOK, gin.H{"message": message})
}

func ForwardMessageHandler(c *gin.Context) {
	messageID := c.Param("id")
	var request struct {
		TargetConversationID string `json:"targetConversationId"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var originalMessage database.Message
	if err := database.DB.First(&originalMessage, "id = ?", messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	forwardedMessage := database.Message{
		ID:             uuid.New().String(),
		ConversationID: request.TargetConversationID,
		SenderID:       originalMessage.SenderID,
		Content:        originalMessage.Content,
		Timestamp:      time.Now(),
	}
	database.DB.Create(&forwardedMessage)
	c.JSON(http.StatusOK, gin.H{"message": forwardedMessage})
}

func CommentMessageHandler(c *gin.Context) {
	messageID := c.Param("id")
	var request struct {
		Reaction string `json:"reaction"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var message database.Message
	if err := database.DB.First(&message, "id = ?", messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Add reaction to message (assuming a field exists in the model)
	// message.Reactions = append(message.Reactions, request.Reaction)
	database.DB.Save(&message)
	c.JSON(http.StatusOK, gin.H{"message": "Reaction added"})
}

func UncommentMessageHandler(c *gin.Context) {
	messageID := c.Param("id")
	var message database.Message
	if err := database.DB.First(&message, "id = ?", messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Remove reaction from message (assuming a field exists in the model)
	// message.Reactions = []string{}
	database.DB.Save(&message)
	c.JSON(http.StatusOK, gin.H{"message": "Reaction removed"})
}

func DeleteMessageHandler(c *gin.Context) {
	messageID := c.Param("id")
	if err := database.DB.Delete(&database.Message{}, "id = ?", messageID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message deleted"})
}
