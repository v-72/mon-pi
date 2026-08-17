package main

import (
	"log"
	"net/http"

	"mon-pi/internal/collector"
	"mon-pi/internal/view"
)

func main() {
	log.Println("PiMon starting — collecting initial metrics...")
	col := collector.New()
	go col.Run()

	handler := view.SetupRoutes(col)
	addr := ":8080"
	log.Printf("PiMon running at http://0.0.0.0%s", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
