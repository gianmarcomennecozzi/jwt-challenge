package main

import (
	"log"
	"net/http"
	chal "jwt-challenge"
)

func main()  {
	myApp := chal.NewApp("menne")
	log.Fatal(http.ListenAndServe(":8080", myApp.Handler()))
}