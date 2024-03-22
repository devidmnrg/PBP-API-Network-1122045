package controllers

type Book struct {
	BookId   int    `json:"bookid"`
	BookName string `json:"bookname"`
	Pages    int    `json:"pages"`
	Year     int    `json:"int"`
}

type BookResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Book   `json:"data"`
}

type BooksResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Book `json:"data"`
}
