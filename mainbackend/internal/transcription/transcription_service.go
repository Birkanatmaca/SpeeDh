package transcription

import "mainbackend/internal/model"

// ITranscriptionService arayüzü, transkript ile ilgili iş mantığı operasyonlarını tanımlar.
type ITranscriptionService interface {
	SaveTranscription(userID uint, text, originalFilename, filepath string) (*model.Transcription, error)
}

// transcriptionService, repository'yi kullanarak işlemleri gerçekleştirir.
type transcriptionService struct {
	repo ITranscriptionRepository
}

// NewTranscriptionService, yeni bir servis örneği oluşturur.
func NewTranscriptionService(repo ITranscriptionRepository) ITranscriptionService {
	return &transcriptionService{repo: repo}
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
