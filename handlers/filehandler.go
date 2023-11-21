package handlers

import (
	"gincounter/cmd"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"time"
	"bytes"
	"io"
)

func FileData(c *gin.Context) {
	startTime := time.Now()
	var Routines int
	Routines, err := strconv.Atoi(c.Query("Routines"))
	if err != nil {
		fmt.Println("Invalid Go Routines")
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["file"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	file, err := files[0].Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	fileData := buf.String()
	totalLines, totalWords, totalVowels, totalPunctuations := cmd.FileReader(fileData, Routines)
	totalTime:=time.Since(startTime)
	totalTimeString := totalTime.String()
	c.JSON(http.StatusOK, gin.H{
		"Go Routines":        Routines,
		"Execution Time":     totalTimeString,
		"Total Words ":       totalWords,
		"Total Lines":        totalLines,
		"Total Vowels":       totalVowels,
		"Total Punctuations": totalPunctuations,
	})
}