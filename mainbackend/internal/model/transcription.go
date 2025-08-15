package model

import "time"

type Transcription struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Title             string    `gorm:"not null" json:"title"`
	OriginalFilename  string    `json:"original_filename"`
	AudioFilepath     string    `gorm:"not null" json:"-"`                   // Ses dosyasının sunucudaki yolu
	TranscriptionText string    `gorm:"type:text" json:"transcription_text"` // Uzun metinler için
	CreatedAt         time.Time `json:"created_at"`

	// Bir transkript bir kullanıcıya aittir.
	UserID uint `gorm:"not null" json:"user_id"`
	User   User `json:"-"` // GORM ilişkisi için
}
