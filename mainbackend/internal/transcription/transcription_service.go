package transcription

import (
	"errors"
	"fmt"
	"log"
	"mainbackend/internal/model"
	"mainbackend/internal/platform/ai"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ITranscriptionService arayüzü, transkript ile ilgili iş mantığı operasyonlarını tanımlar.
type ITranscriptionService interface {
	CreateTranscription(userID uint, originalFilePath string) (*model.Transcription, error) // <-- BU SATIRI EKLEYİN
	SaveTranscription(userID uint, text, originalFilename, filepath string) (*model.Transcription, error)
	GetTranscriptsByUserID(userID uint) ([]model.Transcription, error)
	GetTranscriptByID(id uint) (*model.Transcription, error)
	DeleteTranscript(id uint, userID uint) error
}

// transcriptionService, repository'yi kullanarak işlemleri gerçekleştirir.
// transcriptionService struct'ına aiClient alanını ekliyoruz.
type transcriptionService struct {
	repo     ITranscriptionRepository
	aiClient ai.IAIClient
}

// NewTranscriptionService fonksiyonunu güncelliyoruz, artık aiClient parametresi alıyor.
func NewTranscriptionService(repo ITranscriptionRepository, aiClient ai.IAIClient) ITranscriptionService {
	return &transcriptionService{
		repo:     repo,
		aiClient: aiClient,
	}
}

func (s *transcriptionService) GetTranscriptsByUserID(userID uint) ([]model.Transcription, error) {
	return s.repo.FindByUserID(userID)
}
func (s *transcriptionService) GetTranscriptByID(id uint) (*model.Transcription, error) {
	return s.repo.FindByID(id)
}

// SaveTranscription, yeni bir transkript nesnesi oluşturur ve veritabanına kaydeder.
func (s *transcriptionService) SaveTranscription(userID uint, text, originalFilename, filepath string) (*model.Transcription, error) {
	transcription := &model.Transcription{
		Title:             originalFilename, // Başlık olarak orijinal dosya adını kullanabiliriz
		OriginalFilename:  originalFilename,
		AudioFilepath:     filepath,
		TranscriptionText: text,
		UserID:            userID,
	}

	err := s.repo.Create(transcription)
	if err != nil {
		return nil, err
	}

	return transcription, nil
}
func (s *transcriptionService) DeleteTranscript(id uint, userID uint) error {
	// Önce transkripti bulalım ki dosya yolunu alabilelim ve yetki kontrolü yapabilelim.
	transcript, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("silinecek transkript bulunamadı")
	}

	// Güvenlik: Kullanıcı sadece kendi transkriptini silebilir.
	if transcript.UserID != userID {
		return errors.New("bu kaydı silme yetkiniz yok")
	}

	// 1. Ses dosyasını diskten sil.
	if err := os.Remove(transcript.AudioFilepath); err != nil {
		// Hata olursa loglayabiliriz ama işleme devam edebiliriz (belki dosya zaten yoktur).
		fmt.Printf("Uyarı: Ses dosyası silinemedi: %s, Hata: %v\n", transcript.AudioFilepath, err)
	}

	// 2. Veritabanı kaydını sil.
	return s.repo.DeleteByID(id)
}
func (s *transcriptionService) ProcessAudioFile(originalFilePath string) (string, error) {
	// --- YENİ ADIM: GÜRÜLTÜ TEMİZLEME ---

	// 1. Temizlenmiş dosya için yeni bir yol oluştur.
	// Örneğin: "uploads/dosya.mp3" -> "uploads/dosya_cleaned.wav"
	ext := filepath.Ext(originalFilePath)
	cleanedFilePath := strings.TrimSuffix(originalFilePath, ext) + "_cleaned.wav"

	fmt.Printf("FFmpeg ile gürültü temizleme başlatılıyor: %s -> %s\n", originalFilePath, cleanedFilePath)

	// 2. FFmpeg komutunu hazırla.
	// -i: input dosyası
	// -af afftdn: Gürültü azaltma için "Audio Fast Fourier Transform Denoise" filtresi. En basit ve etkili olanlardan biridir.
	// cleanedFilePath: çıktı dosyası
	cmd := exec.Command("ffmpeg", "-i", originalFilePath, "-af", "afftdn", cleanedFilePath)

	// 3. Komutu çalıştır ve hataları yakala.
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Eğer FFmpeg hata verirse, hatayı ve çıktıyı loglayalım.
		fmt.Printf("FFmpeg komutu başarısız oldu. Çıktı: %s\n", string(output))
		return "", fmt.Errorf("ffmpeg error: %w", err)
	}

	fmt.Println("Gürültü temizleme tamamlandı.")

	// Geçici dosyayı daha sonra silmek için defer kullanabiliriz.
	defer func() {
		fmt.Printf("Geçici dosya siliniyor: %s\n", cleanedFilePath)
		os.Remove(cleanedFilePath)
	}()

	// --- MEVCUT WHISPER ADIMI ---

	// 4. Whisper'ı TEMİZLENMİŞ dosya ile çağır.
	// BURASI SİZİN MEVCUT WHISPER ÇAĞIRMA KODUNUZUN GELECEĞİ YER
	// Örnek:
	// transcript, err := callWhisper(cleanedFilePath)
	// if err != nil {
	//     return "", err
	// }
	// return transcript, nil

	// Şimdilik sadece dosya yolunu döndürüyoruz.
	// Bu satırı kendi whisper çağırma mantığınızla değiştirin.
	fmt.Printf("Whisper işlemi '%s' dosyası ile devam edecek.\n", cleanedFilePath)
	// Örnek: return s.whisperClient.Transcribe(cleanedFilePath)

	// Bu kısım sizin projenize göre düzenlenmelidir.
	// Geçici olarak boş bir string ve nil hata döndürüyorum.
	return "Transkripsiyon metni buraya gelecek", nil
}

// CreateTranscription, tam iş akışını yönetir.
func (s *transcriptionService) CreateTranscription(userID uint, originalFilePath string) (*model.Transcription, error) {
	defer cleanupFile(originalFilePath)
	ext := filepath.Ext(originalFilePath)
	cleanedFilePath := strings.TrimSuffix(originalFilePath, ext) + "_cleaned.wav"
	defer cleanupFile(cleanedFilePath)

	// Adım 1: FFmpeg
	log.Printf("Adım 1/4: FFmpeg ile ses işleniyor...")
	cmd := exec.Command("ffmpeg", "-i", originalFilePath, "-af", "afftdn", "-ar", "16000", "-ac", "1", cleanedFilePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("FFmpeg hatası: %s", string(output))
		return nil, fmt.Errorf("ses dosyası işlenemedi")
	}
	log.Println(">>> FFmpeg başarılı.")

	// Adım 2: Whisper
	log.Printf("Adım 2/4: Whisper ile transkripsiyon...")
	transcriptFilePath := strings.TrimSuffix(cleanedFilePath, ".wav") + ".txt"
	defer cleanupFile(transcriptFilePath)
	whisperCmd := exec.Command("whisper", cleanedFilePath, "--model", "base", "--language", "Turkish", "--output_format", "txt", "--output_dir", filepath.Dir(cleanedFilePath))
	if output, err := whisperCmd.CombinedOutput(); err != nil {
		log.Printf("Whisper hatası: %s", string(output))
		return nil, fmt.Errorf("transkripsiyon işlemi başarısız oldu")
	}
	log.Println(">>> Whisper başarılı.")

	contentBytes, err := os.ReadFile(transcriptFilePath)
	if err != nil {
		return nil, fmt.Errorf("transkript dosyası okunamadı: %w", err)
	}
	rawTranscriptText := string(contentBytes)

	// Adım 3: AI Düzeltme
	log.Printf("Adım 3/4: AI ile metin düzenleniyor...")
	refinedText, err := s.aiClient.RefineTranscription(rawTranscriptText)
	if err != nil {
		log.Printf("AI metin düzenleme hatası: %v. Ham metin kullanılacak.", err)
		refinedText = rawTranscriptText
	}
	log.Println(">>> AI başarılı.")

	// Adım 4: Veritabanına Kaydetme
	log.Printf("Adım 4/4: Veritabanına kaydediliyor...")
	finalText := strings.TrimSpace(refinedText)
	transcription := &model.Transcription{
		UserID:            userID,
		TranscriptionText: finalText,
		AudioFilepath:     originalFilePath,
		CreatedAt:         time.Now(),
	}
	if err := s.repo.Create(transcription); err != nil {
		return nil, fmt.Errorf("transkript veritabanına kaydedilemedi: %w", err)
	}
	log.Println(">>> Kayıt başarılı.")
	return transcription, nil
}

// Geçici dosyaları silmek için yardımcı fonksiyon
func cleanupFile(path string) {
	err := os.Remove(path)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Uyarı: Geçici dosya silinemedi: %s, Hata: %v", path, err)
		}
	} else {
		log.Printf("Geçici dosya silindi: %s", path)
	}
}
