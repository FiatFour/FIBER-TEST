package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/fiatfour/fiber-test/docs" // load generated docs
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func checkMiddleware(c *fiber.Ctx) error {
	// start := time.Now()
	// fmt.Printf("URL = %s, Method = %s, Time = %s\n", c.OriginalURL(), c.Method(), start)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	fmt.Print(claims)
	return c.Next()
}

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html") // path folder of html file

	// Pass engine to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	books = append(books, Book{ID: 1, Title: "Nickname", Author: "Fiat"})
	books = append(books, Book{ID: 2, Title: "Full name", Author: "Anfat Nilaingan"})

	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Use(checkMiddleware)

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)
	app.Get("test-html", testHTML)

	app.Get("/config", getEnv)

	app.Listen("localhost:8080")
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/"+file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File upload complete!")
}

func testHTML(c *fiber.Ctx) error {
	// Render index template
	return c.Render("index", fiber.Map{ // index.html file
		"Title": "Hello, World!",
		"Name":  "Fiat",
	})
}

func getEnv(c *fiber.Ctx) error {
	// if value, exists := os.LookupEnv("SECRET"); exists {
	// 	return c.JSON(fiber.Map{
	// 		"SECRET": value,
	// 	})
	// }
	secret := os.Getenv("SECRET")

	if secret == "" {
		secret = "defaultsecret"
	}

	return c.JSON(fiber.Map{
		"SECRET": secret,
	})
}

// Dummy user for example
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var memberUser = User{
	Email:    "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// user, pass don't match (Unauthorized)
	if user.Email != memberUser.Email || user.Password != memberUser.Password {
		return fiber.ErrUnauthorized
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "Login success",
		"token":   t,
	})
}
