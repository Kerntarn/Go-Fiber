package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	dbname   = "mydatabase"
	username = "myuser"
	password = "mypassword"
)

var db *sql.DB

func dbInit() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable\n", host, port, username, password, dbname)

	print(psqlInfo)

	ldb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = ldb.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db = ldb
	fmt.Println("Database Connection: Success")

	book, err := dbUpdateBook(&Book{Title: "Demon Slayer", Price: 310}, 1)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(book)
}

func dbCreateBook(book *Book) error {
	_, err := db.Exec( //Not Return anything
		"INSERT INTO public.books(title, price) VALUES ($1, $2);",
		book.Title,
		book.Price,
	)

	return err
}

func dbGetBook(id int) (Book, error) {
	var b Book
	row := db.QueryRow("SELECT id, title, price FROM books WHERE id=$1;", id)
	err := row.Scan(&b.ID, &b.Title, &b.Price)
	if err != nil {
		return Book{}, err
	}
	return b, nil
}

func dbUpdateBook(book *Book, id int) (Book, error) {
	var b Book
	row := db.QueryRow( //Return Data
		"UPDATE public.books SET title=$1, price=$2 WHERE id=$3 RETURNING id, title, price;",
		book.Title,
		book.Price,
		id,
	)
	err := row.Scan(&b.ID, &b.Title, &b.Price)
	if err != nil {
		return Book{}, err
	}
	return b, err
}

func dbDeleteBook(id int) error {
	_, err := db.Exec( //Not Return anything
		"DELETE FROM public.books WHERE id=$1;",
		id,
	)

	return err
}
