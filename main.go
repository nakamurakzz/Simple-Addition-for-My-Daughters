package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Calc struct {
	first  string
	second string
	answer string
}

var calcs []Calc

func init() {
	//	create a list of calculations
	//	answer: 0~9
	calcs = make([]Calc, 0)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			answer := i + j
			if answer < 10 {
				calcs = append(calcs, Calc{strconv.Itoa(i), strconv.Itoa(j), strconv.Itoa(answer)})
			}
		}
	}
}

func randCalc() Calc {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return calcs[r.Intn(len(calcs))]
}

type Game struct {
	runes []rune
	title string
	text  string
	calc  Calc
	score int
}

func (g *Game) check(a string) bool {
	if a == g.calc.answer {
		return true
	}
	return false
}

func (g *Game) Update() error {
	// ユーザの入力値の取得
	g.runes = ebiten.AppendInputChars(g.runes[:0])

	if len(g.runes) > 0 && g.runes[0] >= 48 && g.runes[0] <= 57 {
		// check g.runes[0] is a number
		if g.check(string(g.runes)) {
			g.score++
			g.calc = randCalc()
		}
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	score := strconv.Itoa(g.score)
	t := "\n " + g.calc.first + " + " + g.calc.second + " = _" + "\n" + "\n Score: " + string(score)
	ebitenutil.DebugPrint(screen, t)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	g := &Game{
		title: "Type on the keyboard:\n",
		calc:  randCalc(),
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
