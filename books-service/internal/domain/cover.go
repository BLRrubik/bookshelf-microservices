package domain

import "time"

// CoverStatus — статус обработки обложки
type CoverStatus string

const (
	CoverStatusNone       CoverStatus = "none"
	CoverStatusProcessing CoverStatus = "processing"
	CoverStatusReady      CoverStatus = "ready"
	CoverStatusFailed     CoverStatus = "failed"
)

type Cover struct {
	ID           string      `json:"id" db:"id"`
	BookID       string      `json:"book_id" db:"book_id"`
	Status       CoverStatus `json:"status" db:"status"`
	OriginalPath string      `json:"-" db:"original_path"`       // Путь к оригиналу в MinIO
	CoverPath    string      `json:"-" db:"cover_path"`          // Путь к обложке (400x600)
	ThumbPath    string      `json:"-" db:"thumb_path"`          // Путь к миниатюре (100x150)
	CoverURL     string      `json:"cover_url,omitempty"`        // Публичный URL обложки
	ThumbURL     string      `json:"thumb_url,omitempty"`        // Публичный URL миниатюры
	Error        string      `json:"error,omitempty" db:"error"` // Сообщение об ошибке
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`
	CompletedAt  *time.Time  `json:"completed_at,omitempty" db:"completed_at"`
}

// CoverUploadResponse — ответ на загрузку (202 Accepted)
type CoverUploadResponse struct {
	CoverID string      `json:"cover_id"`
	Status  CoverStatus `json:"status"`
	Message string      `json:"message"`
}

// CoverResponse — информация об обложке
type CoverResponse struct {
	Status   CoverStatus `json:"status"`
	CoverURL string      `json:"cover_url,omitempty"`
	ThumbURL string      `json:"thumb_url,omitempty"`
	Message  string      `json:"message,omitempty"`
}

// CoverStatusResponse — детальный статус для polling
type CoverStatusResponse struct {
	CoverID     string      `json:"cover_id"`
	Status      CoverStatus `json:"status"`
	CoverURL    string      `json:"cover_url,omitempty"`
	ThumbURL    string      `json:"thumb_url,omitempty"`
	Error       string      `json:"error,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	CompletedAt *time.Time  `json:"completed_at,omitempty"`
}
