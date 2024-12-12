package controllers

import (
	"log"
	"ortho_vision_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateAppointmentTime(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID лікаря з параметрів
	doctorID := c.Params("doctor_id")
	if doctorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Doctor ID is required",
		})
	}

	// Перевіряємо, чи існує лікар
	var doctor models.User
	if err := db.First(&doctor, "id = ? AND role = ?", doctorID, "doctor").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Doctor with this ID does not exist or is not a doctor")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Doctor not found or user is not a doctor",
			})
		}
		log.Println("Error finding doctor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error verifying doctor",
		})
	}

	// Отримуємо дані з тіла запиту
	var appointmentTime models.AppointmentTimes
	if err := c.BodyParser(&appointmentTime); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Перевіряємо формат доступного часу
	availableTimeStr := c.FormValue("available_time")
	if availableTimeStr != "" {
		// Проводимо парсинг часу в форматі ISO 8601
		availableTime, err := time.Parse(time.RFC3339, availableTimeStr)
		if err != nil {
			log.Println("Error parsing available_time:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid AvailableTime format",
			})
		}
		appointmentTime.AvailableTime = availableTime
	}

	// Перевірка на правильність значення available_time
	if appointmentTime.AvailableTime.IsZero() {
		log.Println("Invalid or missing AvailableTime")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "AvailableTime is required or invalid",
		})
	}

	// Встановлюємо ID лікаря
	appointmentTime.DoctorID = doctor.ID

	// Перевіряємо, чи клініка існує
	var clinic models.Clinic
	if err := db.First(&clinic, "id = ?", appointmentTime.ClinicID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Clinic not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Clinic not found",
			})
		}
		log.Println("Error finding clinic:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error verifying clinic",
		})
	}

	// Зберігаємо доступний час у базу
	if err := db.Create(&appointmentTime).Error; err != nil {
		log.Println("Error saving appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving appointment time",
		})
	}

	// Відповідь про успішне створення
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Appointment time created successfully",
		"data":    appointmentTime,
	})
}

// UpdateAppointmentTime - редагування доступного часу для лікаря
func UpdateAppointmentTime(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID лікаря та ID часу з параметрів
	doctorID := c.Params("doctor_id")
	appointmentTimeID := c.Params("appointment_time_id")
	if doctorID == "" || appointmentTimeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Doctor ID and Appointment Time ID are required",
		})
	}

	var doctor models.User
	if err := db.First(&doctor, "id = ? AND role = ?", doctorID, "doctor").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Doctor with this ID does not exist or is not a doctor")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Doctor not found or user is not a doctor",
			})
		}
		log.Println("Error finding doctor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error verifying doctor",
		})
	}

	var appointmentTime models.AppointmentTimes
	if err := db.First(&appointmentTime, "id = ? AND doctor_id = ?", appointmentTimeID, doctorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Appointment time not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Appointment time not found",
			})
		}
		log.Println("Error finding appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding appointment time",
		})
	}

	// Отримуємо нові дані з тіла запиту
	if err := c.BodyParser(&appointmentTime); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Оновлюємо доступний час
	if err := db.Save(&appointmentTime).Error; err != nil {
		log.Println("Error updating appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating appointment time",
		})
	}

	// Відповідь про успішне редагування
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Appointment time updated successfully",
		"data":    appointmentTime,
	})
}

// GetAllAppointmentTimesForDoctor - отримання всіх доступних часів для лікаря
func GetAllAppointmentTimesForDoctor(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID лікаря з параметрів
	doctorID := c.Params("doctor_id")
	if doctorID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Doctor ID is required",
		})
	}

	var doctor models.User
	if err := db.First(&doctor, "id = ? AND role = ?", doctorID, "doctor").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Doctor with this ID does not exist or is not a doctor")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Doctor not found or user is not a doctor",
			})
		}
		log.Println("Error finding doctor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error verifying doctor",
		})
	}

	var appointmentTimes []models.AppointmentTimes
	if err := db.Where("doctor_id = ?", doctorID).Find(&appointmentTimes).Error; err != nil {
		log.Println("Error retrieving appointment times:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error retrieving appointment times",
		})
	}

	// Відповідь про успішне отримання всіх доступних часів
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Appointment times retrieved successfully",
		"data":    appointmentTimes,
	})
}

// DeleteAppointmentTime - видалення доступного часу для лікаря
func DeleteAppointmentTime(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID лікаря та ID часу з параметрів
	doctorID := c.Params("doctor_id")
	appointmentTimeID := c.Params("appointment_time_id")
	if doctorID == "" || appointmentTimeID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Doctor ID and Appointment Time ID are required",
		})
	}

	var doctor models.User
	if err := db.First(&doctor, "id = ? AND role = ?", doctorID, "doctor").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Doctor with this ID does not exist or is not a doctor")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Doctor not found or user is not a doctor",
			})
		}
		log.Println("Error finding doctor:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error verifying doctor",
		})
	}

	var appointmentTime models.AppointmentTimes
	if err := db.First(&appointmentTime, "id = ? AND doctor_id = ?", appointmentTimeID, doctorID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Appointment time not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Appointment time not found",
			})
		}
		log.Println("Error finding appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding appointment time",
		})
	}

	// Видаляємо запис
	if err := db.Delete(&appointmentTime).Error; err != nil {
		log.Println("Error deleting appointment time:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting appointment time",
		})
	}

	// Відповідь про успішне видалення
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Appointment time deleted successfully",
	})
}
