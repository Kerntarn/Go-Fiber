package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/Kerntarn/Go-Fiber/docs" // load generated docs
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func srv() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error")
	}
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use("/books", authRequired)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "*",
		AllowHeaders: "*",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Post("/register", register)
	app.Post("/login", login)

	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(os.Getenv("JWT_SECRET")),
	// }))

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)
	app.Listen(":3000")
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
	return c.SendString("File's uploaded completely")
}

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	t, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !t.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim := t.Claims.(jwt.MapClaims)

	fmt.Println(claim)
	return c.Next()
}

func login(c *fiber.Ctx) error {
	user := new(User)
	userInDb := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	result := db.Where("email = ?", user.Email).First(userInDb)
	if result.Error != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(userInDb.Password), []byte(user.Password)); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userInDb.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"message": "Login Success",
		"token":   t,
	})
}

func register(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := dbCreateUser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(fiber.Map{
		"message": "Register Successful",
	})
}
