package main

import "errors"

type CellValue string

const (
	Empty CellValue = ""
	X     CellValue = "X"
	O     CellValue = "O"
)

type GameStatus string

const (
	StatusWaiting GameStatus = "waiting"
	StatusPlaying GameStatus = "playing"
	StatusWon     GameStatus = "won"
	StatusDraw    GameStatus = "draw"
)

type Board [3][3]CellValue

type Game struct {
	Board        Board
	Turn         CellValue
	Status       GameStatus
	Winner       CellValue
	WinningCells [][2]int
}

func NewGame() *Game {
	return &Game{
		Turn:   X,
		Status: StatusWaiting,
	}
}

func (g *Game) ApplyMove(player CellValue, row, col int) error {
	if g.Status != StatusPlaying {
		return errors.New("game is not in progress")
	}
	if player != g.Turn {
		return errors.New("it is not your turn")
	}
	if row < 0 || row > 2 || col < 0 || col > 2 {
		return errors.New("invalid cell position")
	}
	if g.Board[row][col] != Empty {
		return errors.New("cell is already taken")
	}

	g.Board[row][col] = player

	if won, cells := g.checkWin(); won {
		g.Status = StatusWon
		g.Winner = player
		g.WinningCells = cells
		return nil
	}

	if g.checkDraw() {
		g.Status = StatusDraw
		return nil
	}

	if g.Turn == X {
		g.Turn = O
	} else {
		g.Turn = X
	}
	return nil
}

func (g *Game) checkWin() (bool, [][2]int) {
	b := g.Board
	lines := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},
		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}
	for _, line := range lines {
		a, bc, c := line[0], line[1], line[2]
		if b[a[0]][a[1]] != Empty &&
			b[a[0]][a[1]] == b[bc[0]][bc[1]] &&
			b[bc[0]][bc[1]] == b[c[0]][c[1]] {
			return true, [][2]int{a, bc, c}
		}
	}
	return false, nil
}

func (g *Game) checkDraw() bool {
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == Empty {
				return false
			}
		}
	}
	return true
}

// Reset clears the board and starts a new round with firstTurn going first.
func (g *Game) Reset(firstTurn CellValue) {
	g.Board = Board{}
	g.Turn = firstTurn
	g.Winner = Empty
	g.WinningCells = nil
	g.Status = StatusPlaying
}
