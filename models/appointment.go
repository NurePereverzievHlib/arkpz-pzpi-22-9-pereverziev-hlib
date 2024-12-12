package models

import (
	"time"
)

// Модель для таблиці Appointments
type Appointment struct {
	ID                uint   `gorm:"primary_key"`
	AppointmentTimeID uint   `gorm:"not null"` // Посилання на доступний час
	Status            string `gorm:"not null;default:'pending';check:status IN ('pending', 'confirmed', 'cancelled', 'completed')"`
	Reason            string `gorm:""` // Причина запису
	CreatedAt         time.Time
}
