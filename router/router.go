package router

import (
	"go-pgsql/middleware"

	"github.com/gorilla/mux"
)

func Router () *mux.Router 	{

  router := mux.NewRouter()

  router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
  router.HandleFunc("/api/stock/", middleware.GetAllStock).Methods("GET", "OPTIONS")
  router.HandleFunc("/api/new-stock", middleware.CreateStock).Methods("POST", "OPTIONS")
  router.HandleFunc("/api/stock/{od}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
  router.HandleFunc("/api/delete-stock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

}