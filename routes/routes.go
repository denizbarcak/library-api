package routes

import (
    "github.com/gofiber/fiber/v2"
    "library-api/handlers"
)

func UserRoutes(app *fiber.App) {
    app.Post("/users", handlers.CreateUser)
    app.Get("/users", handlers.GetUsers)
    app.Delete("/users/:id", handlers.DeleteUser)
    app.Put("/users/:id/books", handlers.AddBookToUser) // Kitap ekleme rotası
    app.Post("/login", handlers.Login)
}

func BookRoutes(app *fiber.App) {
    app.Post("/books", handlers.CreateBook)
    app.Get("/books", handlers.GetBooks)
    app.Delete("/books/:id", handlers.DeleteBook)
}
