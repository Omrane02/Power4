package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	rows = 6
	cols = 7
)

type Game struct {
	Board       [][]string
	Current     string
	Winner      string
	GameOver    bool
	MoveMessage string
	IsDraw      bool
}

var game Game

func init() {
	game = loadOrNew()
}

func newGame() Game {
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
	}
	return Game{Board: board, Current: "red"}
}

func ensureData() {
	os.MkdirAll("data", 0755)
}

func saveGame() {
	ensureData()
	f, _ := os.Create(filepath.Join("data", "save.json"))
	defer f.Close()
	json.NewEncoder(f).Encode(game)
}

func loadOrNew() Game {
	ensureData()
	f, err := os.Open(filepath.Join("data", "save.json"))
	if err != nil {
		return newGame()
	}
	defer f.Close()
	var g Game
	if json.NewDecoder(f).Decode(&g) != nil || g.Board == nil || len(g.Board) != rows {
		return newGame()
	}
	return g
}

func loadScores() map[string]int {
	ensureData()
	file := filepath.Join("data", "scores.json")
	f, err := os.Open(file)
	if err != nil {
		return map[string]int{"red": 0, "yellow": 0}
	}
	defer f.Close()
	var scores map[string]int
	if json.NewDecoder(f).Decode(&scores) != nil {
		return map[string]int{"red": 0, "yellow": 0}
	}
	if _, ok := scores["red"]; !ok {
		scores["red"] = 0
	}
	if _, ok := scores["yellow"]; !ok {
		scores["yellow"] = 0
	}
	return scores
}

func saveScores(scores map[string]int) {
	ensureData()
	f, _ := os.Create(filepath.Join("data", "scores.json"))
	defer f.Close()
	json.NewEncoder(f).Encode(scores)
}

func dropDisc(col int) {
	if game.GameOver {
		return
	}
	for r := rows - 1; r >= 0; r-- {
		if game.Board[r][col] == "" {
			game.Board[r][col] = game.Current
			checkWinner(r, col)
			if !game.GameOver {
				if game.Current == "red" {
					game.Current = "yellow"
				} else {
					game.Current = "red"
				}
			}
			saveGame()
			return
		}
	}
	game.MoveMessage = "Colonne pleine, choisis-en une autre."
}

func countDir(r, c, dr, dc int, player string) int {
	count := 0
	for {
		r += dr
		c += dc
		if r < 0 || r >= rows || c < 0 || c >= cols {
			break
		}
		if game.Board[r][c] != player {
			break
		}
		count++
	}
	return count
}

func checkWinner(r, c int) {
	player := game.Board[r][c]
	if player == "" {
		return
	}
	directions := [][2]int{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
	}
	for _, d := range directions {
		count := 1
		count += countDir(r, c, d[0], d[1], player)
		count += countDir(r, c, -d[0], -d[1], player)
		if count >= 4 {
			game.Winner = player
			game.GameOver = true
			scores := loadScores()
			scores[player]++
			saveScores(scores)
			os.Remove(filepath.Join("data", "save.json"))
			return
		}
	}
	full := true
	for _, row := range game.Board {
		for _, cell := range row {
			if cell == "" {
				full = false
				break
			}
		}
	}
	if full {
		game.IsDraw = true
		game.GameOver = true
	}
}

func resetGame() {
	game = newGame()
	saveGame()
}

func GameHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, game)
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	colStr := r.FormValue("col")
	col, err := strconv.Atoi(colStr)
	if err == nil && col >= 0 && col < cols {
		dropDisc(col)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
	resetGame()
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	scores := loadScores()
	tmpl := template.Must(template.ParseFiles("template/scoreboard.html"))
	tmpl.Execute(w, scores)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/about.html"))
	tmpl.Execute(w, nil)
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/contact.html"))
	tmpl.Execute(w, nil)
}
