package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	X = 1400
	Y = 800
)

func main() {
	pixelgl.Run(run)
}

func run() {
	var wg sync.WaitGroup

	cfg := pixelgl.WindowConfig{
		Title:  "Life",
		Bounds: pixel.R(0, 0, X, Y),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	imd := imdraw.New(nil)

	timestamp := time.Now().UTC()

	board := make([][]bool, X)
	for x := 0; x < X; x++ {
		board[x] = make([]bool, Y)
		for y := 0; y < Y; y++ {
			board[x][y] = rand.Intn(2) == 0
		}
	}

	for !win.Closed() {
		imd.Reset()
		imd.Clear()
		win.Clear(colornames.White)

		boardNext := make([][]bool, X)
		for x := 0; x < X; x++ {
			boardNext[x] = make([]bool, Y)
		}

		wg.Add(16)
		go life(0,     X/4,   0,   Y/4,   board, boardNext, &wg)
		go life(X/4,   X/2,   0,   Y/4,   board, boardNext, &wg)
		go life(X/2,   3*X/4, 0,   Y/4,   board, boardNext, &wg)
		go life(3*X/4, X,     0,   Y/4,   board, boardNext, &wg)
		go life(0,     X/4,   Y/4, Y/2,   board, boardNext, &wg)
		go life(X/4,   X/2,   Y/4, Y/2,   board, boardNext, &wg)
		go life(X/2,   3*X/4, Y/4, Y/2,   board, boardNext, &wg)
		go life(3*X/4, X,     Y/4, Y/2,   board, boardNext, &wg)
		go life(0,     X/4,   Y/2, 3*Y/4, board, boardNext, &wg)
		go life(X/4,   X/2,   Y/2, 3*Y/4, board, boardNext, &wg)
		go life(X/2,   3*X/4, Y/2, 3*Y/4, board, boardNext, &wg)
		go life(3*X/4, X,     Y/2, 3*Y/4, board, boardNext, &wg)
		go life(0,     X/4,   3*Y/4, Y,   board, boardNext, &wg)
		go life(X/4,   X/2,   3*Y/4, Y,   board, boardNext, &wg)
		go life(X/2,   3*X/4, 3*Y/4, Y,   board, boardNext, &wg)
		go life(3*X/4, X,     3*Y/4, Y,   board, boardNext, &wg)

		for x := 0; x < X; x++ {
			for y := 0; y < Y; y++ {
				if board[x][y] {
					imd.Color = pixel.RGB(0, 0, 0)
					imd.Push(pixel.V(float64(x), float64(y)))
					imd.Push(pixel.V(float64(x+1), float64(y+1)))
					imd.Rectangle(0)
				}

			}
		}
		wg.Wait()
		board = boardNext
		imd.Draw(win)

		basicTxt := text.New(pixel.V(100, 500), basicAtlas)
		basicTxt.Color = colornames.Red
		fmt.Fprintf(basicTxt, "%4.1f", 1.0/time.Since(timestamp).Seconds())
		timestamp = time.Now().UTC()
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 3))

		win.Update()
	}
}

func life(xlo int, xhi int, ylo int, yhi int, board [][]bool, boardNext [][]bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := xlo; x < xhi; x++ {
		for y := ylo; y < yhi; y++ {
			neighbors := 0
			for nx := x - 1; nx <= x+1; nx++ {
				for ny := y - 1; ny <= y+1; ny++ {
					x0 := nx
					y0 := ny
					if x0 == -1 {
						x0 = X - 1
					}
					if y0 == -1 {
						y0 = Y - 1
					}
					if x0 == X {
						x0 = 0
					}
					if y0 == Y {
						y0 = 0
					}
					//fmt.Println(x,nx,x0,y,ny,y0)
					if !(x0 == x && y0 == y) && board[x0][y0] {
						neighbors++
					}
				}
			}
			if board[x][y] {
				if neighbors == 2 || neighbors == 3 {
					boardNext[x][y] = true
				}
			} else {
				if neighbors == 3 {
					boardNext[x][y] = true
				}
			}
		}
	}
}
