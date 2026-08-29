package core

import "testing"

// Test_loadFont calls cpu.loadFont and checks that the font is loaded into the correct memory location.
func Test_loadFont(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadBytes(FontStart, font[:])

	var expectedFontValue uint8 = font[0]
	var actualFontValue uint8 = cpu.ram[FontStart]
	if actualFontValue != expectedFontValue {
		t.Errorf(`First font memory address contains %d instead of %d`, actualFontValue, expectedFontValue)
	}
}

func Test_0x8XY6(t *testing.T) {
	var initialVY uint8 = 0xAB
	var expectedVX uint8 = 0x55

	payload := []uint8{0x81, 0x26}
	cpu := NewCPU()
	cpu.LoadBytes(cpu.pc, payload)

	cpu.v[0x2] = initialVY
	cpu.Cycle()

	// vf should contain a flag
	if cpu.v[0xF] != 0x1 {
		t.Error(`VF should contain a flag after the bit shift`)
	}

	// vx should become equal to vy plus one bit-shift right
	var actualVX uint8 = cpu.v[0x1]
	if actualVX != expectedVX {
		t.Errorf(`VX should be %d after bit-shift right but was %d`, actualVX, expectedVX)
	}
}

func Test_0x8XYE(t *testing.T) {
	var initialVY uint8 = 0xF
	var expectedVX uint8 = 0x1e

	payload := []uint8{0x81, 0x2E}
	cpu := NewCPU()
	cpu.LoadBytes(cpu.pc, payload)

	cpu.v[0x2] = initialVY
	cpu.Cycle()

	// vf should contain a flag
	if cpu.v[0xF] != 0x1 {
		t.Error(`VF should contain a flag after the bit shift`)
	}

	// vx should become equal to vy plus one bit-shift left
	var actualVX uint8 = cpu.v[0x1]
	if actualVX != expectedVX {
		t.Errorf(`VX should be %d after bit-shift right but was %d`, actualVX, expectedVX)
	}
}

func TestCPU_op0xFX33(t *testing.T) {
	payload := []uint8{0xF1, 0x33}
	cpu := NewCPU()
	cpu.LoadBytes(cpu.pc, payload)
	cpu.v[1] = 123

	cpu.Cycle()

	if cpu.ram[cpu.idxReg] != 0x1 {
		t.Errorf(`Address I should contain %d but was %d`, 1, cpu.ram[cpu.idxReg])
	}
	if cpu.ram[cpu.idxReg+1] != 0x2 {
		t.Errorf(`Address I should contain %d but was %d`, 2, cpu.ram[cpu.idxReg])
	}
	if cpu.ram[cpu.idxReg+2] != 0x3 {
		t.Errorf(`Address I should contain %d but was %d`, 3, cpu.ram[cpu.idxReg])
	}
}

func TestCPU_op0xFX55(t *testing.T) {
	payload := []uint8{0xF2, 0x55}
	cpu := NewCPU()
	cpu.LoadBytes(cpu.pc, payload)
	cpu.v[0] = 1
	cpu.v[1] = 2
	cpu.v[2] = 3

	cpu.Cycle()

	if cpu.ram[cpu.idxReg] != 0x1 {
		t.Errorf(`Address I should contain %d but was %d`, 1, cpu.ram[cpu.idxReg])
	}
	if cpu.ram[cpu.idxReg+1] != 0x2 {
		t.Errorf(`Address I should contain %d but was %d`, 2, cpu.ram[cpu.idxReg])
	}
	if cpu.ram[cpu.idxReg+2] != 0x3 {
		t.Errorf(`Address I should contain %d but was %d`, 3, cpu.ram[cpu.idxReg])
	}
}

func TestCPU_op0xFX65(t *testing.T) {
	payload := []uint8{0xF2, 0x65}
	cpu := NewCPU()
	cpu.LoadBytes(cpu.pc, payload)
	cpu.ram[cpu.idxReg] = 1
	cpu.ram[cpu.idxReg+1] = 2
	cpu.ram[cpu.idxReg+2] = 3

	cpu.Cycle()

	if cpu.v[0] != 0x1 {
		t.Errorf(`Address I should contain %d but was %d`, 1, cpu.ram[cpu.idxReg])
	}
	if cpu.v[1] != 0x2 {
		t.Errorf(`Address I should contain %d but was %d`, 2, cpu.ram[cpu.idxReg])
	}
	if cpu.v[2] != 0x3 {
		t.Errorf(`Address I should contain %d but was %d`, 3, cpu.ram[cpu.idxReg])
	}
}
