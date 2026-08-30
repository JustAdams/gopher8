package core

import (
	"fmt"
	"io"
	"log"
)

const MaxROMSize uint16 = 3584

type ROM struct {
	Size uint16 // size of the ROM not to exceed MaxROMSize
	Data [MaxROMSize]byte
}

// Creates a ROM from a .ch8 file.
func CreateROM(romReader io.Reader) (*ROM, error) {
	rom := &ROM{}

	romLen, err := io.ReadFull(romReader, rom.Data[:])
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, fmt.Errorf("failed to read rom")
	}
	if romLen > int(MaxROMSize) {
		log.Fatal("Invalid ROM size provided")
	}
	rom.Size = uint16(romLen)

	return rom, nil
}
