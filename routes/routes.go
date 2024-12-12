package routes

import (
	"ortho_vision_api/controllers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes налаштовує маршрути для вашого API.
func SetupRoutes(app *fiber.App) {

	app.Post("/register", controllers.RegisterUser) // Реєстрація користувача

	app.Get("/login", controllers.LoginUser) // Вхід користувача

	app.Put("/users/:id", controllers.UpdateUser) // Оновлення даних користувача

	app.Delete("/admin/users/:id", controllers.DeleteUser) // Видалення користувача за ID

	app.Post("/admin/clinics", controllers.AddClinic) // Додавання клініки

	app.Get("/admin/clinics", controllers.GetAllClinics) // Отримання всіх клінік

	app.Get("/admin/clinics/:name", controllers.GetClinicByName) // Отримання клініки за назвою

	app.Delete("/admin/clinics/:id", controllers.DeleteClinic) // Видалення клініки за ID

	app.Post("/doctor/:doctor_id/appointment_times", controllers.CreateAppointmentTime) // Створення вільного часу доктора

	app.Put("/doctor/:doctor_id/appointment_times/:appointment_time_id", controllers.UpdateAppointmentTime) // Редагування вільного часу доктора

	app.Get("/doctor/:doctor_id/appointment_times", controllers.CreateAppointmentTime) // Отримання всіх вільних місць доктора

	app.Delete("/doctor/:doctor_id/appointment_times/:appointment_time_id", controllers.DeleteAppointmentTime) // Видалення конкретного вільного часу

	// Запис на прийом
	// app.Post("/appointments", controllers.CreateAppointment)

	// // Отримання всіх прийомів
	// app.Get("/appointments", controllers.GetAppointments)

	// // Оновлення інформації про прийом
	// app.Put("/appointments/:id", controllers.UpdateAppointment)

	// // Видалення прийому
	// app.Delete("/appointments/:id", controllers.DeleteAppointment)
}
