package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sudoku-api/pkg/sudoku"
	"sudoku-api/pkg/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const server_port = "0.0.0.0:80"
const apiKeyHeader = "X-API-Key"

var attempts = 0

func generateBoardHandler(w http.ResponseWriter, r *http.Request) {
	utils.InfoLog.Println("Received a request to genearate board ('GET':'/sudoku/generate')")

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	difficulty := strings.ToLower(query.Get("difficulty"))
	if difficulty == "" {
		difficulty = sudoku.Difficulties[0]
	}
	if !sudoku.IsValidDifficulty(difficulty) {
		utils.WarningLog.Println("Invalid difficulty requested: '" + difficulty + "'")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid difficulty requested: '" + difficulty + "'\n"))
		return
	}

	seed := query.Get("seed")
	if seed == "" {
		seed = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if !sudoku.IsValidSeed(seed) {
		utils.WarningLog.Println("Invalid seed requested: " + seed)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid seed requested: '" + seed + "'\n"))
		return
	}

	board, ok := sudoku.GenerateBoard(difficulty, seed, attempts)
	if ok {
		res, err := json.Marshal(struct {
			Board string `json:"board"`
		}{
			Board: board,
		})

		if err != nil {
			utils.ErrorLog.Printf("Failed to encode the Board into JSON - %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to encode the Board into JSON\n"))
		} else {
			utils.InfoLog.Println("Sent the generated Board ('GET':'/sudoku/generate')")
			w.WriteHeader(http.StatusOK)
			w.Write(res)
		}
	} else {
		utils.WarningLog.Println("Failed to generate Board with difficulty: '" + difficulty + "' and seed: '" + seed + "'")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to generate Board with difficulty: '" + difficulty + "' and seed: '" + seed + "'\n"))
	}
}

func solveBoardHandler(w http.ResponseWriter, r *http.Request) {
	utils.InfoLog.Println("Received a request to solve Board ('GET':'/sudoku/solve')")

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	board := query.Get("board")
	if board == "" {
		utils.WarningLog.Println("No Board provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No Board provided\n"))
		return
	}
	if !sudoku.IsValidBoard(board) {
		utils.WarningLog.Println("Invalid Board provided: '" + board + "'")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Board provided: '" + board + "'\n"))
		return
	}

	solution, counter := sudoku.SolveBoard(board)
	result := sudoku.CounterToMsg(counter)

	res, err := json.Marshal(struct {
		Solution string `json:"solution"`
		Result   string `json:"result"`
	}{
		Solution: solution,
		Result:   result,
	})

	if err != nil {
		utils.ErrorLog.Printf("Failed to encode the Board Solution into JSON - %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode the Board Solution into JSON\n"))
	} else {
		utils.InfoLog.Println("Sent the solved Board ('GET':'/sudoku/solve')")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WarningLog.Println("Reqested non-existent route")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Route not found, please check your URL\n"))

}

var registerAPIRoutes = func(router *mux.Router) {
	router.HandleFunc("/sudoku/generate", generateBoardHandler).Methods("GET")
	router.HandleFunc("/sudoku/solve", solveBoardHandler).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
}

func apiKeyMiddleware(apiKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get(apiKeyHeader)
		if key != apiKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			utils.WarningLog.Println("Received Invalid API key")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getConfsFromEnv() (apiKey string, port string, totalAttempts int) {
	if value, ok := os.LookupEnv("API_KEY"); !ok {
		panic("API_KEY environment variable not set, crashing")
	} else {
		apiKey = value
	}

	if value, ok := os.LookupEnv("PORT"); !ok {
		port = server_port
	} else {
		port = value
	}

	totalAttempts = sudoku.TotalAttempts
	if value, ok := os.LookupEnv("TOTAL_ATTEMPTS"); ok {
		value, err := strconv.Atoi(value)
		if err == nil {
			totalAttempts = value
		}
	}

	return apiKey, port, totalAttempts
}

func main() {
	if err := godotenv.Load(); err != nil {
		utils.WarningLog.Println("Error loading .env file")
	}
	apiKey, port, totalAttempts := getConfsFromEnv()

	utils.InfoLog.Println("Using total attempts: " + strconv.Itoa(totalAttempts))
	utils.InfoLog.Println("Using API key: " + apiKey)
	attempts = totalAttempts

	r := mux.NewRouter()
	registerAPIRoutes(r)
	securityMiddleware := apiKeyMiddleware(apiKey, r)

	utils.InfoLog.Println("Starting server at: " + port)
	if err := http.ListenAndServe(port, securityMiddleware); err != nil {
		utils.ErrorLog.Fatalln(err)
	}
}
