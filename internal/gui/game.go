package gui

import (
	"gopher8/internal/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	CPU *core.CPU
}

func (g *Game) Update() error {

	// user input at this frame
	handleInput(g)

	// run ~600 cycles
	for i := 0; i < 10; i++ {
		g.CPU.Cycle()
	}

	// reduce delay timer by 60Hz
	g.CPU.ReduceDelay()
	return nil
}

var keyMap = map[ebiten.Key]ebiten.Key{
	ebiten.Key1: 0x1, ebiten.Key2: 0x2, ebiten.Key3: 0x3, ebiten.Key4: 0xC,
	ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0xD,
	ebiten.KeyA: 0x7, ebiten.KeyS: 0x8, ebiten.KeyD: 0x9, ebiten.KeyF: 0xE,
	ebiten.KeyZ: 0xA, ebiten.KeyX: 0x0, ebiten.KeyC: 0xB, ebiten.KeyV: 0xF,
}

func handleInput(g *Game) {
	for ebitenKey, chip8Key := range keyMap {
		g.CPU.Keypad[chip8Key] = ebiten.IsKeyPressed(ebitenKey)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for r := range core.Height {
		for c := range core.Width {
			idx := r*core.Width + c
			if g.CPU.Display[idx] == true {
				screen.Set(c, r, color.White)
			} else {
				screen.Set(c, r, color.Black)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return core.Width, core.Height
}
