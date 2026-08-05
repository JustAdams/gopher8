package main

import (
	"fmt"
	"gopher8/internal/core"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	cpu core.CPU
}

// ebitengine method which attempts to run at 60Hz
func (g *Game) Update() error {

	// user input
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.cpu.SetCurrentKey(0x1)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.cpu.SetCurrentKey(0x2)
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.cpu.SetCurrentKey(0x3)
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.cpu.SetCurrentKey(0xC)
	} else if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.cpu.SetCurrentKey(0x4)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cpu.SetCurrentKey(0x5)
	} else if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.cpu.SetCurrentKey(0x6)
	} else if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.cpu.SetCurrentKey(0xD)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cpu.SetCurrentKey(0x7)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cpu.SetCurrentKey(0x8)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cpu.SetCurrentKey(0x9)
	} else if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.cpu.SetCurrentKey(0xE)
	} else if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.cpu.SetCurrentKey(0xA)
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.cpu.SetCurrentKey(0x0)
	} else if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.cpu.SetCurrentKey(0xB)
	} else if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.cpu.SetCurrentKey(0xF)
	} else {
		g.cpu.SetCurrentKey(core.NoInput)
	}

	// delay timer is reduced at a rate of 60Hz until it reaches zero
	g.cpu.ReduceDelay()
	g.cpu.Cycle()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for r := range core.Height {
		for c := range core.Width {
			idx := r*core.Width + c
			if g.cpu.Display[idx] == true {
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

func main() {
	fmt.Println("Go Chip8")

	rom, err := core.CreateROM("roms/test_opcode.ch8")
	if err != nil {
		log.Fatal(err)
	}
	cpu := core.NewCPU()
	cpu.LoadBytes(core.RomStart, rom.Data[:])

	game := &Game{
		cpu: *cpu,
	}

	// scales the game to fit this resolution
	ebiten.SetWindowSize(640, 320)
	ebiten.SetWindowTitle("Gopher8")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
