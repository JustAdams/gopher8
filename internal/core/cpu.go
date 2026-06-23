package core

import (
	"fmt"
)

const FontStart = 0x50
const RomStart = 0x200

const Height = 32
const Width = 64

var font = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type CPU struct {
	ram     [4096]uint8
	v       [16]uint8
	idxReg  uint16
	pc      uint16
	sp      uint16
	Display [Height][Width]bool
}

func NewCPU() *CPU {
	cpu := CPU{
		pc: 0x200,
		sp: 0x00,
	}

	cpu.LoadBytes(FontStart, font[:])

	return &cpu
}

func (cpu *CPU) Cycle() {
	// fetch instruction
	var instruction uint16 = uint16(cpu.ram[cpu.pc]) << 8
	instruction |= uint16(cpu.ram[cpu.pc+1])
	cpu.pc += 2

	opCode := NewOpCode(instruction)

	// execute instruction
	cpu.execute(opCode)
}

func (cpu *CPU) LoadBytes(loadPos uint16, payload []uint8) {
	copy(cpu.ram[loadPos:loadPos+uint16(len(payload))], payload[:])
}

// Performs action based on the provided opcode.
func (cpu *CPU) execute(opcode *OpCode) {

	switch opcode.W {
	case 0x0:
		switch opcode.NNN {
		case 0x0E0:
			cpu.opClear()
			cpu.pc += 2
		}
	case 0x1:
		cpu.opJump(opcode.NNN)
	case 0x3:
		cpu.opSkipVxNN(opcode.X, opcode.NN)
	case 0x5:
		cpu.opSkipVxEqualVy(opcode.X, opcode.Y)
	case 0x6:
		cpu.opSet(opcode.X, opcode.NN)
	case 0x7:
		cpu.opAdd(opcode.X, opcode.NN)
	case 0x8:
		switch opcode.N {
		case 0x0:
			cpu.opSetVxVy(opcode.X, opcode.Y)
		case 0x1:
			cpu.opOr(opcode.X, opcode.Y)
		case 0x2:
			cpu.opAnd(opcode.X, opcode.Y)
		case 0x3:
			cpu.opXor(opcode.X, opcode.Y)
		case 0x4:
			cpu.opAddVxVy(opcode.X, opcode.Y)
		}
	case 0x9:
		cpu.opSkipVxNotEqualVy(opcode.X, opcode.Y)
	case 0xA:
		cpu.opSetIndex(opcode.NNN)
	case 0xD:
		cpu.opDraw(opcode.X, opcode.Y, opcode.N)
	default:
		fmt.Printf("Unable to find opcode %d", opcode)
	}
}

// 0x00E0 - clears the screen.
func (cpu *CPU) opClear() {
	for r := range Height {
		for c := range Width {
			cpu.Display[r][c] = false
		}
	}
}

// 0x1NNN - jumps pc to pos.
func (cpu *CPU) opJump(pos uint16) {
	cpu.pc = pos
}

// 0x3XNN - skip if vX equals NN
func (cpu *CPU) opSkipVxNN(x, nn uint8) {
	if cpu.v[x] == nn {
		cpu.pc += 2
	}
}

// 0x4XNN - skip if vX doesnt equals NN
func (cpu *CPU) opNotSkipVxNN(x, nn uint8) {
	if cpu.v[x] != nn {
		cpu.pc += 2
	}
}

// 0x5XY0 - skip if values in vX and vY are equal
func (cpu *CPU) opSkipVxEqualVy(x, y uint8) {
	if cpu.v[x] == cpu.v[y] {
		cpu.pc += 2
	}
}

// 0x6XNN - set register vx
func (cpu *CPU) opSet(x, nn uint8) {
	cpu.v[x] = nn
}

// 0x7XNN - add nn to vx
func (cpu *CPU) opAdd(x, nn uint8) {
	cpu.v[x] += nn
}

// 0x8XY0 - set vx to the value of vy
func (cpu *CPU) opSetVxVy(x, y uint8) {
	cpu.v[x] = cpu.v[y]
}

// 0x8XY1 - set vx to OR of vx and vy
func (cpu *CPU) opOr(x, y uint8) {
	cpu.v[x] |= cpu.v[y]
}

// 0x8XY2 - set vx to AND of vx and vy
func (cpu *CPU) opAnd(x, y uint8) {
	cpu.v[x] &= cpu.v[y]
}

// 0x8XY3 - set vx to XOR of vx and vy
func (cpu *CPU) opXor(x, y uint8) {
	cpu.v[x] ^= cpu.v[y]
}

// 0x8XY4 - adds vy to vx
func (cpu *CPU) opAddVxVy(x, y uint8) {
	cpu.v[x] += cpu.v[y]
	// todo: finish overflow logic
}

// 0x9XY0 - skip if values in vX and vY are not equal
func (cpu *CPU) opSkipVxNotEqualVy(x, y uint8) {
	if cpu.v[x] != cpu.v[y] {
		cpu.pc += 2
	}
}

// 0xANNN - set index register to nnn
func (cpu *CPU) opSetIndex(nnn uint16) {
	cpu.idxReg = nnn
}

// 0xDXYN - draws to the display.
func (cpu *CPU) opDraw(x, y, n uint8) {
	xCoord := cpu.v[x]
	yCoord := cpu.v[y]

	for r := range n {
		spriteData := cpu.ram[cpu.idxReg+uint16(r)]

		for c := range 8 {
			spritePixel := (spriteData>>(7-c))&0x1 == 1
			if !spritePixel {
				continue
			}

			xPos := xCoord + uint8(c)
			if xPos >= Width {
				break
			}

			yPos := (yCoord + r) & (Height - 1)

			currPixel := cpu.Display[yPos][xPos]
			if currPixel {
				cpu.v[0xF] = 0x1
			}
			cpu.Display[yPos][xPos] = !cpu.Display[yPos][xPos]
		}
	}
}
