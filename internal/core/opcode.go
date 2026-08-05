package core

type OpCode struct {
	W    uint8  // first nibble
	X    uint8  // second nibble
	Y    uint8  // third nibble
	N    uint8  // fourth nibble
	NN   uint8  // third and fourth nibbles
	NNN  uint16 // second, third, and fourth nibbles
	NNNN uint16 // full opcode
}

func NewOpCode(instruction uint16) *OpCode {
	return &OpCode{
		W:    uint8((instruction & 0xF000) >> 12),
		X:    uint8((instruction & 0x0F00) >> 8),
		Y:    uint8((instruction & 0x00F0) >> 4),
		N:    uint8(instruction & 0x000F),
		NN:   uint8(instruction & 0x00FF),
		NNN:  (instruction & 0x0FFF),
		NNNN: instruction,
	}
}
