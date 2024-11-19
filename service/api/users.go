package api

import (
	"net/http"

	"github.com/AleksK26/WASA_AleksK_2024-25/service/database"
	"github.com/gin-gonic/gin"
)

func SetMyPhotoHandler(c *gin.Context) {
	userID := c.GetString("userID")
	var request struct {
		PhotoURL string `json:"photoUrl"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.PhotoURL = request.PhotoURL
	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Photo updated"})
}

func LoginHandler(c *gin.Context) {
	var request struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user database.User
	result := database.DB.FirstOrCreate(&user, database.User{Name: request.Name})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": user.ID})
}

func SetMyUserName(c *gin.Context) {
	userID := c.GetString("userID")
	var request struct {
		NewName string `json:"newName"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user database.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Name = request.NewName
	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Name updated"})
}
