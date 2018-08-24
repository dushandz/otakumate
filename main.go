package main

import (
	"log"
	"net/http"
	"otakumate/comic"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/comic/top/{type}/{pageNum}", comic.TopHandler)
	r.HandleFunc("/comic/{comicID}", comic.ComicHandler)
	r.HandleFunc("/comic/vols/{volID}", comic.VolHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
	// comic.UpdateWholeData()

	//定时任务
	// c := cron.New()
	// c.AddFunc("@midnight", main)
	// c.Start()
}
