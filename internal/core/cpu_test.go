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
