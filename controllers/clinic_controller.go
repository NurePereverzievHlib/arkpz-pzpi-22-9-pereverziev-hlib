package controllers

import (
	"log"
	"ortho_vision_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AddClinic - функція для додавання нової клініки
func AddClinic(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо нову структуру для клініки
	var clinic models.Clinic

	// Парсимо тіло запиту в структуру clinic
	if err := c.BodyParser(&clinic); err != nil {
		log.Println("BodyParser error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Додаткові перевірки, наприклад, наявність обов'язкових полів
	if clinic.Name == "" || clinic.Address == "" || clinic.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name, Address, and Location are required fields",
		})
	}

	// Заповнюємо поля часу створення та оновлення
	clinic.CreatedAt = time.Now()
	clinic.UpdatedAt = time.Now()

	// Зберігаємо клініку в базу даних
	if err := db.Create(&clinic).Error; err != nil {
		log.Println("Database save error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving clinic to database",
		})
	}

	// Відповідь про успішне додавання клініки
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Clinic added successfully",
		"clinic":  clinic,
	})
}

// GetAllClinics - функція для отримання всіх клінік
func GetAllClinics(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо всі клініки з бази
	var clinics []models.Clinic
	if err := db.Find(&clinics).Error; err != nil {
		log.Println("Error fetching clinics:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching clinics",
		})
	}

	// Повертаємо список клінік
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Clinics retrieved successfully",
		"clinics": clinics,
	})
}

// GetClinicByName - функція для отримання клініки за назвою
func GetClinicByName(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо назву клініки з параметрів
	clinicName := c.Params("name")
	if clinicName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Clinic name is required",
		})
	}

	// Пошук клініки за назвою
	var clinic models.Clinic
	if err := db.Where("name = ?", clinicName).First(&clinic).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Clinic not found",
			})
		}
		log.Println("Error fetching clinic by name:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching clinic",
		})
	}

	// Повертаємо знайдену клініку
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Clinic retrieved successfully",
		"clinic":  clinic,
	})
}

func DeleteClinic(c *fiber.Ctx) error {
	// Отримуємо з'єднання з базою даних із контексту
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо ID клініки з параметрів
	clinicID := c.Params("id")
	if clinicID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Clinic ID is required",
		})
	}

	// Видалення клініки за ID
	if err := db.Delete(&models.Clinic{}, clinicID).Error; err != nil {
		log.Println("Error deleting clinic:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting clinic",
		})
	}

	// Відповідь про успішне видалення
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Clinic deleted successfully",
	})
}
