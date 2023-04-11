package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createRouter(){
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/parse-pdf", ginParsePDF)
	router.Run(":8080")
}

func ginParsePDF(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	filename := "./temp.pdf"
	c.SaveUploadedFile(file, filename)

	output := ParsePDF(filename)

	c.JSON(http.StatusOK, output)
}