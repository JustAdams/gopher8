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
