package main

import (
	"fmt"
	"log"
	"os"
)

var registers_8 = []string{"AL", "CL", "DL", "BL", "AH", "CH", "DH", "BH"}
var registers_16 = []string{"AX", "CX", "DX", "BX", "SP", "BP", "SI", "DI"}

func main() {
	data, err := os.ReadFile("listing_0037_single_register_mov")

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(data)
	// fmt.Printf("%08b\n", data)

	for len(data) > 0 {
		firstByte := data[0]
		if len(data) > 0 {
			data = data[1:]
		}
		// mov operation register/memory to/from register
		if firstByte>>2 == 0b00100010 {
			secondByte := data[0]

			d := (firstByte >> 1) & 1

			w := firstByte & 1
			mod := (secondByte >> 6) & 0b11
			reg := (secondByte >> 3) & 0b111
			rm := secondByte & 0b111

			if mod != 0b11 {
				panic("suported only register to register move")
			}
			registers := registers_8
			if w == 0b1 {
				registers = registers_16
			}
			dst := registers[rm]
			src := registers[reg]
			if d == 0b1 {
				dst, src = src, dst
			}

			if len(data) > 0 {
				data = data[1:]
			}
			fmt.Printf("mov %v, %v \n", dst, src)
		} else {
			panic("Not supported")
		}

	}
}
