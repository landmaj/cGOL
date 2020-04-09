package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
	"math/rand"
	"time"
)

const width = 640
const height = 480
const TPS = 60

type Game struct {
	board  []bool
	width  int
	height int
}

func NewGame(w, h int) Game {
	g := Game{
		board:  make([]bool, w*h),
		width:  w,
		height: h,
	}
	var src = rand.NewSource(time.Now().UnixNano())
	var r = rand.New(src)
	for i := range g.board {
		g.board[i] = r.Intn(2) == 0
	}
	return g
}

func (g *Game) LivingNeighbors(index int) (livingNeighbors int) {
	// safe directions
	up := index >= g.width
	down := index/g.width < g.height-1
	left := index%g.width > 0
	right := index%g.width < g.width-1
	var neighbors []bool
	if up {
		neighbors = append(neighbors, g.board[index-g.width])
	}
	if right {
		neighbors = append(neighbors, g.board[index+1])
	}
	if down {
		neighbors = append(neighbors, g.board[index+g.width])
	}
	if left {
		neighbors = append(neighbors, g.board[index-1])
	}
	if up && right {
		neighbors = append(neighbors, g.board[index-g.width+1])
	}
	if right && down {
		neighbors = append(neighbors, g.board[index+g.width+1])
	}
	if down && left {
		neighbors = append(neighbors, g.board[index+g.width-1])
	}
	if left && up {
		neighbors = append(neighbors, g.board[index-g.width-1])
	}
	for _, n := range neighbors {
		if n {
			livingNeighbors++
		}
	}
	return
}

func (g *Game) Evolve() {
	newBoard := make([]bool, len(g.board))
	for index, cell := range g.board {
		ln := g.LivingNeighbors(index)
		if cell && 1 < ln && ln < 4 {
			// Any live cell with two or three neighbors survives.
			newBoard[index] = true
		} else if !cell && ln == 3 {
			// Any dead cell with three live neighbors becomes a live cell.
			newBoard[index] = true
		} else {
			newBoard[index] = false
		}
	}
	g.board = newBoard
}

func (g *Game) Draw(img *ebiten.Image) {
	for i, pxl := range g.board {
		x, y := i%g.width, i/g.width
		if pxl {
			img.Set(x, y, color.White)
		} else {
			img.Set(x, y, color.Black)
		}
	}
}

func (g *Game) Update(img *ebiten.Image) error {
	g.Evolve()
	if ebiten.IsDrawingSkipped() {
		return nil
	}
	g.Draw(img)
	return nil
}

func main() {
	g := NewGame(width, height)
	ebiten.SetMaxTPS(TPS)
	if err := ebiten.Run(g.Update, width, height, 2, "cGOL"); err != nil {
		log.Fatal(err)
	}
}
