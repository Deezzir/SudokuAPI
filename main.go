package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sudoku-api/pkg/sudoku"
	"sudoku-api/pkg/utils"
	"time"

	"github.com/gorilla/mux"
)

var server_port = "0.0.0.0:6969"

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

	board := sudoku.GenerateBoard(difficulty, seed)
	res, err := json.Marshal(board)
	if err != nil {
		utils.ErrorLog.Printf("Failed to encode the Board into JSON - %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode the Board into JSON\n"))
	} else {
		utils.InfoLog.Println("Sent the generated Board ('GET':'/sudoku/generate')")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func solveBoardHandler(w http.ResponseWriter, r *http.Request) {
	utils.InfoLog.Println("Received a request to solve Board ('GET':'/sudoku/solve')")

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()
	board := query.Get("board")
	if !sudoku.IsValidBoard(board) {
		utils.WarningLog.Println("Invalid Board requested: '" + board + "'")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Board requested: '" + board + "'\n"))
		return
	}

	solution := sudoku.SolveBoard(board)
	res, err := json.Marshal(solution)
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

func main() {
	r := mux.NewRouter()
	registerAPIRoutes(r)
	http.Handle("/", r)

	utils.InfoLog.Println("Starting server at: " + server_port)
	if err := http.ListenAndServe(server_port, r); err != nil {
		utils.ErrorLog.Fatalln(err)
	}
}
