package controllers

import (
	"ortho_vision_api/config"
	"ortho_vision_api/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AddSmartGlassesData - Функція для додавання нових даних смарт-окулярів
func AddSmartGlassesData(c *fiber.Ctx) error {
	// Завжди використовуємо user_id = 4
	userID := 4

	// Отримуємо дані з тіла запиту
	var data models.SmartGlassesData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data",
		})
	}

	// Перевірка наявності користувача в базі даних
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Встановлюємо дані для запису
	data.UserID = uint(userID)
	data.Timestamp = time.Now()

	// Додаємо новий запис у таблицю
	if err := config.DB.Create(&data).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create smart glasses data",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(data)
}

// GetSmartGlassesStatistics - Функція для отримання статистики по даних смарт-окулярів за вказаний день
func GetSmartGlassesStatistics(c *fiber.Ctx) error {
	userID := 4

	// Отримуємо дату з тіла запиту (потрібно, щоб дата була у форматі "YYYY-MM-DD")
	dateParam := c.Query("date") // Читаємо параметр "date" з запиту

	// Перевірка на коректність формату дати
	date, err := time.Parse("2006-01-02", dateParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD.",
		})
	}

	// Отримуємо дані з таблиці smart_glasses_data для заданого користувача
	var data []models.SmartGlassesData
	if err := config.DB.Where("user_id = ?", userID).Find(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No data found for the user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch smart glasses data",
		})
	}

	// Ініціалізація змінних для підрахунку часу
	var timeHeadTiltExceeds45 float64
	var timeLowLight float64
	var timeHighLight float64

	// Змінні для збереження часу початку порушення норми
	var startHeadTiltExceeds45 time.Time
	var startLowLight time.Time
	var startHighLight time.Time

	// Перебираємо всі записи
	for _, entry := range data {
		// Перевіряємо, чи є дані для кожного запису
		if entry.Timestamp.IsZero() || entry.PostureAngle == 0 || entry.EyeStrain == 0 {
			continue // Пропускаємо цей запис, якщо є некоректні дані
		}

		// Перевірка, чи потрапляє запис в заданий день
		// Округлюємо час до початку дня і перевіряємо, чи він співпадає
		entryDate := entry.Timestamp.Truncate(24 * time.Hour) // Округлюємо до початку дня
		if entryDate != date {
			continue // Пропускаємо записи, які не належать до вказаного дня
		}

		// Якщо нахил голови більше 45°
		if entry.PostureAngle > 45 {
			// Якщо раніше не було початку порушення, фіксуємо його
			if startHeadTiltExceeds45.IsZero() {
				startHeadTiltExceeds45 = entry.Timestamp
			}
		} else {
			// Якщо нахил голови став нормальним і був запис початку
			if !startHeadTiltExceeds45.IsZero() {
				duration := entry.Timestamp.Sub(startHeadTiltExceeds45).Seconds()
				timeHeadTiltExceeds45 += duration
				startHeadTiltExceeds45 = time.Time{} // Скидаємо початок порушення
			}
		}

		// Якщо освітлення менше 100 lux
		if entry.EyeStrain < 100 {
			// Якщо раніше не було початку порушення, фіксуємо його
			if startLowLight.IsZero() {
				startLowLight = entry.Timestamp
			}
		} else {
			// Якщо освітлення стало нормальним і був запис початку
			if !startLowLight.IsZero() {
				duration := entry.Timestamp.Sub(startLowLight).Seconds()
				timeLowLight += duration
				startLowLight = time.Time{} // Скидаємо початок порушення
			}
		}

		// Якщо освітлення більше 1000 lux
		if entry.EyeStrain > 1000 {
			// Якщо раніше не було початку порушення, фіксуємо його
			if startHighLight.IsZero() {
				startHighLight = entry.Timestamp
			}
		} else {
			// Якщо освітлення стало нормальним і був запис початку
			if !startHighLight.IsZero() {
				duration := entry.Timestamp.Sub(startHighLight).Seconds()
				timeHighLight += duration
				startHighLight = time.Time{} // Скидаємо початок порушення
			}
		}
	}

	// Створюємо відповідь з підрахованим часом
	stats := fiber.Map{
		"time_head_tilt_exceeds_45": timeHeadTiltExceeds45,
		"time_low_light":            timeLowLight,
		"time_high_light":           timeHighLight,
	}

	// Повертаємо результат
	return c.Status(fiber.StatusOK).JSON(stats)
}
