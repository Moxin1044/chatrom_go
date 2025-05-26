package main

import (
	"chatrom/controllers"
	"chatrom/models"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	// 创建上传文件夹（如果不存在）
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// 初始化数据库（创建表）
	models.InitDB()

	r := gin.Default()

	// 静态文件托管，供前端访问上传的图片和文件
	r.Static("/uploads", "./uploads")
	r.Static("/static", "./static") // 静态资源路径

	// 渲染前端页面（你可以替换成你自己的前端html）
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 发送消息接口
	r.POST("/send_message", controllers.SendMessage)

	// 获取消息接口
	r.GET("/get_messages", controllers.GetMessages)

	r.Run(":8080")
}
