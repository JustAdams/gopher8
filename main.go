package main

import (
	"fmt"
	"gopher8/internal/core"
	"log"
)

func main() {
	fmt.Println("Go Chip8")

	rom, err := core.CreateROM("roms/IBM Logo.ch8")
	if err != nil {
		log.Fatal(err)
	}
	cpu := core.NewCPU()
	cpu.LoadBytes(core.RomStart, rom.Data[:])

	for {
		cpu.Cycle()

		// draw to terminal
		for r := range core.Height {
			for c := range core.Width {
				if cpu.Display[r][c] == true {
					fmt.Print(" *")
				} else {
					fmt.Print("  ")
				}
			}
			fmt.Println()
		}

	}
}
