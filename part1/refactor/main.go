package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
)

type instruction string
type addressMode string

const (
	MOVE instruction = "mov"
)

const (
	RegisterMemoryToFromRegister addressMode = "RegisterMemoryToFromRegister"
	ImmediateToRegisterMemory    addressMode = "ImmediateToRegisterMemory"
	ImmediateToRegister          addressMode = "ImmediateToRegister"
	MemoryToAccumulator          addressMode = "MemoryToAccumulator"
	AccumulatorToMemory          addressMode = "AccumulatorToMemory"
)

func (cpu *Cpu) parseInstruction() (instruction, addressMode) {
	var instruction instruction
	var addressMode addressMode

	cpu.readNextByte()
	if cpu.currentByte>>2 == 0b00100010 {
		// register/memory to/from register
		// d = cpu.currentByte >> 1 & 1
		// w = cpu.currentByte & 1
		instruction = MOVE
		addressMode = RegisterMemoryToFromRegister
	} else if cpu.currentByte>>1 == 0b01100011 {
		// immediate to register/memory
		// w = cpu.currentByte & 1
		instruction = MOVE
		addressMode = ImmediateToRegisterMemory
	} else if cpu.currentByte>>4 == 0b00001011 {
		// immediate to register
		// w = cpu.currentByte >> 3 & 1
		// reg = cpu.currentByte & 111
		instruction = MOVE
		addressMode = ImmediateToRegister
	} else if cpu.currentByte>>1 == 0b01010000 {
		// memory to accumulator
		// w = cpu.currentByte & 1
		instruction = MOVE
		addressMode = MemoryToAccumulator
	} else if cpu.currentByte>>1 == 0b01010001 {
		// accumulator to memory
		// w = instructionData & 1
		instruction = MOVE
		addressMode = AccumulatorToMemory
	}
	return instruction, addressMode
}

var registers_8 = []string{"AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH"}
var registers_16 = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}
var memoryRegisters = []string{
	"BX + SI",
	"BX + DI",
	"BP + SI",
	"BP + DI",
	"SI",
	"DI",
	"BP",
	"BX",
}

func (cpu *Cpu) getRegisterOrMemory() string {

	if cpu.mod == nil {
		panic("mod is required in get registerOrMemory")
	}

	if *cpu.mod == 0b00 {
		// memory mode no displacement
		if *cpu.rm == 0b110 {
			// special case
			cpu.readNextByte()
			firstByte := cpu.currentByte
			cpu.readNextByte()
			val := int16(binary.LittleEndian.Uint16(append([]byte{}, firstByte, cpu.currentByte)))
			return fmt.Sprintf("[%v]", val)
		}
		return fmt.Sprintf("[%v]", memoryRegisters[*cpu.rm])
	} else if *cpu.mod == 0b01 {
		// memory mode, 8 bit displacement follows
		cpu.readNextByte()
		val := int8(cpu.currentByte)
		if val < 0 {
			return fmt.Sprintf("[%v - %v]", memoryRegisters[*cpu.rm], math.Abs(float64(val)))
		}
		return fmt.Sprintf("[%v + %v]", memoryRegisters[*cpu.rm], val)
	} else if *cpu.mod == 0b10 {
		cpu.readNextByte()
		firstByte := cpu.currentByte
		cpu.readNextByte()
		val := int16(binary.LittleEndian.Uint16(append([]byte{}, firstByte, cpu.currentByte)))
		if val < 0 {
			return fmt.Sprintf("[%v - %v]", memoryRegisters[*cpu.rm], math.Abs(float64(val)))
		}
		return fmt.Sprintf("[%v + %v]", memoryRegisters[*cpu.rm], val)
	} else if *cpu.mod == 0b11 {
		if cpu.w == nil {
			panic("w is required for mod 11")
		}
		// register mode (no displacement)
		if *cpu.w == 0 {
			return registers_8[*cpu.rm]
		} else {
			return registers_16[*cpu.rm]
		}
	}
	fmt.Errorf("mod: %v, rm: %v, w: %v", cpu.mod, cpu.rm, cpu.w)
	panic("Should never happend")
}

func (cpu *Cpu) getRegister() string {
	if cpu.reg == nil || cpu.w == nil {
		fmt.Println("reg", cpu.reg, "w", cpu.w)
		panic("reg and w is required in get register")
	}

	if *cpu.reg > 0b111 {
		panic("reg should be at most 111")
	}

	if *cpu.w > 0b1 {
		panic("w should be only 0 or 1")
	}

	if *cpu.w == 0 {
		return registers_8[*cpu.reg]
	} else {
		return registers_16[*cpu.reg]
	}
}

func (cpu *Cpu) getImmediateData() string {
	if cpu.w == nil {
		panic("get immediate data requires w")
	}

	if *cpu.w == 0b1 {
		cpu.readNextByte()
		firstByte := cpu.currentByte
		cpu.readNextByte()
		val := int16(binary.LittleEndian.Uint16(append([]byte{}, firstByte, cpu.currentByte)))
		return fmt.Sprintf("%v", val)
	} else {
		cpu.readNextByte()
		val := int8(cpu.currentByte)
		return fmt.Sprintf("%v", val)
	}
}

func (cpu *Cpu) getAccumulatorData() string {

	cpu.readNextByte()
	firstByte := cpu.currentByte
	cpu.readNextByte()
	val := int16(binary.LittleEndian.Uint16(append([]byte{}, firstByte, cpu.currentByte)))
	return fmt.Sprintf("[%v]", val)
}

func (cpu *Cpu) handleImmediateToRegisterMemory(instruction instruction, addressMode addressMode) (string, string) {

	w := cpu.currentByte & 0b1
	cpu.readNextByte()
	mod := cpu.currentByte >> 6 & 0b11
	rm := cpu.currentByte & 0b111

	cpu.w = &w
	cpu.mod = &mod
	cpu.rm = &rm

	dst := cpu.getRegisterOrMemory()
	src := cpu.getImmediateData()
	if *cpu.w == 0b0 {
		src = fmt.Sprintf("byte %v", src)
	} else {
		src = fmt.Sprintf("word %v", src)
	}

	return dst, src
}

func (cpu *Cpu) handleRegisterMemoryToFromRegister(instruction instruction, addressMode addressMode) (string, string) {
	d := cpu.currentByte >> 1 & 0b1

	w := cpu.currentByte & 0b1
	cpu.readNextByte()
	mod := cpu.currentByte >> 6 & 0b11
	reg := cpu.currentByte >> 3 & 0b111
	rm := cpu.currentByte & 0b111

	cpu.d = &d
	cpu.w = &w
	cpu.mod = &mod
	cpu.reg = &reg
	cpu.rm = &rm

	dst := cpu.getRegisterOrMemory()
	src := cpu.getRegister()

	return dst, src
}

func (cpu *Cpu) handleImmediateToRegister(instruction instruction, addressMode addressMode) (string, string) {
	w := (cpu.currentByte >> 3) & 0b1
	reg := cpu.currentByte & 0b111

	cpu.w = &w
	cpu.reg = &reg

	dst := cpu.getRegister()
	src := cpu.getImmediateData()

	return dst, src
}

func (cpu *Cpu) handleMemoryToAccumulator(instruction instruction, addressMode addressMode) (string, string) {
	w := cpu.currentByte & 0b1

	cpu.w = &w

	dst := "ax"
	src := cpu.getAccumulatorData()

	return dst, src
}

func (cpu *Cpu) handleAccumulatorToMemory(instruction instruction, addressMode addressMode) (string, string) {
	w := cpu.currentByte & 0b1

	cpu.w = &w

	src := "ax"
	dst := cpu.getAccumulatorData()

	return dst, src
}

func (cpu *Cpu) exectuteInstruction(instruction instruction, addressMode addressMode) (string, string) {
	switch addressMode {
	case AccumulatorToMemory:
		return cpu.handleAccumulatorToMemory(instruction, addressMode)
	case MemoryToAccumulator:
		return cpu.handleMemoryToAccumulator(instruction, addressMode)
	case ImmediateToRegister:
		return cpu.handleImmediateToRegister(instruction, addressMode)
	case ImmediateToRegisterMemory:
		return cpu.handleImmediateToRegisterMemory(instruction, addressMode)
	case RegisterMemoryToFromRegister:
		return cpu.handleRegisterMemoryToFromRegister(instruction, addressMode)
	default:
		fmt.Printf("%0b %v %v \n", cpu.currentByte, instruction, addressMode)
		panic("not handled instruction")
	}

	panic("not handled instruction")
}

type Cpu struct {
	data        []byte
	currentByte byte
	rm          *byte
	d           *byte
	w           *byte
	reg         *byte
	mod         *byte
}

func (cpu *Cpu) readNextByte() {
	if len(cpu.data) == 0 {
		panic("Cant read data of empty arr")
	}

	res := cpu.data[0]
	cpu.data = cpu.data[1:]
	cpu.currentByte = res

	// return res
}
func (cpu *Cpu) resetIdentifiers() {
	cpu.reg = nil
	cpu.rm = nil
	cpu.w = nil
	cpu.mod = nil
	cpu.d = nil
}

func (cpu *Cpu) runInstructionSet() {
	for len(cpu.data) > 0 {
		cpu.resetIdentifiers()
		instruction, addressMode := cpu.parseInstruction()
		dst, src := cpu.exectuteInstruction(instruction, addressMode)
		if cpu.d != nil && *cpu.d == 1 {
			dst, src = src, dst
		}

		fmt.Printf("%v %v, %v \n", instruction, dst, src)
	}
}

func main() {
	data, err := os.ReadFile("./data/listing_0040_challenge_movs")

	if err != nil {
		log.Fatal(err)
	}

	cpu := Cpu{
		data: data,
	}
	cpu.runInstructionSet()

}
