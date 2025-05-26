package controllers

import (
	"chatrom/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var allowedExtensions = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".txt": true, ".pdf": true, ".docx": true,
}

func isAllowedExtension(ext string) bool {
	_, ok := allowedExtensions[ext]
	return ok
}

// 发送消息接口
func SendMessage(c *gin.Context) {
	username := c.PostForm("username")
	messageType := c.PostForm("message_type")
	messageContent := c.PostForm("message")

	file, _ := c.FormFile("file")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	// 处理文件上传
	if file != nil {
		ext := filepath.Ext(file.Filename)
		if !isAllowedExtension(ext) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
			return
		}

		// 用时间戳重命名，防止覆盖
		newFileName := strconv.FormatInt(time.Now().UnixNano(), 10) + ext
		savePath := filepath.Join("uploads", newFileName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
			return
		}

		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			messageType = "image"
		} else {
			messageType = "file"
		}
		messageContent = "/uploads/" + newFileName
	}

	if messageType == "text" && messageContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty text message"})
		return
	}

	msg := models.Message{
		Username:    username,
		Content:     messageContent,
		MessageType: messageType,
		CreatedAt:   time.Now(),
	}

	if err := models.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": messageType + " sent"})
}

// 获取消息接口
func GetMessages(c *gin.Context) {
	var messages []models.Message
	if err := models.DB.Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
