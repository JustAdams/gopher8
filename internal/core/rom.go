package core

import (
	"log"
	"os"
)

const MaxROMSize uint16 = 3584

type ROM struct {
	Size uint16 // size of the ROM not to exceed MaxROMSize
	Data [MaxROMSize]uint8
}

// Creates a ROM from a .ch8 file.
func CreateROM(romPath string) (*ROM, error) {
	file, err := os.Open(romPath)
	if err != nil {
		log.Fatalf("%s", "Cannot find "+romPath)
	}
	defer file.Close()

	rom := &ROM{}

	size, err := file.Read(rom.Data[:])
	if err != nil {
		log.Fatal(err)
	}
	if size > int(MaxROMSize) {
		log.Fatal("Invalid ROM size provided")
	}
	rom.Size = uint16(size)

	return rom, nil
}
