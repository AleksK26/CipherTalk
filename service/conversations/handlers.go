package conversations

import (
	"net/http"

	"github.com/AleksK26/wasatext/service/db"
	"github.com/gin-gonic/gin"
)

func GetConversationsHandler(c *gin.Context) {
	userID := c.GetString("userID")
	var conversations []db.Conversation
	db.DB.Preload("Members").Where("EXISTS (SELECT 1 FROM conversation_members WHERE user_id = ?)", userID).Find(&conversations)
	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func GetConversationHandler(c *gin.Context) {
	conversationID := c.Param("id")
	var conversation db.Conversation
	if err := db.DB.Preload("Members").First(&conversation, "id = ?", conversationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"conversation": conversation})
}
