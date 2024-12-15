package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	books = append(books, Book{ID: 1, Title: "Demon Slayer", Author: "A Sensei"})
	books = append(books, Book{ID: 2, Title: "One Piece", Author: "B Sensei"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Listen(":3000")
}
