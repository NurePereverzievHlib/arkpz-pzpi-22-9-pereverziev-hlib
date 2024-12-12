package controllers

import (
	"ortho_vision_api/models"

	"github.com/gofiber/fiber/v2"
)

// Створення прийому
func CreateAppointment(c *fiber.Ctx) error {
	var appointment models.Appointment
	if err := c.BodyParser(&appointment); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}
	// Логіка для створення прийому
	return c.Status(fiber.StatusCreated).JSON(appointment)
}

// Отримання всіх прийомів
func GetAppointments(c *fiber.Ctx) error {
	var appointments []models.Appointment
	// Логіка для отримання всіх прийомів з БД
	return c.Status(fiber.StatusOK).JSON(appointments)
}

// Оновлення прийому
func UpdateAppointment(c *fiber.Ctx) error {
	// Логіка для оновлення прийому
	return c.Status(fiber.StatusOK).SendString("Appointment updated")
}

// Видалення прийому
func DeleteAppointment(c *fiber.Ctx) error {
	// Логіка для видалення прийому
	return c.Status(fiber.StatusOK).SendString("Appointment deleted")
}
