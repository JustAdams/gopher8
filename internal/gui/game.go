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

func handleInput(g *Game) {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.CPU.SetCurrentKey(0x1)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.CPU.SetCurrentKey(0x2)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.CPU.SetCurrentKey(0x3)
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.CPU.SetCurrentKey(0xC)
	} else if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.CPU.SetCurrentKey(0x4)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.CPU.SetCurrentKey(0x5)
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.CPU.SetCurrentKey(0x6)
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.CPU.SetCurrentKey(0xD)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.CPU.SetCurrentKey(0x7)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.CPU.SetCurrentKey(0x8)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.CPU.SetCurrentKey(0x9)
	} else if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.CPU.SetCurrentKey(0xE)
	} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.CPU.SetCurrentKey(0xA)
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.CPU.SetCurrentKey(0x0)
	} else if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.CPU.SetCurrentKey(0xB)
	} else if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.CPU.SetCurrentKey(0xF)
	} else {
		g.CPU.SetCurrentKey(core.NoInput)
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
