package transcription

import (
	"gorm.io/gorm"
	"mainbackend/internal/model"
)

// ITranscriptionRepository arayüzü, transkript veritabanı işlemleri için metotları tanımlar.
type ITranscriptionRepository interface {
	Create(transcription *model.Transcription) error
	FindByUserID(userID uint) ([]model.Transcription, error)
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
func (r *transcriptionRepository) FindByUserID(userID uint) ([]model.Transcription, error) {
	var transcriptions []model.Transcription
	// 'user_id' alanı userID ile eşleşen kayıtları bul ve CreatedAt'a göre en yeniden eskiye doğru sırala
	result := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&transcriptions)
	return transcriptions, result.Error
}
