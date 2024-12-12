package models

import (
	"time"
)

// Модель для таблиці SmartGlassesData
type SmartGlassesData struct {
	ID                     uint      `gorm:"primary_key"`
	UserID                 uint      `gorm:"not null"`
	PostureAngle           int       `gorm:"not null;check:posture_angle >= 0 AND posture_angle <= 180"`
	EyeStrain              int       `gorm:"not null;check:eye_strain >= 0 AND eye_strain <= 100"`
	LowLightDuration       int       `gorm:"default:0"`
	HighLightDuration      int       `gorm:"default:0"`
	IncorrectPostureDuration int    `gorm:"default:0"`
	Timestamp              time.Time `gorm:"not null"`
}
