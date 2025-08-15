package transcription

import (
	"gorm.io/gorm"
	"mainbackend/internal/model"
)

// ITranscriptionRepository arayüzü, transkript veritabanı işlemleri için metotları tanımlar.
type ITranscriptionRepository interface {
	Create(transcription *model.Transcription) error
}

// transcriptionRepository, veritabanı bağlantısını tutar.
type transcriptionRepository struct {
	db *gorm.DB
}

// NewTranscriptionRepository, yeni bir repository örneği oluşturur.
func NewTranscriptionRepository(db *gorm.DB) ITranscriptionRepository {
	return &transcriptionRepository{db: db}
}

// Create, veritabanına yeni bir transkript kaydı ekler.
func (r *transcriptionRepository) Create(transcription *model.Transcription) error {
	return r.db.Create(transcription).Error
}
