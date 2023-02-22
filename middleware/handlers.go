package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-pgsql/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID int64 `json:"id,omitempty"`
  Message string `json:"message,omitempty"`
} 

func createConnectio() *sql.DB {
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatal("Error Loading .env File")
  }

  db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

  if err != nil {
    panic(err)
  }

  err  = db.Ping()
  if err != nil {
    panic(err)
  }

  fmt.Println("Successfully connected to pg..")
  return db
}


func CreateStock(w http.ResponseWriter, r *http.Request) {
  var stock models.Stock

  err := json.NewDecoder(r.Body).Decode(&stock)
  if err != nil {
    log.Fatalf("Unable to decode request body. %v", err)
  }

  insertID := insertStokc(stock)

  res := response{
    ID: insertID,
    Message: "Stock created successfully.",
  }

  json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r);

  id, err := strconv.Atoi(params["id"])

  if err != nil {
    log.Fatalf("Unable to convert string into int.., %v", err)
  }

  stock, err := getStock(int64(id))

  if err != nil {
    log.Fatalf("Unable to get stock from db ..., %v", err)
  }

  json.NewEncoder(w).Encode(stock)

}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
  stocks, err := getAllStocks()
  if err != nil {
    log.Fatalf("Unable to get all stocks.., %v", err)
  }

  json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  
  id, err := strconv.Atoi(params["id"])

  if err != nil {
    log.Fatalf("Unable to convert string into int.., %v", err)
  }

  var stock models.Stock
  err  = json.NewDecoder(r.Body).Decode(&stock)

  if err != nil {
    log.Fatalf("Unable to decode request body.., %v", err)
  }
  updatedRows := updateStock(int64(id), stock)
  
  msg := fmt.Sprintf("Stock updated successfully, Total rows affected %v", updatedRows)
  res := response {
    ID: int64(id),
    Message: msg,
  }
  json.NewEncoder(w).Encode(res)
  
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r)
  id, err := strconv.ParseInt(params["id"], 10, 64)
  if err != nil {
    log.Fatalf("Unable to convert string into int.., %v", err)
  }

  deletedRows := deletStock(int64(id))
  msg := fmt.Sprintf("Stock deleted successfully, Total rows affected %v", deletedRows)
  res := response {
    ID: id,
    Message: msg,
  }
  json.NewEncoder(w).Encode(res) 
}




func insertStokc(stock models.Stock) int64 {
  db := createConnectio()
  defer db.Close()
  sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $2) RETURNING stock_id`
  var id int64

  err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

  if err != nil {
    log.Fatalf("Unable to insert new row in db.. %v", err)
  }
  fmt.Printf("Inserted single row, %v", id)
  return id
}

func getStock(id int64) (models.Stock, error) {
  db := createConnectio()
  defer db.Close()
  var stock models.Stock
  sqlStatement := `SELECT * FROM stocks WHERE stocl_id = $1`
  row := db.QueryRow(sqlStatement, id)
  err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

  switch err {
  case sql.ErrNoRows:
    fmt.Println("No data found!..")
    return stock, nil
  case nil:
    return stock, nil
  default:
    log.Fatalf("Unable to retrive data from db... %v", err)
  }

  return stock, err

}

func getAllStocks() ([]models.Stock, error) {
  db := createConnectio()
  defer db.Close()
  var stocks[] models.Stock
  sqlStatement := `SELECT * FROM stocks`

  rows, err := db.Query(sqlStatement)

  if err != nil {
    log.Fatalf("Unable to execute query!... %v", err)
  }

  for rows.Next() {
    var stock models.Stock
    err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
    if err != nil {
      log.Fatalf("Unable to scan object!... %v", err)
      stocks = append(stocks, stock)
    }
  }

  return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {
  db := createConnectio()
  defer db.Close()
  sqlStatement := `UPDATE stocks set name=$!2, price=$3, company=$4 WHERE stock_id = $1`

  res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
  if err != nil {
    log.Fatalf("Unable to update row.. %v", err)
  }

  rowsAffected, err := res.RowsAffected()
  if err != nil {
    log.Fatalf("Unable to update row.. %v", err)
  }

  fmt.Println("Row updated")
  return rowsAffected
}

func deletStock(id int64) int64 {
  db := createConnectio()
  defer db.Close()
  sqlStatement := `DELETE FROM stocks WHERE stock_id = $1`

  res, err := db.Exec(sqlStatement, id)
  if err != nil {
    log.Fatalf("Unable to delete row.. %v", err)
  }

  rowsAffected, err := res.RowsAffected()
  if err != nil {
    log.Fatalf("Unable to delete row.. %v", err)
  }

  fmt.Println("Row deleted")
  return rowsAffected

}
