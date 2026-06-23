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

func (g *Game) Update() error {
	g.cpu.Cycle()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for r := range core.Height {
		for c := range core.Width {
			if g.cpu.Display[r][c] == true {
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

	ebiten.SetWindowSize(core.Width*10, core.Height*10)
	ebiten.SetWindowTitle("Gopher8")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
