package core

import "testing"

func Test_NewOpCode(t *testing.T) {
	var instruction uint16 = 0xABCD
	opCode := NewOpCode(instruction)

	if opCode.W != 0xA {
		t.Errorf("First nibble is %d, expected %d", opCode.W, 0xA)
	}
	if opCode.X != 0xB {
		t.Errorf("Second nibble is %d, expected %d", opCode.X, 0xB)
	}
	if opCode.Y != 0xC {
		t.Errorf("Third nibble is %d, expected %d", opCode.Y, 0xC)
	}
	if opCode.N != 0xD {
		t.Errorf("Fourth nibble is %d, expected %d", opCode.N, 0xB)
	}
	if opCode.NN != 0xCD {
		t.Errorf("NN is %d, expected %d", opCode.NN, 0xCD)
	}
	if opCode.NNN != 0xBCD {
		t.Errorf("NNN is %d, expected %d", opCode.NNN, 0xBCD)
	}
}
