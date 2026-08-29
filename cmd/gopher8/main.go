package main

import (
	"gopher8/internal/core"
	"gopher8/internal/gui"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 640
	windowHeight = 320

	chip8Width  = 64
	chip8Height = 32
)

func main() {
	// read rom from input otherwise default to test
	romPath := "roms/test_opcode.ch8"
	if len(os.Args) > 1 {
		romPath = "roms/" + os.Args[len(os.Args)-1] + ".ch8"
	}

	game := NewGame(romPath)

	// scales the game to fit this resolution
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Gopher8")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func NewGame(romPath string) *gui.Game {
	rom, err := core.CreateROM(romPath)
	if err != nil {
		log.Fatalf("%s", "Error loading "+romPath)
	}
	cpu := core.NewCPU()
	cpu.LoadBytes(core.RomStart, rom.Data[:])

	game := &gui.Game{
		CPU: cpu,
	}
	return game
}
