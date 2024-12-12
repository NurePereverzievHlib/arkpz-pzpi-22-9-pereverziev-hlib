package models

import "time"

// Модель для таблиці Clinics
type Clinic struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	Phone     string    `gorm:"size:20"`
	Location  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
