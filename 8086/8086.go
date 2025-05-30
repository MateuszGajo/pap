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
	MOVE   instruction = "mov"
	ADD    instruction = "add"
	SUB    instruction = "sub"
	CMP    instruction = "cmp"
	JNZ    instruction = "jnz"
	JE     instruction = "je"
	JL     instruction = "jl"
	JLE    instruction = "jle"
	JB     instruction = "jb"
	JBE    instruction = "jbe"
	JP     instruction = "jp"
	JO     instruction = "jo"
	JS     instruction = "js"
	JNE    instruction = "jne"
	JNL    instruction = "jnl"
	JG     instruction = "jg"
	JNB    instruction = "jnb"
	JA     instruction = "ja"
	JNP    instruction = "jnp"
	JNO    instruction = "jno"
	JNS    instruction = "jns"
	LOOP   instruction = "loop"
	LOOPZ  instruction = "loopz"
	LOOPNZ instruction = "loopnz"
	JCXZ   instruction = "jcxz"
)

const (
	RegisterMemoryToFromRegister                addressMode = "RegisterMemoryToFromRegister"
	ImmediateToRegisterMemory                   addressMode = "ImmediateToRegisterMemory"
	ImmediateToRegister                         addressMode = "ImmediateToRegister"
	MemoryToAccumulator                         addressMode = "MemoryToAccumulator"
	AccumulatorToMemory                         addressMode = "AccumulatorToMemory"
	RegisterMemoryWithRegisterToEither          addressMode = "RegisterMemoryWithRegisterToEither"
	ImmediateToRegisterMemoryWithSignExtenstion addressMode = "ImmediateToRegisterMemoryWithSignExtenstion"
	ImmediateToAccumulator                      addressMode = "ImmediateToAccumulator"
	SignedIncrementToInstructionPointer         addressMode = "SignedIncrementToInstructionPointer"
)

func (cpu *Cpu) parseInstruction() (instruction, addressMode) {
	var instruction instruction
	var addressMode addressMode

	cpu.readNextByte()
	if cpu.currentByte>>2 == 0b00100010 {
		instruction = MOVE
		addressMode = RegisterMemoryToFromRegister
	} else if cpu.currentByte>>1 == 0b01100011 {
		instruction = MOVE
		addressMode = ImmediateToRegisterMemory
	} else if cpu.currentByte>>4 == 0b00001011 {
		instruction = MOVE
		addressMode = ImmediateToRegister
	} else if cpu.currentByte>>1 == 0b01010000 {
		instruction = MOVE
		addressMode = MemoryToAccumulator
	} else if cpu.currentByte>>1 == 0b01010001 {
		instruction = MOVE
		addressMode = AccumulatorToMemory
	} else if cpu.currentByte>>2 == 0b000000 {
		instruction = ADD
		addressMode = RegisterMemoryWithRegisterToEither
	} else if cpu.currentByte>>2 == 0b00100000 && (cpu.nextByte>>3&0b111) == 0b000 {
		instruction = ADD
		addressMode = ImmediateToRegisterMemoryWithSignExtenstion
	} else if cpu.currentByte>>1 == 0b00000010 {
		instruction = ADD
		addressMode = ImmediateToAccumulator
	} else if cpu.currentByte>>2 == 0b00001010 {
		instruction = SUB
		addressMode = RegisterMemoryWithRegisterToEither
	} else if cpu.currentByte>>2 == 0b00100000 && (cpu.nextByte>>3&0b111) == 0b101 {
		instruction = SUB
		addressMode = ImmediateToRegisterMemoryWithSignExtenstion
	} else if cpu.currentByte>>1 == 0b00010110 {
		instruction = SUB
		addressMode = ImmediateToAccumulator
	} else if cpu.currentByte>>2 == 0b00001110 {
		instruction = CMP
		addressMode = RegisterMemoryWithRegisterToEither
	} else if cpu.currentByte>>2 == 0b00100000 && (cpu.nextByte>>3&0b111) == 0b111 {
		instruction = CMP
		addressMode = ImmediateToRegisterMemoryWithSignExtenstion
	} else if cpu.currentByte>>1 == 0b00011110 {
		instruction = CMP
		addressMode = ImmediateToAccumulator
	} else if cpu.currentByte == 0b01110101 {
		instruction = JNZ
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110100 {
		instruction = JE
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111100 {
		instruction = JL
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111110 {
		instruction = JLE
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110010 {
		instruction = JB
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110110 {
		instruction = JBE
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111010 {
		instruction = JP
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110000 {
		instruction = JO
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111000 {
		instruction = JS
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111101 {
		instruction = JNL
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111111 {
		instruction = JG
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110011 {
		instruction = JNB
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110111 {
		instruction = JA
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111011 {
		instruction = JNP
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01110001 {
		instruction = JNO
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b01111001 {
		instruction = JNS
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b011100010 {
		instruction = LOOP
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b011100001 {
		instruction = LOOPZ
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b011100000 {
		instruction = LOOPNZ
		addressMode = SignedIncrementToInstructionPointer
	} else if cpu.currentByte == 0b011100011 {
		instruction = JCXZ
		addressMode = SignedIncrementToInstructionPointer
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
			cpu.readNextByte()
			val := int16(binary.LittleEndian.Uint16(append([]byte{}, cpu.previousByte, cpu.currentByte)))
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
		cpu.readNextByte()
		val := int16(binary.LittleEndian.Uint16(append([]byte{}, cpu.previousByte, cpu.currentByte)))
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
		cpu.readNextByte()
		val := int16(binary.LittleEndian.Uint16(append([]byte{}, cpu.previousByte, cpu.currentByte)))
		return fmt.Sprintf("%v", val)
	} else {
		cpu.readNextByte()
		val := int8(cpu.currentByte)
		return fmt.Sprintf("%v", val)
	}
}

func (cpu *Cpu) getImmediateDataWithSignExtenstion() string {
	if cpu.w == nil || cpu.s == nil {
		panic("get immediate data with sign extenstion requires w and s")
	}

	if *cpu.w == 0b1 && *cpu.s == 0b0 {
		cpu.readNextByte()
		cpu.readNextByte()
		val := int16(binary.LittleEndian.Uint16(append([]byte{}, cpu.previousByte, cpu.currentByte)))
		return fmt.Sprintf("%v", val)
	} else {
		cpu.readNextByte()
		val := int8(cpu.currentByte)
		return fmt.Sprintf("%v", val)
	}
}

func (cpu *Cpu) getAccumulatorData() string {

	cpu.readNextByte()
	cpu.readNextByte()
	val := int16(binary.LittleEndian.Uint16(append([]byte{}, cpu.previousByte, cpu.currentByte)))
	return fmt.Sprintf("[%v]", val)
}

func (cpu *Cpu) handleSignedIncrementToInstructionPointer(instruction instruction, addressMode addressMode) (string, string) {
	cpu.readNextByte()
	return fmt.Sprintf("%v", cpu.currentByte), ""
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
	if *cpu.mod != 0b11 {
		if *cpu.w == 0b0 {
			dst = fmt.Sprintf("byte %v", dst)
		} else {
			dst = fmt.Sprintf("word %v", dst)
		}
	}

	return dst, src
}

func (cpu *Cpu) handleImmediateToRegisterMemoryWithSignExtenstion(instruction instruction, addressMode addressMode) (string, string) {

	s := (cpu.currentByte >> 1) & 0b1
	w := cpu.currentByte & 0b1
	cpu.readNextByte()
	mod := cpu.currentByte >> 6 & 0b11
	rm := cpu.currentByte & 0b111

	cpu.w = &w
	cpu.mod = &mod
	cpu.rm = &rm
	cpu.s = &s

	dst := cpu.getRegisterOrMemory()
	if *cpu.mod != 0b11 {
		if *cpu.w == 0b0 {
			dst = fmt.Sprintf("byte %v", dst)
		} else {
			dst = fmt.Sprintf("word %v", dst)
		}
	}

	src := cpu.getImmediateDataWithSignExtenstion()

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
	if w == 0 {
		dst = "al"
	}
	src := cpu.getAccumulatorData()

	return dst, src
}

func (cpu *Cpu) handleAccumulatorToMemory(instruction instruction, addressMode addressMode) (string, string) {
	w := cpu.currentByte & 0b1

	cpu.w = &w

	dst := cpu.getAccumulatorData()
	src := "ax"
	if w == 0 {
		src = "al"
	}

	return dst, src
}

func (cpu *Cpu) handleRegisterMemoryWithRegisterToEither(instruction instruction, addressMode addressMode) (string, string) {
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

func (cpu *Cpu) handleImmediateToAccumulator(instruction instruction, addressMode addressMode) (string, string) {

	w := cpu.currentByte & 0b1

	cpu.w = &w

	dst := "AX"
	if w == 0 {
		dst = "AL"
	}
	src := cpu.getImmediateData()

	return dst, src
}

func (cpu *Cpu) exectuteInstruction(instruction instruction, addressMode addressMode) (string, string) {
	switch addressMode {
	case AccumulatorToMemory:
		return cpu.handleAccumulatorToMemory(instruction, addressMode)
	case SignedIncrementToInstructionPointer:
		return cpu.handleSignedIncrementToInstructionPointer(instruction, addressMode)
	case MemoryToAccumulator:
		return cpu.handleMemoryToAccumulator(instruction, addressMode)
	case ImmediateToRegister:
		return cpu.handleImmediateToRegister(instruction, addressMode)
	case ImmediateToRegisterMemory:
		return cpu.handleImmediateToRegisterMemory(instruction, addressMode)
	case RegisterMemoryToFromRegister:
		return cpu.handleRegisterMemoryToFromRegister(instruction, addressMode)
	case RegisterMemoryWithRegisterToEither:
		return cpu.handleRegisterMemoryWithRegisterToEither(instruction, addressMode)
	case ImmediateToRegisterMemoryWithSignExtenstion:
		return cpu.handleImmediateToRegisterMemoryWithSignExtenstion(instruction, addressMode)
	case ImmediateToAccumulator:
		return cpu.handleImmediateToAccumulator(instruction, addressMode)
	default:
		fmt.Printf("%0b %v %v \n", cpu.currentByte, instruction, addressMode)
		panic("not handled instruction")
	}

}

type Cpu struct {
	data         []byte
	currentByte  byte
	nextByte     byte
	previousByte byte
	rm           *byte
	d            *byte
	s            *byte
	w            *byte
	reg          *byte
	mod          *byte
}

func (cpu *Cpu) readNextByte() {
	if len(cpu.data) == 0 {
		panic("Cant read data of empty arr")
	}
	cpu.previousByte = cpu.currentByte

	res := cpu.data[0]
	cpu.data = cpu.data[1:]
	cpu.currentByte = res

	if len(cpu.data) > 1 {
		cpu.nextByte = cpu.data[0]
	}

	// return res
}
func (cpu *Cpu) resetIdentifiers() {
	cpu.reg = nil
	cpu.rm = nil
	cpu.w = nil
	cpu.mod = nil
	cpu.d = nil
	cpu.s = nil
}

func (cpu *Cpu) runInstructionSet() {
	for len(cpu.data) > 0 {
		cpu.resetIdentifiers()
		instruction, addressMode := cpu.parseInstruction()
		dst, src := cpu.exectuteInstruction(instruction, addressMode)
		if cpu.d != nil && *cpu.d == 1 {
			dst, src = src, dst
		}
		if src == "" {
			fmt.Printf("%v %v\n", instruction, dst)
		} else {
			fmt.Printf("%v %v, %v \n", instruction, dst, src)
		}

	}
}

func main() {
	data, err := os.ReadFile("./data/listing_0041_add_sub_cmp_jnz")
	// data, err := os.ReadFile("./data/listing_0040_challenge_movs")
	// data, err := os.ReadFile("./data/listing_0039_more_movs")

	if err != nil {
		log.Fatal(err)
	}

	cpu := Cpu{
		data: data,
	}
	cpu.runInstructionSet()

}
