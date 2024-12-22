package controllers

import (
	"log"
	"ortho_vision_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateAppointment - створення запису на прийом
func CreateAppointment(c *fiber.Ctx) error {
	var appointment models.Appointment

	// Отримуємо дані з тіла запиту
	if err := c.BodyParser(&appointment); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	// Перевіряємо, чи є доступний час
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Перевіряємо доступний час
	var availableTime models.AppointmentTimes
	err := db.Where("id = ? AND is_booked = ?", appointment.AppointmentTimeID, false).First(&availableTime).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "The selected time is not available",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error checking available time",
		})
	}

	// Якщо час доступний, змінюємо його статус на "booked"
	availableTime.IsBooked = true
	if err := db.Save(&availableTime).Error; err != nil {
		log.Println("Error updating appointment time status:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating appointment time status",
		})
	}

	// Встановлюємо статус прийому на "pending" і додаємо новий запис
	appointment.Status = "pending"
	if err := db.Create(&appointment).Error; err != nil {
		log.Println("Error saving appointment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving appointment",
		})
	}

	// Відповідь про успішний запис
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Appointment created successfully",
		"data":    appointment,
	})
}

// DeleteAppointment - видалення запису на прийом
func DeleteAppointment(c *fiber.Ctx) error {
	// Отримуємо ID запису на прийом з параметрів запиту
	appointmentID := c.Params("id")

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Знаходимо запис на прийом
	var appointment models.Appointment
	err := db.First(&appointment, "id = ?", appointmentID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Appointment not found",
			})
		}
		log.Println("Error finding appointment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding appointment",
		})
	}

	// Оновлюємо статус часу на "available" (is_booked = false)
	var availableTime models.AppointmentTimes
	err = db.First(&availableTime, "id = ?", appointment.AppointmentTimeID).Error
	if err != nil {
		log.Println("Error finding appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding appointment time",
		})
	}

	availableTime.IsBooked = false
	if err := db.Save(&availableTime).Error; err != nil {
		log.Println("Error updating appointment time status:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating appointment time status",
		})
	}

	// Видаляємо запис на прийом
	if err := db.Delete(&appointment).Error; err != nil {
		log.Println("Error deleting appointment:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting appointment",
		})
	}

	// Відповідь про успішне видалення
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Appointment deleted successfully",
	})
}

// GetAppointmentsByPatientID - отримання всіх записів на прийом для конкретного пацієнта
func GetAppointmentsByPatientID(c *fiber.Ctx) error {
	// Отримуємо ID пацієнта з параметрів запиту
	patientID := c.Params("patientID")

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Знаходимо всі записи для пацієнта
	var appointments []models.Appointment
	err := db.Where("patient_id = ?", patientID).Find(&appointments).Error
	if err != nil {
		log.Println("Error finding appointments:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding appointments",
		})
	}

	// Якщо записів немає, повертаємо відповідь про їх відсутність
	if len(appointments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No appointments found for the specified patient",
		})
	}

	// Повертаємо знайдені записи
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Appointments retrieved successfully",
		"data":    appointments,
	})
}
