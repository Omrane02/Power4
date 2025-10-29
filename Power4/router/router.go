package router

import (
	"net/http"
	"power4/controller"
)

func SetupRoutes() {
	http.HandleFunc("/", controller.GameHandler)
	http.HandleFunc("/play", controller.PlayHandler)
	http.HandleFunc("/reset", controller.ResetHandler)
	http.HandleFunc("/scoreboard", controller.ScoreboardHandler)
	http.HandleFunc("/about", controller.AboutHandler)
	http.HandleFunc("/contact", controller.ContactHandler)
}
