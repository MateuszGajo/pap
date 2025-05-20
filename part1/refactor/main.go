package main

import (
	"fmt"
	"log"
	"os"
)

type instruction string
type addressMode string

const (
	MOVE instruction = "MOVE"
)

const (
	RegisterMemoryToFromRegister addressMode = "RegisterMemoryToFromRegister"
	ImmediateToRegisterMemory    addressMode = "ImmediateToRegisterMemory"
	ImmediateToRegister          addressMode = "ImmediateToRegister"
	MemoryToAccumulator          addressMode = "MemoryToAccumulator"
	AccumulatorToMemory          addressMode = "AccumulatorToMemory"
)

func parseInstruction(data []byte) (instruction, addressMode) {
	var instruction instruction
	var addressMode addressMode
	if data[0]>>2 == 0b00100010 {
		// register/memory to/from register
		// d = data[0] >> 1 & 1
		// w = data[0] & 1
		instruction = MOVE
		addressMode = RegisterMemoryToFromRegister
	} else if data[0]>>1 == 0b01100011 {
		// immediate to register/memory
		// w = data[0] & 1
		instruction = MOVE
		addressMode = ImmediateToRegisterMemory
	} else if data[0]>>4 == 0b00001011 {
		// immediate to register
		// w = data[0] >> 3 & 1
		// reg = data[0] & 111
		instruction = MOVE
		addressMode = ImmediateToRegister
	} else if data[0]>>1 == 0b01010000 {
		// memory to accumulator
		// w = data[0] & 1
		instruction = MOVE
		addressMode = MemoryToAccumulator
	} else if data[0]>>1 == 0b01010001 {
		// accumulator to memory
		// w = data[0] & 1
		instruction = MOVE
		addressMode = AccumulatorToMemory
	}
	return instruction, addressMode
}

var registers_8 = []string{"AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH"}
var registers_16 = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}

func getRegisterOrMemory(mod, rm byte, w *byte) string {

	if (mod == 0b11) != (w != nil) {

		panic("w should be only defined for mod 11")
	}

	if mod > 0b111 || rm > 0b111 {
		fmt.Println(rm > 0b111)
		fmt.Println("mod: %v, rm: %v, w: %v", mod, rm, w)
		panic("mod and rm should at most 111")
	}

	if mod == 0b00 {
		// memory mode no displacement
	} else if mod == 0b01 {
		// memory mode, 8 bit displacement follows
	} else if mod == 0b10 {
		// memory mode, 16 bit displacement follows
	} else if mod == 0b11 {
		if *w > 0b1 {
			panic("w should be only 0 or 1")
		}
		// register mode (no displacement)
		if *w == 0 {
			return registers_8[rm]
		} else {
			return registers_16[rm]
		}
	}
	fmt.Errorf("mod: %v, rm: %v, w: %v", mod, rm, w)
	panic("Should never happend")
}

func getRegister(reg, w byte) string {
	if reg > 0b111 {
		panic("reg should be at most 111")
	}

	if w > 0b1 {
		panic("w should be only 0 or 1")
	}

	if w == 0 {
		return registers_8[reg]
	} else {
		return registers_16[reg]
	}

}

func handleImmediateToRegister(instruction instruction, addressMode addressMode, data []byte) {
	// w := data[0] >> 3
}

func handleRegisterMemoryToFromRegister(instruction instruction, addressMode addressMode, data []byte) (string, string) {
	d := data[0] >> 1 & 0b1
	w := data[0] & 1
	mod := data[1] >> 6 & 0b11
	reg := data[1] >> 3 & 0b111
	rm := data[1] & 0b111

	dst := getRegisterOrMemory(mod, rm, &w)
	src := getRegister(reg, w)

	if d == 1 {
		dst, src = src, dst
	}

	return dst, src
}

func exectuteInstruction(instruction instruction, addressMode addressMode, data []byte) (string, string) {
	switch addressMode {
	case ImmediateToRegister:
		handleImmediateToRegister(instruction, addressMode, data)
	case RegisterMemoryToFromRegister:
		return handleRegisterMemoryToFromRegister(instruction, addressMode, data)
	}
	return "", ""
}

func main() {
	data, err := os.ReadFile("./data/listing_0037_single_register_mov")

	if err != nil {
		log.Fatal(err)
	}
	instruction, addressMode := parseInstruction(data)

	dst, src := exectuteInstruction(instruction, addressMode, data)

	fmt.Printf("%v %v, %v \n", instruction, dst, src)

}
