package main

import (
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

	Time       int64
	GameTimer  int32
	TimerLimit int32

	Score int

	Field [][]int32

	Player PlayerStruct
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
	g.TimerLimit = 1000
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
				g.Rotate(-1)
			} else if raylib.IsKeyPressed('W') {
				g.Rotate(1)
			}

			if g.GameTimer > g.TimerLimit {
				g.Drop()
			}
		}
	} else {
		if raylib.IsKeyPressed(raylib.KeyEnter) {
			g.Init()
		}
	}
}

// Draw the game
func (g *Game) Draw() {
	raylib.BeginDrawing()

	raylib.ClearBackground(raylib.RayWhite)

	if g.Pause {
		raylib.DrawText("PAUSED", g.ScreenWidth/2-raylib.MeasureText("PAUSED", 40)/2, g.ScreenHeight/2-40, 40, raylib.Black)
	} else if g.GameOver {
		finalScore := fmt.Sprintf("Final Score: %d", g.Score)
		raylib.DrawText("GAME OVER", g.ScreenWidth/2-raylib.MeasureText("GAME OVER", 38)/2, g.ScreenHeight/2-40, 38, raylib.Black)
		raylib.DrawText(finalScore, g.ScreenWidth/2-raylib.MeasureText(finalScore, 26)/2, g.ScreenHeight/2, 26, raylib.Black)
		raylib.DrawText("Press [ENTER] to restart",
			g.ScreenWidth/2-raylib.MeasureText("Press [ENTER] to restart", 16)/2,
			g.ScreenHeight/2+30, 16, raylib.Black)
	} else {
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

		scoreText := fmt.Sprintf("Score: %d", g.Score)
		raylib.DrawText(scoreText, 5, 5, 14, raylib.White)
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

func (g *Game) Rotate(dir int) {
	offset := 1
	g.Player.RotateBlock(dir)
	for t := 0; t < 10; t++ {
		if g.CheckCollissions() {
			g.Player.Position.X += offset
			if offset > 0 {
				offset = -(offset + 1)
			} else {
				offset = -(offset - 1)
			}
		}
	}
}

func (g *Game) CheckCollissions() bool {
	for y := range g.Player.Shape {
		for x := range g.Player.Shape[0] {
			if g.Player.Shape[y][x] != 0 {
				if x+g.Player.Position.X >= len(g.Field[0]) || x+g.Player.Position.X < 0 {
					return true
				} else if y+g.Player.Position.Y == gameHeight {
					return true
				} else if g.Field[y+g.Player.Position.Y][x+g.Player.Position.X] > 0 {
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

	g.CheckForCompletes()

	g.Player = PlayerStruct{
		Position: Point{4, 0},
		Shape:    generateBlock(),
	}

	if g.CheckCollissions() {
		g.GameOver = true
	}
}

func (g *Game) CheckForCompletes() {
	lineCount := 0
	for y := range g.Field {
		isFull := true
		for x := range g.Field[0] {
			if g.Field[y][x] == 0 {
				isFull = false
				continue
			}
		}

		if isFull {
			g.Field = append(g.Field[:y], g.Field[y+1:]...)
			newLine := make([][]int32, 1)
			newLine[0] = make([]int32, gameWidth)
			g.Field = append(newLine, g.Field...)
			lineCount++
		}
	}

	g.AddScore(lineCount)
}

func (g *Game) AddScore(count int) {
	for i := 0; i < count; i++ {
		g.Score += 10 * (i + 1)
		g.GameTimer -= 10
	}
}
