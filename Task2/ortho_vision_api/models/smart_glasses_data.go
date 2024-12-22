package models

import (
	"time"
)

// Модель для таблиці SmartGlassesData
type SmartGlassesData struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id"`
	PostureAngle float64   `json:"posture_angle"`
	EyeStrain    float64   `json:"eye_strain"`
	Timestamp    time.Time `json:"timestamp"`
}

// Вказуємо правильну назву таблиці
func (SmartGlassesData) TableName() string {
	return "smartglassesdata" // вказуємо точну назву таблиці
}
