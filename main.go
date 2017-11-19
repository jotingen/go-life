package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const (
	X     = 1024
	Y     = 768
)

func main() {
	pixelgl.Run(run)
}

func run() {
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

	var board [X][Y]bool
	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			if rand.Intn(2) == 1 {
				board[x][y] = true
			}
		}
	}

	for !win.Closed() {
		imd.Reset()
		imd.Clear()
		win.Clear(colornames.White)

		var boardNext [X][Y]bool

		for x := 0; x < X; x++ {
			for y := 0; y < Y; y++ {
				if board[x][y] {
					imd.Color = pixel.RGB(0, 0, 0)
					imd.Push(pixel.V(float64(x), float64(y)))
					imd.Push(pixel.V(float64(x+1), float64(y+1)))
					imd.Rectangle(0)
				}
				neighbors := 0
				if x > 0 {
					if y > 0 && board[x-1][y-1] {
						neighbors++
					}
					if y < Y-1 && board[x-1][y+1] {
						neighbors++
					}
					if board[x-1][y] {
						neighbors++
					}
				}

				if x < X-1 {
					if y > 0 && board[x+1][y-1] {
						neighbors++
					}
					if y < Y-1 && board[x+1][y+1] {
						neighbors++
					}
					if board[x+1][y] {
						neighbors++
					}
				}
				if y > 0 && board[x][y-1] {
					neighbors++
				}
				if y < Y-1 && board[x][y+1] {
					neighbors++
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
