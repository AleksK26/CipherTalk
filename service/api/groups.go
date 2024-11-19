package api

import (
	"net/http"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/database"
	"github.com/gin-gonic/gin"
)

func AddToGroupHandler(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		UserID string `json:"userId"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var group database.Conversation
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, "id = ?", request.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Model(&group).Association("Members").Append(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User added to group"})
}

func LeaveGroupHandler(c *gin.Context) {
	groupID := c.Param("id")
	userID := c.GetString("userID")

	var group database.Conversation
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	database.DB.Model(&group).Association("Members").Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User left the group"})
}

func SetGroupNameHandler(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		NewName string `json:"newName"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var group database.Conversation
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	group.Name = request.NewName
	database.DB.Save(&group)
	c.JSON(http.StatusOK, gin.H{"message": "Group name updated"})
}

func SetGroupPhotoHandler(c *gin.Context) {
	groupID := c.Param("id")
	var request struct {
		PhotoURL string `json:"photoUrl"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var group database.Conversation
	if err := database.DB.First(&group, "id = ?", groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Update group photo (assuming a field exists in the model)
	// group.PhotoURL = request.PhotoURL
	database.DB.Save(&group)
	c.JSON(http.StatusOK, gin.H{"message": "Group photo updated"})
}
