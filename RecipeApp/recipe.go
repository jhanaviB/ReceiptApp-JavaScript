package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
    ShortDescription string  `json:"shortDescription"`
    Price             string `json:"price"`
}

type Receipt struct {
    Retailer    string  `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items       []Item  `json:"items"`
    Total       string  `json:"total"`
}

var receiptStore = make(map[string]Receipt)
var mu sync.Mutex

func getPoints(receipt Receipt) int {
    var points float64 = 0
    regex := regexp.MustCompile(`(?i)[a-z0-9]`)

    for _, char := range receipt.Retailer {
        if regex.MatchString(string(char)) {
            points += 1
        }
    }

    total, err := strconv.ParseFloat(receipt.Total, 64)
    if err != nil {
        return 0
    }

    if total == float64(int(total)) {
        points += 50
    }

    if total != 0 && int(total*100)%25 == 0 {
        points += 25
    }

    points += float64(int(len(receipt.Items) / 2) * 5)

    for _, item := range receipt.Items {
        if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
            price, err := strconv.ParseFloat(item.Price, 64)
            if err != nil {
                continue
            }
            points += float64(math.Ceil(price * 0.2))
        }
    }
    layout := "2006-01-02"
    parsedTime, err := time.Parse(layout, receipt.PurchaseDate)
    if parsedTime.Weekday()%2 != 0 {
        points += 6
    }

    timeParts := strings.Split(receipt.PurchaseTime, ":")
    if len(timeParts) > 0 {
        hour, err := strconv.Atoi(timeParts[0])
        if err == nil && hour >= 14 && hour < 16 {
            points += 10
        }
    }

    return int(points)
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
    var receipt Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, "Invalid receipt data", http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    id := uuid.New().String()
    receiptStore[id] = receipt

    response := map[string]string{"message": "Receipt created!", "id": id}
    json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    mu.Lock()
    receipt, exists := receiptStore[id]
    mu.Unlock()

    if !exists {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    points := getPoints(receipt)
    json.NewEncoder(w).Encode(map[string]int{"points": points})
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/receipts/process", processReceipt).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

    http.Handle("/", r)
    fmt.Println("Server running at http://localhost:8081")
    http.ListenAndServe(":8081", nil)
}
