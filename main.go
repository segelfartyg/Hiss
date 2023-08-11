package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
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

type PlayerTile struct {
	Color       int16
	PosX        int
	PosY        int
	TrailingMap map[int]TrailingTile
}

type TrailingTile struct {
	Color int16
	PosX  int
	PosY  int
	NextX int
	NextY int
}

var gameBoard Board
var frameRate int64 = 10
var secondOnLastCall int64 = 0
var secondsSinceLastCall int64 = 222
var memory = 0
var render = true
var direction = "right"

var tailLength = 3

func main() {

	gameBoard := new(Board)
	gameBoard.Rows = 50
	gameBoard.Columns = 50

	playerTile := new(PlayerTile)
	playerTile.PosX = 20
	playerTile.PosY = 20

	firstTrailing := new(TrailingTile)
	firstTrailing.PosX = 19
	firstTrailing.PosY = 20
	firstTrailing.NextX = 20
	firstTrailing.NextY = 20

	secondTrailing := new(TrailingTile)
	secondTrailing.PosX = 18
	secondTrailing.PosY = 20
	secondTrailing.NextX = 19
	secondTrailing.NextY = 20

	playerTile.TrailingMap = make(map[int]TrailingTile)

	playerTile.TrailingMap[1] = *firstTrailing
	playerTile.TrailingMap[2] = *secondTrailing

	go func() {
		w := app.NewWindow(
			app.Size(unit.Dp(500), unit.Dp(500)),
			//app.MaxSize(unit.Dp(200), unit.Dp(200)),
			app.MinSize(unit.Dp(500), unit.Dp(500)),
		)

		err := run(w, *gameBoard, playerTile)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)

	}()

	app.Main()

}

func run(w *app.Window, gameBoard Board, playerTile *PlayerTile) error {
	var ops op.Ops

	for windowEvent := range w.Events() {
		switch winE := windowEvent.(type) {

		case system.DestroyEvent:
			return winE.Err
		case system.FrameEvent:

			gtx := layout.NewContext(&ops, winE)

			for _, gtxEvent := range gtx.Events(0) {

				switch gtxE := gtxEvent.(type) {

				case key.Event:
					switch gtxE.Name {
					case "W":
						direction = "up"
					case "A":
						direction = "left"
					case "S":
						direction = "down"
					case "D":
						direction = "right"
					}
				}

			}

			eventArea := clip.Rect(
				image.Rectangle{
					Min: image.Point{0, 0},
					Max: image.Point{gtx.Constraints.Max.X, gtx.Constraints.Max.Y},
				},
			).Push(gtx.Ops)

			key.InputOp{
				Keys: key.Set("(Shift)-F|(Shift)-S|(Shift)-U|(Shift)-D|(Shift)-J|(Shift)-K|(Shift)-W|(Shift)-N|Space"),
				Tag:  0,
			}.Add(gtx.Ops)

			eventArea.Pop()

			draw(gtx.Ops, gameBoard, playerTile)
			winE.Frame(gtx.Ops)
		}

	}
	return nil
}

func draw(ops *op.Ops, gameBoard Board, playerTile *PlayerTile) {

	var startX float32 = 0
	var startY float32 = 0
	var current bool = false

	if secondsSinceLastCall == 222 {
		secondOnLastCall = time.Now().UnixMilli()

		secondsSinceLastCall = 0
		//fmt.Println(secondOnLastCall)

	} else {
		if time.Now().UnixMilli()-secondOnLastCall >= frameRate {

			secondsSinceLastCall = 222
			render = true
			// if secondOnLastCall >= 2 {
			// 	render = true

			// }

		}
	}

	if render {

		var playerLastPositionX int = playerTile.PosX
		var playerLastPositionY int = playerTile.PosY

		switch direction {
		case "right":
			playerLastPositionX = playerTile.PosX
			playerTile.PosX++
		case "down":
			playerLastPositionY = playerTile.PosY
			playerTile.PosY++
		case "left":
			playerLastPositionX = playerTile.PosX
			playerTile.PosX--
		case "up":
			playerLastPositionY = playerTile.PosY
			playerTile.PosY--
		}

		// SAVING TRAILING TILES POSITION IN TEMP VARIABLES
		var trailingX int = playerTile.TrailingMap[1].PosX
		var trailingY int = playerTile.TrailingMap[1].PosY

		for i := 1; i <= len(playerTile.TrailingMap); i++ {
			if i == 1 {
				// CREATING NEW TILES WITH THE VALUE OF PLAYER TILE'S LAST POSITION
				tempTrail := new(TrailingTile)
				tempTrail.PosX = playerLastPositionX
				tempTrail.PosY = playerLastPositionY
				playerTile.TrailingMap[i] = *tempTrail

			} else {
				// ANOTHER VARIABLE FOR STORING POSITION OF FURTHER TRAILING TILES
				tempTrailingX := trailingX
				tempTrailingY := trailingY

				tempTrail := new(TrailingTile)
				tempTrail.PosX = tempTrailingX
				tempTrail.PosY = tempTrailingY
				playerTile.TrailingMap[i] = *tempTrail

				// PREPARING TEMP VARIABLES FOR NEXT ITERATION
				trailingX = tempTrail.PosX
				trailingY = tempTrail.PosY
			}

		}

		for i := 0; i < gameBoard.Rows; i++ {

			for j := 0; j < gameBoard.Columns; j++ {

				if playerTile.PosX == j && playerTile.PosY == i {
					current = true
				}

				for _, v := range playerTile.TrailingMap {
					if v.PosX == j && v.PosY == i {
						current = true
					}
				}

				drawTile(ops, startX, startY, current)
				current = false
				startX += 10
			}
			startX = 0
			startY += 10

		}
		render = false

		op.InvalidateOp{}.Add(ops)

	} else {
		for i := 0; i < gameBoard.Rows; i++ {
			for j := 0; j < gameBoard.Columns; j++ {
				if playerTile.PosX == j && playerTile.PosY == i {
					current = true
				}
				for _, v := range playerTile.TrailingMap {
					if v.PosX == j && v.PosY == i {
						current = true
					}
				}
				drawTile(ops, startX, startY, current)
				current = false
				startX += 10
			}
			startX = 0
			startY += 10
		}
		op.InvalidateOp{}.Add(ops)
	}

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

//func drawTrailing
