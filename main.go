package main

import (
	"math/rand"
	"time"

	"fmt"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	blockSize  = 20
	gameWidth  = 12
	gameHeight = 20
)

var (
	colors = [8]raylib.Color{
		raylib.Black,
		raylib.Blue,
		raylib.Pink,
		raylib.Orange,
		raylib.Purple,
		raylib.Green,
		raylib.Brown,
		raylib.RayWhite,
	}
)

// Game is the Gamestate
type Game struct {
	ScreenWidth  int32
	ScreenHeight int32

	GameOver bool
	Pause    bool

	Time      int64
	GameTimer int32

	Field [][]int32

	Player PlayerStruct
}

// PlayerStruct is the block in play
type PlayerStruct struct {
	Shape    [][]int32
	Position Point
}

type Point struct {
	X int
	Y int
}

func main() {
	game := Game{}
	game.Init()

	raylib.InitWindow(game.ScreenWidth, game.ScreenHeight, "sample game: tetris")
	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Update()
		game.Draw()
	}

	raylib.CloseWindow()
}

// Init will initialize the game
func (g *Game) Init() {
	g.ScreenWidth = gameWidth * blockSize
	g.ScreenHeight = gameHeight * blockSize

	g.GameOver = false
	g.Pause = false

	g.Field = make([][]int32, gameHeight)
	for i := range g.Field {
		g.Field[i] = make([]int32, gameWidth)
	}

	g.Player = PlayerStruct{
		Position: Point{2, 2},
		Shape:    generateBlock(),
	}

	g.Time = time.Now().UnixNano()
	g.GameTimer = -300
}

// Update the gamestate
func (g *Game) Update() {
	deltaTime := int32((time.Now().UnixNano() - g.Time) / (int64(time.Millisecond) / int64(time.Nanosecond)))
	g.Time = time.Now().UnixNano()
	if !g.GameOver {
		if raylib.IsKeyPressed('P') {
			g.Pause = !g.Pause
		}

		if !g.Pause {
			g.GameTimer += deltaTime
			if raylib.IsKeyPressed(raylib.KeyLeft) {
				g.Move(-1)
			} else if raylib.IsKeyPressed(raylib.KeyRight) {
				g.Move(1)
			} else if raylib.IsKeyPressed(raylib.KeyDown) {
				g.Drop()
			} else if raylib.IsKeyPressed('Q') {
				g.Player.Shape = generateBlock()
			} else if raylib.IsKeyPressed('W') {
				g.Player.Rotate(1)
			}

			if g.GameTimer > 1000 {
				g.Drop()
			}
		}
	}
}

// Draw the game
func (g *Game) Draw() {
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.RayWhite)

	for j := range g.Field {
		for i := range g.Field[0] {
			raylib.DrawRectangle(
				int32(i*blockSize),
				int32(j*blockSize),
				blockSize,
				blockSize,
				colors[g.Field[j][i]],
			)
		}
	}

	for j := range g.Player.Shape {
		for i := range g.Player.Shape[0] {
			if g.Player.Shape[j][i] != 0 {
				raylib.DrawRectangle(
					int32(i*blockSize+int(g.Player.Position.X*blockSize)),
					int32(j*blockSize+int(g.Player.Position.Y*blockSize)),
					blockSize,
					blockSize,
					colors[g.Player.Shape[j][i]],
				)
			}

		}
	}

	raylib.EndDrawing()
}

func (g *Game) Move(dir int) {
	g.Player.Position.X += dir
	if g.CheckCollissions() {
		g.Player.Position.X -= dir
	}
}

func (g *Game) Drop() {
	g.Player.Position.Y++
	g.GameTimer = 0
	if g.CheckCollissions() {
		g.Player.Position.Y--
		g.NextBlock()
	}
}

func (p *PlayerStruct) Rotate(dir int) {
	for y := 0; y < len(p.Shape); y++ {
		for x := 0; x < len(p.Shape); x++ {
			fmt.Printf("%d %d\t", p.Shape[y][x], p.Shape[x][y])
			p.Shape[y][x], p.Shape[x][y] = p.Shape[x][y], p.Shape[y][x]
			fmt.Printf("%d %d\n", p.Shape[y][x], p.Shape[x][y])
		}
	}
}

func (g *Game) CheckCollissions() bool {
	for y := range g.Player.Shape {
		for x := range g.Player.Shape[0] {
			if g.Player.Shape[y][x] != 0 {
				if x+g.Player.Position.X >= len(g.Field[0]) || x+g.Player.Position.X < 0 {
					fmt.Println("Wall")
					return true
				} else if y+g.Player.Position.Y == gameHeight {
					return true
				} else if g.Field[y+g.Player.Position.Y][x+g.Player.Position.X] > 0 {
					fmt.Println("Collide")
					return true
				}
			}
		}
	}

	return false
}

func (g *Game) NextBlock() {
	for y := range g.Player.Shape {
		for x := range g.Player.Shape[0] {
			if g.Player.Shape[y][x] != 0 {
				g.Field[y+g.Player.Position.Y][x+g.Player.Position.X] = g.Player.Shape[y][x]
			}
		}
	}

	g.Player = PlayerStruct{
		Position: Point{2, 2},
		Shape:    generateBlock(),
	}
}

func generateBlock() [][]int32 {
	switch i := rand.Int31n(7); i {
	// = block
	case 0:
		return [][]int32{
			[]int32{1, 1},
			[]int32{1, 1},
		}
		// T block
	case 1:
		return [][]int32{
			[]int32{0, 0, 0},
			[]int32{2, 2, 2},
			[]int32{0, 2, 0},
		}
	//L Block
	case 2:
		return [][]int32{
			[]int32{0, 3, 0},
			[]int32{0, 3, 0},
			[]int32{0, 3, 3},
		}
	//J Block
	case 3:
		return [][]int32{
			[]int32{0, 4, 0},
			[]int32{0, 4, 0},
			[]int32{4, 4, 0},
		}
	//S Block
	case 4:
		return [][]int32{
			[]int32{0, 5, 5},
			[]int32{5, 5, 0},
			[]int32{0, 0, 0},
		}
	//Z Block
	case 5:
		return [][]int32{
			[]int32{6, 6, 0},
			[]int32{0, 6, 6},
			[]int32{0, 0, 0},
		}
	//I Block
	default:
		return [][]int32{
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
		}
	}
}
