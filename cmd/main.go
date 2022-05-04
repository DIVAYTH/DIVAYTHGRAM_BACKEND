package main

import (
	"DIVAYTHGRAM_BACKEND/internal/database"
	"DIVAYTHGRAM_BACKEND/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
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
		port = "8080"
	}
	log.Println(port)
	http.ListenAndServe(":"+port, handler)
}
