package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

// Handler functions
// getBooks godoc
// @Summary Get all books
// @Description Get details of all books
// @Tags books
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} Book
// @Router /books [get]
func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, book := range books {
		if book.ID == bookId {
			return c.JSON(book)
		}
	}
	// return c.SendStatus(fiber.StatusNotFound)
	return c.Status(fiber.StatusNotFound).SendString("Not Found Krub :(")
}

func createBook(c *fiber.Ctx) error {
	book := new(Book)                          // instead *Book reserve Address already but have Book struct prototype (Be like pointer)
	if err := c.BodyParser(book); err != nil { // get value from Body request and map to book(struct)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	books = append(books, *book) // append value pointer to slide

	return c.JSON(book)
}

func updateBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for index, book := range books {
		if book.ID == bookId {
			books[index].Title = bookUpdate.Title
			books[index].Author = bookUpdate.Author
			return c.JSON(books[index])
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Not Found Krub :(")
}

func deleteBook(c *fiber.Ctx) error {
	bookId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for index, book := range books {
		if book.ID == bookId {
			// [1,2,3,4,5]
			// [1,2] + [4,5] = [1,2,4,5]
			books = append(books[:index], books[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("Not Found Krub :(")
}
