package vm

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

// RegularIncrement is the step value for all ops except halt
const RegularIncrement = 3

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^
//
func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {

		pc := registers[0]
		op := memory[pc] // fetch the opcode

		// decode and execute
		switch op {
		case Halt:
			return
		case Load:
			destReg := memory[pc+1]
			srcAddr := memory[pc+2]
			registers[destReg] = memory[srcAddr]
			registers[0] += RegularIncrement
		case Store:
			srcReg := memory[pc+1]
			srcVal := registers[srcReg]
			destAddr := memory[pc+2]
			memory[destAddr] = srcVal
			registers[0] += RegularIncrement
		case Add:
			srcReg1 := memory[pc+1]
			srcReg2 := memory[pc+2]
			registers[srcReg1] = registers[srcReg1] + registers[srcReg2]
			registers[0] += RegularIncrement
		case Sub:
			srcReg1 := memory[pc+1]
			srcReg2 := memory[pc+2]
			registers[srcReg1] = registers[srcReg1] - registers[srcReg2]
			registers[0] += RegularIncrement
		}
	}
}
