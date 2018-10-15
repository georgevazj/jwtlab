package main

import (
	"github.com/gorilla/mux"
	"github.com/georgevazj/jwtlab/app"
	"os"
	"fmt"
	"net/http"
	"github.com/georgevazj/jwtlab/controllers"
)

func main()  {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication) // Attach JWT auth middleware

	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET") //  user/2/contacts

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
