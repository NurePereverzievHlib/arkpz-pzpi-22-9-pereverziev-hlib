package models

import (
	"time"
)

type AppointmentTimes struct {
	ID            uint      `gorm:"primary_key"`
	DoctorID      uint      `gorm:"not null;index"`      // Зв'язок з лікарем
	Doctor        User      `gorm:"foreignkey:DoctorID"` // Лікар (припускається, що модель Users є для лікарів)
	ClinicID      uint      `gorm:"not null;index"`      // Зв'язок з клінікою
	Clinic        Clinic    `gorm:"foreignkey:ClinicID"` // Клініка (припускається, що є модель Clinic)
	AvailableTime time.Time `gorm:"not null"`            // Доступний час прийому
	IsBooked      bool      `gorm:"default:false"`       // Чи заброньований цей час
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
