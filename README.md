# SudokuAPI

## Description

This is a simple API that generates a sudoku puzzle and solves it. It is written in GoLang.

## Usage

```bash
go build main.go
./main
```

## Endpoints

1. GET /sudoku/generate
    - Generates a sudoku puzzle
    - Returns a string of 81 integers representing the puzzle
    - Optional query parameter: difficulty
        - 1: Easy
        - 2: Medium
        - 3: Hard
    - Optional query parameter: seed
        - Integer value to seed the random number generator as string

2. GET /sudoku/solve
    - Solves a sudoku puzzle
    - Returns a string of 81 integers representing the solved puzzle
    - Required query parameter: puzzle
        - string of 81 integers representing the puzzle
