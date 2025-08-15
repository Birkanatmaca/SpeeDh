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

// TranscribeHandler, servis bağımlılığını tutar.
type TranscribeHandler struct {
	service ITranscriptionService
}

// NewTranscribeHandler, yeni bir handler örneği oluşturur.
func NewTranscribeHandler(service ITranscriptionService) *TranscribeHandler {
	return &TranscribeHandler{service: service}
}

// Transcribe metodu, ses dosyasını işler, metne çevirir ve veritabanına kaydeder.
func (h *TranscribeHandler) Transcribe(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı, bu rota korumalı olmalı"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı ID'si formatı hatalı"})
		return
	}

	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ses dosyası bulunamadı: " + err.Error()})
		return
	}

	// Dosyayı geçici bir konuma kaydet
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d-%d%s", userID, time.Now().UnixNano(), ext)
	tempDir := "uploads"
	os.MkdirAll(tempDir, os.ModePerm)
	filePath := filepath.Join(tempDir, newFileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi: " + err.Error()})
		return
	}
	defer os.Remove(filePath) // İşlem bitince dosyayı sil

	// Whisper ile metne dönüştür
	transcribedText, err := runWhisperLocally(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Metne çevirme hatası: " + err.Error()})
		return
	}

	// --- YENİ EKLENEN VERİTABANINA KAYIT KISMI ---
	_, err = h.service.SaveTranscription(userID, strings.TrimSpace(transcribedText), file.Filename, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transkript veritabanına kaydedilemedi: " + err.Error()})
		return
	}
	// --- BİTİŞ ---

	c.JSON(http.StatusOK, gin.H{
		"message":       "Dosya başarıyla metne dönüştürüldü ve kaydedildi.",
		"transcription": strings.TrimSpace(transcribedText),
	})
}

// runWhisperLocally fonksiyonu (değişiklik yok)
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
