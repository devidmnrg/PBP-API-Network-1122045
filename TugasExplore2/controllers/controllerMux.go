package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "pbp/tugasexplore2/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func GetAllBooksMux(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
		http.Error(w, "Something has gone wrong with the Book query", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var books []m.Book
	for rows.Next() {
		var book m.Book
		if err := rows.Scan(&book.BookId, &book.BookName, &book.Pages, &book.Year); err != nil {
			log.Println(err)
			http.Error(w, "Book not found", http.StatusBadRequest)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookMux(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["book_id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid book_id", http.StatusBadRequest)
		return
	}

	var book m.Book
	err = db.QueryRow("SELECT * FROM books WHERE book_id=?", bookID).Scan(&book.BookId, &book.BookName, &book.Pages, &book.Year)
	if err != nil {
		log.Println(err)
		http.Error(w, "Book not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func InsertBookMux(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var book m.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO books (book_name, pages, year) VALUES (?, ?, ?)", book.BookName, book.Pages, book.Year)
	if err != nil {
		log.Println(err)
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Book added successfully")
}

func UpdateBookMux(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	var book m.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["book_id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid book_id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE books SET book_name=?, pages=?, year=? WHERE book_id=?", book.BookName, book.Pages, book.Year, bookID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book updated successfully")
}

func DeleteBookMux(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	params := mux.Vars(r)
	bookID, err := strconv.Atoi(params["book_id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid book_id", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM books WHERE book_id=?", bookID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book deleted successfully")
}
