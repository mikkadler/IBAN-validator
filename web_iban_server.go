package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"web-iban/iban"
)

const port = ":8080"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/validateIBAN", corsMiddleware(ibanAPI))

	err := iban.InitIbanData("./iban/data/")
	if err != nil {
		fmt.Println("Error initializing iban data:", err)
	}

	fmt.Println("Starting server on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"http://127.0.0.1:5500",
			"http://localhost:5500",
			// Add other allowed URLs here as needed
		}
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
				break
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// If it's a preflight request, return early with status code 204
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func ibanAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	type InputArray struct {
		Data []string `json:"data"`
	}

	var inputArray InputArray

	err := json.NewDecoder(r.Body).Decode(&inputArray)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	type result struct {
		IBAN     string `json:"iban"`
		Valid    bool   `json:"valid"`
		BankName string `json:"bank_name"`
		Error    string `json:"error"`
	}

	var resultArray []result

	for _, ibanString := range inputArray.Data {
		valid, err := iban.ValidateIBAN(ibanString)
		if err == nil {
			err = fmt.Errorf("ok") //we want to forward error to user and since it is in struct it can not be nil
		}
		result := result{IBAN: ibanString, Valid: valid, Error: err.Error()}
		bankName := iban.TryGetBankName(ibanString) // return name or "unknown"
		result.BankName = bankName
		resultArray = append(resultArray, result)
	}

	data, err := json.Marshal(resultArray)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Write(data)
}
