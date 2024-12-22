package models

import (
	"time"
)

// Модель для таблиці Users
type User struct {
	ID           uint   `gorm:"primary_key"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"not null;check:role in ('patient', 'admin', 'doctor')"`
	CreatedAt    time.Time
	Password     string `gorm:"-"`
}
