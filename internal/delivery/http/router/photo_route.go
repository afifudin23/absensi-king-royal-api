package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func registerPhotoRoutes(rg *gin.RouterGroup) {
	rg.POST("/uploads", func(c *gin.Context) {

		name := c.PostForm("name")
		userId := c.PostForm("user_id")

		file, err := c.FormFile("photo")
		if err != nil {
			c.JSON(400, gin.H{"error": "foto wajib"})
			return
		}

		c.SaveUploadedFile(file, "./uploads/"+file.Filename)

		c.JSON(200, gin.H{
			"user": name,
			"id":   userId,
		})
	})

	rg.DELETE("/uploads", func(c *gin.Context) {
		var body map[string]string

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"error": "invalid payload",
			})
			return
		}

		filePath := body["path"]
		if filePath == "" {
			c.JSON(400, gin.H{
				"error": "path required",
			})
			return
		}

		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete file",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "file deleted",
		})
	})
}
