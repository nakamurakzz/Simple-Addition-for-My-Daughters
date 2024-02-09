package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// GameState はゲームの状態を表します。
type GameState int

const (
	NotStarted GameState = iota // ゲームがまだ始まっていない状態
	Playing                     // ゲーム中の状態
	Paused                      // 一時停止中の状態
)

const (
	gameInstructions = " Addition Game\n 'q': quit 'r': reset"
	startPrompt      = " 's': start game\n"
	maxGameTime      = 30 // ゲームの最大時間 (秒)
)

// AdditionProblem は足し算の問題を表します。
type AdditionProblem struct {
	Operand1 string
	Operand2 string
	Solution string
}

var problems []AdditionProblem

func init() {
	// 足し算の問題を生成します。解答は0から9の範囲です。
	problems = make([]AdditionProblem, 0)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			sum := i + j
			if sum < 10 {
				problems = append(problems, AdditionProblem{strconv.Itoa(i), strconv.Itoa(j), strconv.Itoa(sum)})
			}
		}
	}
}

// randomProblem はランダムな足し算の問題を返します。
func randomProblem() AdditionProblem {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return problems[r.Intn(len(problems))]
}

// Game はゲームの状態を管理します。
type Game struct {
	InputChars []rune
	Text       string
	Problem    AdditionProblem
	Score      int
	State      GameState
	TimeLeft   int
	StartTime  time.Time
	Result     string
}

func (g *Game) Update() error {
	if g.State == Playing {
		g.TimeLeft = maxGameTime - int(time.Since(g.StartTime).Seconds())
		if g.TimeLeft < 0 {
			g.State = Paused
			g.TimeLeft = 0
			return nil
		}
	}

	// ユーザーの入力を取得
	g.InputChars = ebiten.AppendInputChars(g.InputChars[:0])

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("good bye")
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.State = Playing
		g.TimeLeft = maxGameTime
		g.Score = 0
		g.StartTime = time.Now()
		g.Problem = randomProblem()
	}

	if g.State != Playing {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Score = 0
		g.StartTime = time.Now()
		g.Problem = randomProblem()
	}

	if len(g.InputChars) > 0 && g.InputChars[0] >= '0' && g.InputChars[0] <= '9' {
		if string(g.InputChars) == g.Problem.Solution {
			g.Result = ""
			g.Score++
			g.Problem = randomProblem()
			return nil
		}
		g.Result = "x"
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	displayText := gameInstructions
	switch g.State {
	case NotStarted:
		displayText += startPrompt
	case Playing:
		displayText += "\n\n " + g.Problem.Operand1 + " + " + g.Problem.Operand2 + " = _   " + g.Result + "\n" + "\n Score: " + strconv.Itoa(g.Score) + "  Time: " + strconv.Itoa(g.TimeLeft)
	case Paused:
		displayText += startPrompt + "\n\n\n Score: " + strconv.Itoa(g.Score) + "  Time: " + strconv.Itoa(g.TimeLeft)
	}
	ebitenutil.DebugPrint(screen, displayText)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 120
}

func main() {
	game := &Game{
		Problem:  randomProblem(),
		TimeLeft: maxGameTime,
		State:    NotStarted,
	}
	ebiten.SetWindowSize(640, 240)
	ebiten.SetWindowTitle("Addition Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
