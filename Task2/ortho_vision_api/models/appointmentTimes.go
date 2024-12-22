package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// CustomTime - це кастомний тип для роботи з часом.
type CustomTime time.Time

// Scan для CustomTime реалізує інтерфейс Scanner, щоб використовувати цей тип у базі даних.
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime(time.Time{})
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan CustomTime")
	}
	*ct = CustomTime(t)
	return nil
}

// Value для CustomTime реалізує інтерфейс Valuer, щоб використовувати цей тип при збереженні в базу.
func (ct CustomTime) Value() (driver.Value, error) {
	return time.Time(ct).Format("2006-01-02 15:04:05"), nil
}

// Структура AppointmentTimes
type AppointmentTimes struct {
	ID            uint       `gorm:"primary_key"`
	DoctorID      uint       `gorm:"not null;index"`
	Doctor        User       `gorm:"foreignkey:DoctorID"`
	ClinicID      uint       `gorm:"not null;index" json:"clinic_id"`
	Clinic        Clinic     `gorm:"foreignkey:ClinicID"`
	AvailableTime CustomTime `gorm:"not null" json:"available_time"`
	IsBooked      bool       `gorm:"column:is_booked;default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
