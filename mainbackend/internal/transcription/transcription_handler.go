package transcription

import (
	"fmt"
	"io/ioutil"
	"log"
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
	// --- MEVCUT KODUNUZ (Değişiklik yok) ---
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

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d-%d%s", userID, time.Now().UnixNano(), ext)
	tempDir := "uploads"
	os.MkdirAll(tempDir, os.ModePerm)
	filePath := filepath.Join(tempDir, newFileName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Dosya kaydedilemedi: " + err.Error()})
		return
	}
	// Orijinal dosyayı işlem bittiğinde silmek için defer ekliyoruz.
	defer os.Remove(filePath)
	log.Printf("Orijinal dosya geçici olarak kaydedildi: %s", filePath)
	// --- MEVCUT KODUNUZ BİTTİ ---

	// --- YENİ EKLENEN FFMPEG İLE SES İŞLEME ADIMI ---
	cleanedFilePath := strings.TrimSuffix(filePath, ext) + "_cleaned.wav"
	// Temizlenmiş dosyayı da işlem bittiğinde sil.
	defer os.Remove(cleanedFilePath)

	log.Printf("FFmpeg ile gürültü temizleme başlatılıyor: %s -> %s", filePath, cleanedFilePath)
	cmd := exec.Command("ffmpeg", "-i", filePath, "-af", "afftdn", "-ar", "16000", "-ac", "1", cleanedFilePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("FFmpeg hatası: %s", string(output))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ses dosyası işlenemedi (ffmpeg error)"})
		return
	}
	log.Println("Gürültü temizleme tamamlandı.")
	// --- YENİ ADIM BİTTİ ---

	// --- DEĞİŞTİRİLEN KISIM ---
	// Whisper'a orijinal dosya yerine TEMİZLENMİŞ dosyayı veriyoruz.
	log.Printf("Whisper temizlenmiş dosya ile başlatılıyor: %s", cleanedFilePath)
	transcribedText, err := runWhisperLocally(cleanedFilePath) // <--- Değişiklik burada
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Metne çevirme hatası: " + err.Error()})
		return
	}
	// --- DEĞİŞİKLİK BİTTİ ---

	// --- MEVCUT KODUNUZ (Değişiklik yok) ---
	// Veritabanına orijinal dosyanın adını ve yolunu kaydediyoruz.
	_, err = h.service.SaveTranscription(userID, strings.TrimSpace(transcribedText), file.Filename, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transkript veritabanına kaydedilemedi: " + err.Error()})
		return
	}

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
