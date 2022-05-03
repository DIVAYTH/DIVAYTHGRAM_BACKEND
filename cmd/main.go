package main

import (
	"DIVAYTHGRAM_BACKEND/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {
	//database.Init()
	r := mux.NewRouter()
	handlers.InitHandlers(r)
	handler := cors.AllowAll().Handler(r)
	fileServer := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))
	http.ListenAndServe(os.Getenv("PORT"), handler)
}
