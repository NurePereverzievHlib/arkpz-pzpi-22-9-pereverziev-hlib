package models

import (
	"time"
)

// Модель для таблиці Diseases
type Disease struct {
	ID            uint      `gorm:"primary_key"`
	UserID        uint      `gorm:"not null"`
	DiseaseName   string    `gorm:"not null"`
	Description   string    `gorm:""`
	DiagnosisDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status        string    `gorm:"not null;default:'active';check:status IN ('active', 'inactive')"`
}
