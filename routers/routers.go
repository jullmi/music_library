package routers

import (
	"music_library/handlers"
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/songs", handlers.GetList).Methods("GET")
	r.HandleFunc("/api/songs", handlers.AddSong).Methods("POST")
	r.HandleFunc("", handlers.DeleteSong).Methods("DELETE")
	r.HandleFunc("", handlers.UpdateSong).Methods("PUT")
	r.HandleFunc("", handlers.GetSongTextWithPagination).Methods("GET")

	return r
}