package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"fmt"
	"flag"
	"net/http"
	"net/url"
	"path/filepath"
	"github.com/gin-gonic/gin"
)

const (
	MaxUploadSize = 16 << 20 // 16 mb
)

func main() {
	//port flag
	server_port := flag.Int("p",8080,"The port the server will run in. 8080 by default.")
	
	router := gin.Default()

	router.Static("/files", "./files")
	router.Static("/static","./static")
	router.LoadHTMLGlob("templates/*")
	router.POST("/upload", func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

		if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
			//"_" represents the response from c.Request.ParseMultipartForm and "ok" is the bool value
			if _, ok := err.(*http.MaxBytesError); ok {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{
					"error": fmt.Sprintf("file too large (max: %d bytes)", MaxUploadSize),
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		defer file.Close()
		
		//generate hash
		hasher := sha256.New()
		if _, err := io.Copy(hasher,file); err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		//encode hash
		hash := hex.EncodeToString(hasher.Sum(nil))
		
		file.Seek(0,0)
		ext := filepath.Ext(fileHeader.Filename)
		filename := hash + ext
		destination := filepath.Join("./files/", filename)
		if _, err := os.Stat(destination); err == nil {
			c.String(http.StatusOK, fmt.Sprintf("http://%s/files/%s", c.Request.Host,url.QueryEscape(filename)))
			return
		}
		c.SaveUploadedFile(fileHeader, destination)
		c.String(http.StatusOK, fmt.Sprintf("http://%s/files/%s", c.Request.Host,url.QueryEscape(filename)))
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	flag.Parse()
	router.Run(fmt.Sprintf(":%d",*server_port))
}
