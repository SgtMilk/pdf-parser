package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func createRouter(){
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.GET("/ping", func(c *gin.Context) {
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
	// log.Println(file.Filename)

	// Upload the file to specific dst.
	filename := "./source/temp.pdf"
	err = c.SaveUploadedFile(file, filename)
	catch(err)

	defer os.Remove(filename)

	output := ParsePDF(filename)

	c.Data(http.StatusOK, "application/json", output)
}