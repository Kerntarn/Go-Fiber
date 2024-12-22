package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	dbname   = "mydatabase"
	username = "myuser"
	password = "mypassword"
)

var db *gorm.DB

func dbInit() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable\n", host, port, username, password, dbname)

	ldb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect to database")
	}

	db = ldb

	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Database migration completed!")

}

func dbCreateBook(book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}

	fmt.Println("Create Book Successful")
	return nil
}

func dbGetBook(id int) (Book, error) {
	var b Book
	result := db.First(&b, id)
	if result.Error != nil {
		return Book{}, result.Error
	}
	return b, nil
}

func dbUpdateBook(book *Book) error {
	result := db.Model(&book).Updates(book)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func dbDeleteBook(id int) error {
	var b Book
	result := db.Delete(&b, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func dbGetBooks() ([]Book, error) {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}

	return books, nil
}
