package main

import (
	"fmt"
	"log"
	"net/http"
	"pbp/TugasExplore2/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := gin.Default()

	router.GET("/books", controllers.GetAllBooks)
	router.GET("/book/:book_id", controllers.GetBook)
	router.POST("/book", controllers.InsertBook)
	router.PUT("/book/:book_id", controllers.UpdateBook)
	router.DELETE("/book/:book_id", controllers.DeleteBook)

	_ = router.Run(":8888")
	main2()

}

func main2() {
	router := mux.NewRouter()

	router.HandleFunc("/booksmux", controllers.GetAllBooksMux).Methods("GET")
	router.HandleFunc("/bookmux/{book_id}", controllers.GetBookMux).Methods("GET")
	router.HandleFunc("/bookmux", controllers.InsertBookMux).Methods("POST")
	router.HandleFunc("/bookmux/{book_id}", controllers.UpdateBookMux).Methods("PUT")
	router.HandleFunc("/bookmux/{book_id}", controllers.DeleteBookMux).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}
