package controllers

import (
	"log"
	"ortho_vision_api/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateDisease - створення нового запису про хворобу
func CreateDisease(c *fiber.Ctx) error {
	var disease models.Disease

	// Отримуємо дані з тіла запиту
	if err := c.BodyParser(&disease); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Створюємо запис
	if err := db.Create(&disease).Error; err != nil {
		log.Println("Error saving disease:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving disease",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Disease created successfully",
		"data":    disease,
	})
}

// DeleteDisease - видалення запису про хворобу
func DeleteDisease(c *fiber.Ctx) error {
	diseaseID := c.Params("id")

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Видаляємо запис
	if err := db.Delete(&models.Disease{}, diseaseID).Error; err != nil {
		log.Println("Error deleting disease:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting disease",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Disease deleted successfully",
	})
}

// UpdateDisease - оновлення запису про хворобу
func UpdateDisease(c *fiber.Ctx) error {
	diseaseID := c.Params("id")
	var updatedData models.Disease

	// Отримуємо дані з тіла запиту
	if err := c.BodyParser(&updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Оновлюємо запис
	if err := db.Model(&models.Disease{}).Where("id = ?", diseaseID).Updates(updatedData).Error; err != nil {
		log.Println("Error updating disease:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating disease",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Disease updated successfully",
	})
}

// GetMedicalRecord - отримання всіх хвороб для призначення (appointment) за його ID
func GetMedicalRecord(c *fiber.Ctx) error {
	appointmentID := c.Params("appointmentID")

	// Перевіряємо з'єднання з базою даних
	db, ok := c.Locals("db").(*gorm.DB)
	if !ok || db == nil {
		log.Println("Database connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Database connection error",
		})
	}

	// Отримуємо всі хвороби для призначення (appointment)
	var diseases []models.Disease
	if err := db.Where("appointment_id = ?", appointmentID).Find(&diseases).Error; err != nil {
		log.Println("Error fetching diseases:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error fetching medical record",
		})
	}

	if len(diseases) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No diseases found for the specified appointment",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Medical record retrieved successfully",
		"data":    diseases,
	})
}
