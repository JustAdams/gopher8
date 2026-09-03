package core

import (
	"fmt"
	"math/rand/v2"
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
	ram       [4096]uint8
	v         [16]uint8
	idxReg    uint16
	pc        uint16
	sp        uint16
	callStack stack
	delay     uint8
	sound     uint8
	currKey   uint8
	Display   [Height * Width]bool
}

func NewCPU() *CPU {
	cpu := CPU{
		pc: 0x200,
		sp: 0x00,
	}

	cpu.LoadBytes(FontStart, font[:])

	return &cpu
}

// ReduceDelay reduces the delay timer by 1 to a minimum of 0
func (cpu *CPU) ReduceDelay() {
	if cpu.delay > 0 {
		cpu.delay--
	}

	if cpu.sound > 0 {
		cpu.sound--
		if cpu.sound == 0 {
			// play sound
		}
	} else {
		// stop sound
	}
}

func (cpu *CPU) SetCurrentKey(key uint8) {
	cpu.currKey = key
}

func (cpu *CPU) Cycle() {
	// fetch instruction
	instruction := uint16(cpu.ram[cpu.pc]) << 8
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
		case 0x0EE:
			cpu.op0x00EE()
		}
	case 0x1:
		cpu.op0x1NNN(opcode.NNN)
	case 0x2:
		cpu.op0x2NNN(opcode.NNN)
	case 0x3:
		cpu.op0x3XNN(opcode.X, opcode.NN)
	case 0x4:
		cpu.op0x4XNN(opcode.X, opcode.NN)
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
		case 0x5:
			cpu.op0x8XY5(opcode.X, opcode.Y)
		case 0x6:
			cpu.op0x8XY6(opcode.X, opcode.Y)
		case 0x7:
			cpu.op0x8XY7(opcode.X, opcode.Y)
		case 0xE:
			cpu.op0x8XYE(opcode.X, opcode.Y)
		}
	case 0x9:
		cpu.op0x9XY0(opcode.X, opcode.Y)
	case 0xA:
		cpu.op0xANNN(opcode.NNN)
	case 0xB:
		cpu.op0xBXNN(opcode.NNN)
	case 0xC:
		cpu.op0xCXNN(opcode.X, opcode.NN)
	case 0xD:
		cpu.op0xDXYN(opcode.X, opcode.Y, opcode.N)
	case 0xE:
		switch opcode.NN {
		case 0x9E:
			cpu.op0xEX9E(opcode.X)
		case 0xA1:
			cpu.op0xEXA1(opcode.X)
		}
	case 0xF:
		switch opcode.NN {
		case 0x07:
			cpu.op0xF07(opcode.X)
		case 0x0A:
			cpu.op0xFX0A(opcode.X)
		case 0x15:
			cpu.op0xF15(opcode.X)
		case 0x18:
			cpu.op0xF18(opcode.X)
		case 0x1E:
			cpu.op0xF1E(opcode.X)
		case 0x29:
			cpu.op0xF29(opcode.X)
		case 0x33:
			cpu.op0xFX33(opcode.X)
		case 0x55:
			cpu.op0xFX55(opcode.X)
		case 0x65:
			cpu.op0xFX65(opcode.X)
		}
	default:
		fmt.Printf("Unable to find opcode %d", opcode)
	}
}

// 0x00E0 - clears the screen.
func (cpu *CPU) op0x00E0() {
	for r := range Height {
		for c := range Width {
			idx := r*Width + c
			cpu.Display[idx] = false
		}
	}
}

// 0x00EE - returns from subroutine
func (cpu *CPU) op0x00EE() {
	cpu.pc = cpu.callStack.pop()
}

// 0x1NNN - jumps pc to pos
func (cpu *CPU) op0x1NNN(pos uint16) {
	cpu.pc = pos
}

// 0x2NNN - calls subroutine at memory location NNN
func (cpu *CPU) op0x2NNN(nnn uint16) {
	// push PC to the stack first so it can be resumed after subroutine completes
	cpu.callStack.push(cpu.pc)
	cpu.pc = nnn
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
	sum := uint16(cpu.v[x] + cpu.v[y])
	if sum > 0xFF {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] = uint8(sum)
}

// 0x8XY5 - sets vx to the result of vx - vy
func (cpu *CPU) op0x8XY5(x, y uint8) {
	if cpu.v[x] >= cpu.v[y] {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] -= cpu.v[y]
}

// 0x8XY6 - sets vx equal to vy and shifts vx one bit to the right
func (cpu *CPU) op0x8XY6(x, y uint8) {
	cpu.v[0xF] = cpu.v[y] & 1
	cpu.v[x] = cpu.v[y] >> 1
}

// 0x8XY7 - sets vx to the result of vy - vx
func (cpu *CPU) op0x8XY7(x, y uint8) {
	if cpu.v[y] >= cpu.v[x] {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] = cpu.v[y] - cpu.v[x]
}

// 0x8XYE - sets vx equal to vy and shifts vx one bit to the left
func (cpu *CPU) op0x8XYE(x, y uint8) {
	msb := (cpu.v[y] & 0x80) >> 7
	cpu.v[x] = cpu.v[y] << 1
	cpu.v[0xF] = msb
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

// 0xBXNN - jump to address xnn plus the value in vx
func (cpu *CPU) op0xBXNN(nnn uint16) {
	nextPos := nnn + uint16(cpu.v[0])
	cpu.pc = nextPos
}

// 0xCXNN - generate random number and sets to vx
func (cpu *CPU) op0xCXNN(x, nn uint8) {
	num := uint8(rand.UintN(256)) & nn
	cpu.v[x] = num
}

// 0xDXYN - draws to the display.
func (cpu *CPU) op0xDXYN(x, y, n uint8) {
	// get x and y coordinates from var registers
	xCoord := cpu.v[x] % Width
	yCoord := cpu.v[y] % Height
	cpu.v[0xF] = 0

	for row := uint8(0); row < n; row++ {
		yPos := yCoord + row
		// stop if the end of display is reached
		if yPos >= Height {
			break
		}

		spriteByte := cpu.ram[cpu.idxReg+uint16(row)]

		// for each of the 8 bits in the row
		for col := uint8(0); col < 8; col++ {
			xPos := xCoord + col
			// stop if end of display is reached
			if xPos >= Width {
				break
			}

			if (spriteByte & (0x80 >> col)) != 0 {
				pos := uint16(yPos)*uint16(Width) + uint16(xPos)

				if cpu.Display[pos] {
					cpu.Display[pos] = false
					cpu.v[0xF] = 1
				} else {
					cpu.Display[pos] = true
				}
			}

		}
	}
}

// 0xEX9E - skip if vx equals current key
func (cpu *CPU) op0xEX9E(x uint8) {
	if cpu.v[x] == cpu.currKey {
		cpu.pc += 2
	}
}

// 0xEXA1 - skip if vx doesnt equals current key
func (cpu *CPU) op0xEXA1(x uint8) {
	if cpu.v[x] != cpu.currKey {
		cpu.pc += 2
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

// 0xFX1E
func (cpu *CPU) op0xF1E(x uint8) {
	cpu.idxReg += uint16(cpu.v[x])
}

// 0xFX29 - sets idx register to address of the hexdecimal character in vx
func (cpu *CPU) op0xF29(x uint8) {
	cpu.idxReg = uint16(FontStart + (cpu.v[x] * 5))
}

// 0xFX33 - binary-coded decimal conversion
func (cpu *CPU) op0xFX33(x uint8) {
	// split vx into three digits
	val := cpu.v[x]
	for i := 2; i >= 0; i-- {
		cpu.ram[cpu.idxReg+uint16(i)] = val % 10
		val /= 10
	}
}

// 0xFX55 - store
func (cpu *CPU) op0xFX55(x uint8) {
	// value of each variable register from v0 to vx is stored in successive memory addresses
	vIdx := 0
	for i := uint8(0); i <= x; i++ {
		cpu.ram[cpu.idxReg+uint16(i)] = cpu.v[vIdx]
		vIdx++
	}
}

// 0xFX65 - load
func (cpu *CPU) op0xFX65(x uint8) {
	// value of each memory address from idx to idx + x is stored in in variable memory starting a 0
	vIdx := 0
	for i := uint8(0); i <= x; i++ {
		cpu.v[vIdx] = cpu.ram[cpu.idxReg+uint16(i)]
		vIdx++
	}
}
