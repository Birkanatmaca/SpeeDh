// Dosya: internal/transcription/transcribe.go (GÜNCELLENMİŞ VERSİYON)

package transcription

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func runWhisperLocally(filePath string) (string, error) {
	fmt.Printf(">>> Whisper Working: %s\n", filePath)

	baseName := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	txtFilePath := baseName + ".txt"

	cmd := exec.Command("whisper", filePath, "--model", "base", "--language", "Turkish", "--output_format", "txt", "--output_dir", filepath.Dir(filePath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Whisper not working: %v\nOutput: %s", err, string(output))
	}

	content, err := ioutil.ReadFile(txtFilePath)
	if err != nil {
		return "", fmt.Errorf("transkript file don't read: %v", err)
	}

	os.Remove(txtFilePath)

	fmt.Println(">>> Whisper process is success.")
	return string(content), nil
}
func TranscribeHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı, bu rota korumalı olmalı"})
		return
	}

	// userID'nin doğru tipte olduğunu doğrulayalım.
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı ID'si formatı hatalı"})
		return
	}

	fmt.Printf(">>> Transkript isteği, Kullanıcı ID %d tarafından yapıldı.\n", userIDUint)
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ses dosyası bulunamadı: " + err.Error()})
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), ext)

	tempDir := "uploads"
	os.MkdirAll(tempDir, os.ModePerm)

	filePath := filepath.Join(tempDir, newFileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi: " + err.Error()})
		return
	}
	defer os.Remove(filePath)

	transcribedText, err := runWhisperLocally(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Metne çevirme hatası: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transcription": strings.TrimSpace(transcribedText),
	})
}
