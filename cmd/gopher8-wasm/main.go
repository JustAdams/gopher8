package main

import (
	"bytes"
	_ "embed"
	"gopher8/internal/core"
	"gopher8/internal/gui"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 640
	windowHeight = 320
)

func main() {
	// read rom from input otherwise default to test
	romPath := "outlaw.ch8"
	if len(os.Args) > 1 {
		romPath = os.Args[len(os.Args)-1] + ".ch8"
	}

	game := NewGame(romPath)

	// scales the game to fit this resolution
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Gopher8")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

//go:embed outlaw.ch8
var defaultROM []byte

func NewGame(romPath string) *gui.Game {
	romReader := bytes.NewReader(defaultROM)
	rom, err := core.CreateROM(romReader)
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
