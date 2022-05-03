package main

import (
	"DIVAYTHGRAM_BACKEND/internal/database"
	"DIVAYTHGRAM_BACKEND/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	database.Init()
	r := mux.NewRouter()
	handlers.InitHandlers(r)
	handler := cors.AllowAll().Handler(r)
	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println(port)
	http.ListenAndServe(port, handler)
}
