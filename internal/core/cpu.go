package core

import (
	"fmt"
)

const FontStart = 0x50
const RomStart = 0x200

const Height = 32
const Width = 64

const NoInput = 0xFF

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
	delay   uint8
	sound   uint8
	currKey uint8
	Display [Height * Width]bool
}

func NewCPU() *CPU {
	cpu := CPU{
		pc: 0x200,
		sp: 0x00,
	}

	cpu.LoadBytes(FontStart, font[:])

	return &cpu
}

// reduces the delay timer by 1 to a minimum of 0
func (cpu *CPU) ReduceDelay() {
	if cpu.delay == 0 {
		cpu.delay = 0
	} else {
		cpu.delay -= 1
	}
}

func (cpu *CPU) SetCurrentKey(key uint8) {
	cpu.currKey = key
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
			cpu.op0x00E0()
		}
	case 0x1:
		cpu.op0x1NNN(opcode.NNN)
	case 0x3:
		cpu.op0x3XNN(opcode.X, opcode.NN)
	case 0x5:
		cpu.op0x5XY0(opcode.X, opcode.Y)
	case 0x6:
		cpu.op0x6XNN(opcode.X, opcode.NN)
	case 0x7:
		cpu.op0x7XNN(opcode.X, opcode.NN)
	case 0x8:
		switch opcode.N {
		case 0x0:
			cpu.op0x8XY0(opcode.X, opcode.Y)
		case 0x1:
			cpu.op0x8XY1(opcode.X, opcode.Y)
		case 0x2:
			cpu.op0x8XY2(opcode.X, opcode.Y)
		case 0x3:
			cpu.op0x8XY3(opcode.X, opcode.Y)
		case 0x4:
			cpu.op0x8XY4(opcode.X, opcode.Y)
		}
	case 0x9:
		cpu.op0x9XY0(opcode.X, opcode.Y)
	case 0xA:
		cpu.op0xANNN(opcode.NNN)
	case 0xD:
		cpu.op0xDXYN(opcode.X, opcode.Y, opcode.N)
	case 0xF:
		switch opcode.NN {
		case 0x07:
			cpu.op0xF07(opcode.X)
		case 0x15:
			cpu.op0xF15(opcode.X)
		case 0x18:
			cpu.op0xF18(opcode.X)
		}
	default:
		fmt.Printf("Unable to find opcode %d", opcode)
	}
}

// 0x00E0 - clears the screen.
func (cpu *CPU) op0x00E0() {
	for r := range Height {
		for c := range Width {
			idx := r*Height + c
			cpu.Display[idx] = false
		}
	}
}

// 0x1NNN - jumps pc to pos.
func (cpu *CPU) op0x1NNN(pos uint16) {
	cpu.pc = pos
}

// 0x3XNN - skip if vX equals NN
func (cpu *CPU) op0x3XNN(x, nn uint8) {
	if cpu.v[x] == nn {
		cpu.pc += 2
	}
}

// 0x4XNN - skip if vX doesnt equals NN
func (cpu *CPU) op0x4XNN(x, nn uint8) {
	if cpu.v[x] != nn {
		cpu.pc += 2
	}
}

// 0x5XY0 - skip if values in vX and vY are equal
func (cpu *CPU) op0x5XY0(x, y uint8) {
	if cpu.v[x] == cpu.v[y] {
		cpu.pc += 2
	}
}

// 0x6XNN - set register vx
func (cpu *CPU) op0x6XNN(x, nn uint8) {
	cpu.v[x] = nn
}

// 0x7XNN - add nn to vx
func (cpu *CPU) op0x7XNN(x, nn uint8) {
	cpu.v[x] += nn
}

// 0x8XY0 - set vx to the value of vy
func (cpu *CPU) op0x8XY0(x, y uint8) {
	cpu.v[x] = cpu.v[y]
}

// 0x8XY1 - set vx to OR of vx and vy
func (cpu *CPU) op0x8XY1(x, y uint8) {
	cpu.v[x] |= cpu.v[y]
}

// 0x8XY2 - set vx to AND of vx and vy
func (cpu *CPU) op0x8XY2(x, y uint8) {
	cpu.v[x] &= cpu.v[y]
}

// 0x8XY3 - set vx to XOR of vx and vy
func (cpu *CPU) op0x8XY3(x, y uint8) {
	cpu.v[x] ^= cpu.v[y]
}

// 0x8XY4 - adds vy to vx
func (cpu *CPU) op0x8XY4(x, y uint8) {
	cpu.v[x] += cpu.v[y]
	// todo: finish overflow logic
}

// 0x9XY0 - skip if values in vX and vY are not equal
func (cpu *CPU) op0x9XY0(x, y uint8) {
	if cpu.v[x] != cpu.v[y] {
		cpu.pc += 2
	}
}

// 0xANNN - set index register to nnn
func (cpu *CPU) op0xANNN(nnn uint16) {
	cpu.idxReg = nnn
}

// 0xDXYN - draws to the display.
func (cpu *CPU) op0xDXYN(x, y, n uint8) {
	// starting coordinates from top-left of sprite
	xCoord := cpu.v[x]
	yCoord := cpu.v[y]

	for r := range n {
		// get the nth byte of sprite data from memory address at idxReg
		spriteData := cpu.ram[cpu.idxReg+uint16(r)]
		yPos := (yCoord + r) & (Height - 1)
		if yPos >= Height {
			break
		}

		// for each of the 8 bits in the sprite row
		for c := range 8 {
			spritePixel := (spriteData>>(7-c))&0x1 == 1
			if !spritePixel {
				continue
			}

			xPos := xCoord + uint8(c)
			if xPos >= Width {
				break
			}

			idx := (uint16(yPos) * Width) + uint16(xPos)

			currPixel := cpu.Display[idx]
			if currPixel {
				cpu.v[0xF] = 0x1
			}
			cpu.Display[idx] = !cpu.Display[idx]
		}
	}
}

// 0xFX07 - sets vx to the value of the delay timer
func (cpu *CPU) op0xF07(x uint8) {
	cpu.v[x] = cpu.delay
}

// 0xFX0A - waits for key input and stores it in vx
func (cpu *CPU) op0xFX0A(x uint8) {
	if cpu.currKey == NoInput {
		cpu.pc -= 2
	} else {
		cpu.v[x] = cpu.currKey
	}
}

// 0xFX15 - sets the delay timer to the value at vx
func (cpu *CPU) op0xF15(x uint8) {
	cpu.delay = cpu.v[x]
}

// 0xFX18 - sets the sound timer to the value at vx
func (cpu *CPU) op0xF18(x uint8) {
	cpu.sound = cpu.v[x]
}
