package main

import (
	"fmt"
	"go-pgsql/router"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	r := router.Router();
	fmt.Println("Starting serever on port 5000... ")

	log.Fatal(http.ListenAndServe(":5000", r))
}