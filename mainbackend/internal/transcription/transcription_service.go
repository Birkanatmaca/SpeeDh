package transcription

import (
	"errors"
	"fmt"
	"mainbackend/internal/model"
	"os"
)

// ITranscriptionService arayüzü, transkript ile ilgili iş mantığı operasyonlarını tanımlar.
type ITranscriptionService interface {
	SaveTranscription(userID uint, text, originalFilename, filepath string) (*model.Transcription, error)
	GetTranscriptsByUserID(userID uint) ([]model.Transcription, error)
	GetTranscriptByID(id uint) (*model.Transcription, error)
	DeleteTranscript(id uint, userID uint) error
}

// transcriptionService, repository'yi kullanarak işlemleri gerçekleştirir.
type transcriptionService struct {
	repo ITranscriptionRepository
}

// NewTranscriptionService, yeni bir servis örneği oluşturur.
func NewTranscriptionService(repo ITranscriptionRepository) ITranscriptionService {
	return &transcriptionService{repo: repo}
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
