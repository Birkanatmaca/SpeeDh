package transcription

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

func (h *TranscribeHandler) Transcribe(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı kimliği bulunamadı"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı kimliği formatı geçersiz"})
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ses dosyası alınamadı: " + err.Error()})
		return
	}
	defer file.Close()

	// Geçici bir dosya oluştur ve yolunu servise gönder
	tempDir := "uploads"
	os.MkdirAll(tempDir, os.ModePerm)
	originalFilePath := filepath.Join(tempDir, fmt.Sprintf("%d-%d-%s", userID, time.Now().UnixNano(), header.Filename))
	if err := c.SaveUploadedFile(header, originalFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi: " + err.Error()})
		return
	}

	// --- TEK ÇAĞRI: Tüm iş akışını başlatmak için servisi çağır ---
	transcription, err := h.service.CreateTranscription(userID, originalFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "İşlem başarıyla tamamlandı.",
		"transcription": transcription.TranscriptionText,
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
func (h *TranscribeHandler) GetTranscripts(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı ID'si bulunamadı"})
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı ID'si formatı hatalı"})
		return
	}

	transcripts, err := h.service.GetTranscriptsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transkriptler alınırken bir hata oluştu: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, transcripts)
}
func (h *TranscribeHandler) GetAudioFile(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	transcriptIDStr := c.Param("id")
	transcriptID, err := strconv.ParseUint(transcriptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID formatı"})
		return
	}

	transcript, err := h.service.GetTranscriptByID(uint(transcriptID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transkript bulunamadı"})
		return
	}

	// Güvenlik kontrolü: Kullanıcı sadece kendi ses dosyasını indirebilir.
	if transcript.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu dosyaya erişim yetkiniz yok"})
		return
	}

	// Dosyayı istemciye gönder.
	c.File(transcript.AudioFilepath)
}

func (h *TranscribeHandler) DeleteTranscript(c *gin.Context) {
	userIDValue, _ := c.Get("userID")
	userID := userIDValue.(uint)

	transcriptIDStr := c.Param("id")
	transcriptID, err := strconv.ParseUint(transcriptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID formatı"})
		return
	}

	err = h.service.DeleteTranscript(uint(transcriptID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transkript başarıyla silindi."})
}
