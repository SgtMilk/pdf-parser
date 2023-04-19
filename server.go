package pdfparser

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const MaxMultipartMemory int64 = 8 << 20 // 8 MiB

func CreateRouter() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = MaxMultipartMemory
	router.GET("/ping", func(c *gin.Context) {
		log.Println(c.ClientIP())
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.POST("/parse-pdf", parsePDF)

	catch(router.Run(":8080"))
}

func parsePDF(c *gin.Context) {
	file, err := c.FormFile("file")
	catch(err)

	if file.Filename[len(file.Filename)-4:] != ".pdf" {
		log.Println("Something else than a pdf file was sent:", file.Filename)

		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"message": ".pdf file type required",
		})

		return
	}

	output := ParsePdf(file)

	c.Data(http.StatusOK, "application/json", output)
}
