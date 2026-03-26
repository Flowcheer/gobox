package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
)

const (
	MaxUploadSize = 16 << 20 // 16 mb
)

func main() {
	//port flag
	server_port := flag.Int("p",8080,"The port the server will run in. 8080 by default.")
	server_ip := flag.String("ip", "localhost", "The IP address the server will listen on. localhost by default.")
	
	router := gin.Default()
	
	var validFilename = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	
	router.Static("/static", "./static")
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
	router.GET("/files/:filename", func (c *gin.Context){
		filename := c.Param("filename")
		if !validFilename.MatchString(filename){
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
            return
		}
		if filename == "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not specified"})
			return
		}
		filePath := filepath.Join("./files",filepath.Clean(filename))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.File(filePath)
	})
	flag.Parse()
	router.Run(fmt.Sprintf("%s:%d",*server_ip, *server_port))
}
