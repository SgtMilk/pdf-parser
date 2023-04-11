package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func createRouter(){
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.GET("/ping", func(c *gin.Context) {
		log.Println(c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
		  "message": "pong",
		})
	  })
	router.POST("/parse-pdf", ginParsePDF)
	router.Run(":8080")
}

func ginParsePDF(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	catch(err)
	if file.Filename[len(file.Filename) - 4:] != ".pdf"{
		log.Println("Something else than a pdf file was sent:", file.Filename)
		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": ".pdf file type required",
		  })
		return
	}

	filename := "./source/temp.pdf"
	err = c.SaveUploadedFile(file, filename)
	catch(err)

	defer os.Remove(filename)

	output := ParsePDF(filename)

	c.Data(http.StatusOK, "application/json", output)
}