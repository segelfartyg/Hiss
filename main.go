package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Board struct {
	Rows    int
	Columns int
}

type Tile struct {
	Color int16
	PosX  int
	PosY  int
}

type CurrentTile struct {
	Color int16
	PosX  int
	PosY  int
}

var gameBoard Board

var frameRate int = 1
var sinceLastFrameUpdate int = 0
var secondsSinceLastCall int = 0

func main() {

	gameBoard := new(Board)
	gameBoard.Rows = 50
	gameBoard.Columns = 50

	currentTile := new(CurrentTile)
	currentTile.PosX = 20
	currentTile.PosY = 20

	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(500), unit.Dp(500)),
			//app.MaxSize(unit.Dp(200), unit.Dp(200)),
			app.MinSize(unit.Dp(500), unit.Dp(500)),
		)

		err := run(w, *gameBoard, currentTile)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	}()

	app.Main()

}

func run(w *app.Window, gameBoard Board, currentTile *CurrentTile) error {
	var ops op.Ops
	for {
		e := <-w.Events()

		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:

			ops.Reset()
			draw(&ops, gameBoard, currentTile, &e)

			redraw(&ops)
			e.Frame(&ops)

		}
	}
}

func draw(ops *op.Ops, gameBoard Board, currentTile *CurrentTile, e *system.FrameEvent) {
	//var tileSize float32 = 5
	var startX float32 = 0
	var startY float32 = 0
	var current bool = false

	currentTile.PosX++
	//fmt.Println(currentTile.PosX)

	for i := 0; i < gameBoard.Rows; i++ {

		for j := 0; j < gameBoard.Columns; j++ {

			if currentTile.PosX == j && currentTile.PosY == i {
				current = true

			}

			drawTile(ops, startX, startY, current)
			current = false
			startX += 10
		}
		startX = 0
		startY += 10

	}

}

func redraw(ops *op.Ops) {
	op.InvalidateOp{}.Add(ops)
}

func drawTile(ops *op.Ops, xPos float32, yPos float32, current bool) {

	const r = 10
	rect := image.Rect(int(xPos), int(yPos), int(xPos+10), int(yPos+10))
	grect := clip.Rect(rect)

	randColor := rand.Intn(3)

	var red uint8
	var green uint8
	var blue uint8

	if current {
		red = 0x0
		green = 0x0
		blue = 0x0

	} else {
		switch randColor {
		case 0:
			red = 0xFF
			green = 0xFF
			blue = 0xFF
		case 1:
			red = 0xFF
			green = 0xFF
			blue = 0xFF
		case 2:
			red = 0xFF
			green = 0xFF
			blue = 0xFF
		}
	}

	paint.FillShape(ops, color.NRGBA{R: red, G: green, B: blue, A: 0xFF},
		clip.Stroke{
			Path:  grect.Path(),
			Width: 10,
		}.Op(),
	)

}
