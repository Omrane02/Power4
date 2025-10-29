package main

import (
	"fmt"
	"net/http"
	"power4/router"
)

func main() {
	router.SetupRoutes()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
