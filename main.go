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
var frameRate int64 = 1
var secondOnLastCall int64 = 0
var secondsSinceLastCall int64 = 222
var memory = 0
var render = true

// DIRECTIONS
const (
	Up    string = "up"
	Right        = "right"
	Down         = "down"
	Left         = "left"
)

var direction string = Right

// TILE CODES
const (
	NEWTILE int = 10003
)

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

	playerTile.TrailingMap = make(map[int]TrailingTile)
	playerTile.TrailingMap[1] = *firstTrailing

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

// RUNNING THE GAME
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
						direction = Up
					case "A":
						direction = Left
					case "S":
						direction = Down
					case "D":
						direction = Right
					case "P":

						addTrailing(playerTile)
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

// DRAW LOGIC FOR DISPLAYING THE ENTIRE BOARD
func draw(ops *op.Ops, gameBoard Board, playerTile *PlayerTile) {

	var startX float32 = 0
	var startY float32 = 0
	var current bool = false

	if secondsSinceLastCall == 222 {
		secondOnLastCall = time.Now().UnixMilli()
		secondsSinceLastCall = 0
	} else {
		if time.Now().UnixMilli()-secondOnLastCall >= frameRate {
			secondsSinceLastCall = 222
			render = true
		}
	}

	if render {

		var playerLastPositionX int = playerTile.PosX
		var playerLastPositionY int = playerTile.PosY

		switch direction {
		case Right:
			playerTile.PosX++
		case Down:
			playerTile.PosY++
		case Left:
			playerTile.PosX--
		case Up:
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
				trailingX = playerTile.TrailingMap[i].PosX
				trailingY = playerTile.TrailingMap[i].PosY

				playerTile.TrailingMap[i] = *tempTrail

			} else {

				if playerTile.TrailingMap[i].PosX != NEWTILE {

					// ANOTHER VARIABLE FOR STORING POSITION OF FURTHER TRAILING TILES

					tempTrail := new(TrailingTile)
					tempTrail.PosX = trailingX
					tempTrail.PosY = trailingY

					// PREPARING TEMP VARIABLES FOR NEXT ITERATION
					trailingX = playerTile.TrailingMap[i].PosX
					trailingY = playerTile.TrailingMap[i].PosY

					playerTile.TrailingMap[i] = *tempTrail

				} else {
					newInsertedTrail := new(TrailingTile)
					newInsertedTrail.PosX = trailingX
					newInsertedTrail.PosY = trailingY
					playerTile.TrailingMap[i] = *newInsertedTrail

				}

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

// DRAWING SPECIFIC TILE BASED OF CORDS
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

// ADDING NEW TILE TO THE PLAYER TILE MAP
func addTrailing(playerTile *PlayerTile) {

	lastTrailingKey := len(playerTile.TrailingMap)
	newTrailingKey := lastTrailingKey + 1
	newTrailingEntity := new(TrailingTile)
	newTrailingEntity.PosX = NEWTILE
	newTrailingEntity.PosY = NEWTILE
	playerTile.TrailingMap[newTrailingKey] = *newTrailingEntity

}
