package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "pbp/tugasexplore2/models"

	"github.com/gin-gonic/gin"
)

func GetAllBooks(c *gin.Context) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": "Something has gone wrong with the Movie query"})
		return
	}

	var book m.Book
	for rows.Next() {
		if err := rows.Scan(&book.BookId, &book.BookName, &book.Pages, &book.Year); err != nil {
			log.Println(err)
			SendErrorResponseGIN(c, 400, "Error: Book not found!")
		} else {
			c.IndentedJSON(http.StatusOK, book)
		}
	}
}

func GetBook(c *gin.Context) {
	db := connect()
	defer db.Close()

	idBook := c.Param("book_id")

	var book m.Book
	err := db.QueryRow("SELECT * FROM books WHERE book_id=?", idBook).Scan(&book.BookId, &book.BookName, &book.Pages, &book.Year)
	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 400, "Error: Book not found!")
		return
	}

	c.JSON(http.StatusOK, book)
}

func InsertBook(c *gin.Context) {
	db := connect()
	defer db.Close()

	var book m.Book

	book.BookName = c.Query("book_name")
	pagesStr := c.Query("pages")
	yearStr := c.Query("year")

	pages, err := strconv.Atoi(pagesStr)
	if err != nil {
		SendErrorResponseGIN(c, 404, "Error: invalid pages value!")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		SendErrorResponseGIN(c, 404, "Error: invalid year value!")
		return
	}

	_, err = db.Exec("INSERT INTO books (book_name, pages, year) VALUES (?, ?, ?)", book.BookName, pages, year)
	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 500, "Error: iinsert failed")
		return
	}

	SendSuccessResponseGIN(c, 200, "Insert Success!")
}

func UpdateBook(c *gin.Context) {
	db := connect()
	defer db.Close()

	var book m.Book
	idBook := c.Param("book_id")

	book.BookName = c.PostForm("book_name")
	pagesStr := c.PostForm("pages")
	yearStr := c.PostForm("year")

	id, err := strconv.Atoi(idBook)
	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 404, "Error: invalid book_id!")
		return
	}
	pages, err := strconv.Atoi(pagesStr)
	if err != nil {
		SendErrorResponseGIN(c, 404, "Error: invalid pages value!")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		SendErrorResponseGIN(c, 404, "Error: invalid year value!")
		return
	}

	if err := c.Bind(&book); err != nil {
		fmt.Println(err)
		SendErrorResponseGIN(c, 500, "Error: internal server error!")
		return
	}

	result, err := db.Exec("UPDATE books SET book_name=?, pages=?, year=? WHERE book_id=?", book.BookName, pages, year, id)

	num, _ := result.RowsAffected()

	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 500, "Error: internal server error!")
		return
	}

	if num == 0 {
		log.Println(err)
		SendErrorResponseGIN(c, 404, "Error: no rows affected!")
		return
	}

	SendSuccessResponseGIN(c, 200, "Update Success!")
}

func DeleteBook(c *gin.Context) {
	db := connect()
	defer db.Close()

	var book m.Book
	idBook := c.Param("book_id")
	id, err := strconv.Atoi(idBook)
	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 404, "Error: invalid book_id!")
		return
	}

	if err := c.Bind(&book); err != nil {
		fmt.Print(err)
		SendErrorResponseGIN(c, 400, "Error: bad request!")
		return
	}

	result, err := db.Exec("DELETE FROM books WHERE book_id=?", id)

	num, _ := result.RowsAffected()

	if err != nil {
		log.Println(err)
		SendErrorResponseGIN(c, 400, "Error: delete failed!")
	} else {
		if num == 0 {
			log.Println(err)
			SendErrorResponseGIN(c, 400, "Error: no rows affected!")
		} else {
			SendSuccessResponseGIN(c, 200, "Delete Success")
		}

	}
}

func SendSuccessResponseGIN(
	c *gin.Context, code int,
	message string) {
	c.JSON(code,
		gin.H{
			"code":    code,
			"message": message,
		})
}

func SendErrorResponseGIN(
	c *gin.Context, code int,
	message string) {
	c.JSON(code,
		gin.H{
			"code":    code,
			"message": message,
		})
}
